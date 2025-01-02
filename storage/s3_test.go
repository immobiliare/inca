package storage

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/matryer/is"
)

func TestS3_Tune(t *testing.T) {
	t.Parallel()

	s := &S3{}
	options := map[string]interface{}{
		"endpoint": "https://s3.amazonaws.com",
		"access":   "access_key",
		"secret":   "secret_key",
		"region":   "us-west-2",
	}

	err := s.Tune(options)
	test := is.New(t)
	test.NoErr(err)

	expectedConfig := aws.NewConfig().
		WithEndpoint("https://s3.amazonaws.com").
		WithDisableSSL(false).
		WithRegion("us-west-2").
		WithS3ForcePathStyle(true).
		WithCredentials(credentials.NewStaticCredentials(
			"access_key",
			"secret_key",
			"",
		))

	test.Equal(s.config.Endpoint, expectedConfig.Endpoint)
	test.Equal(s.config.DisableSSL, expectedConfig.DisableSSL)
	test.Equal(s.config.Region, expectedConfig.Region)
	test.Equal(s.config.S3ForcePathStyle, expectedConfig.S3ForcePathStyle)
}

func TestS3_Tune_MissingEndpoint(t *testing.T) {
	t.Parallel()

	s := &S3{}
	options := map[string]interface{}{
		"access": "access_key",
		"secret": "secret_key",
		"region": "us-west-2",
	}

	err := s.Tune(options)
	test := is.New(t)
	test.Equal(err, fmt.Errorf("provider %s: endpoint not defined", s.ID()))
}

func TestS3_Tune_MissingAccess(t *testing.T) {
	t.Parallel()

	s := &S3{}
	options := map[string]interface{}{
		"endpoint": "https://s3.amazonaws.com",
		"secret":   "secret_key",
		"region":   "us-west-2",
	}

	err := s.Tune(options)
	test := is.New(t)
	test.Equal(err, fmt.Errorf("provider %s: access not defined", s.ID()))
}

func TestS3_Tune_MissingSecret(t *testing.T) {
	t.Parallel()

	s := &S3{}
	options := map[string]interface{}{
		"endpoint": "https://s3.amazonaws.com",
		"access":   "access_key",
		"region":   "us-west-2",
	}

	err := s.Tune(options)
	test := is.New(t)
	test.Equal(err, fmt.Errorf("provider %s: secret not defined", s.ID()))
}

func TestS3_Tune_MissingRegion(t *testing.T) {
	t.Parallel()

	s := &S3{}
	options := map[string]interface{}{
		"endpoint": "https://s3.amazonaws.com",
		"access":   "access_key",
		"secret":   "secret_key",
	}

	err := s.Tune(options)
	test := is.New(t)
	test.Equal(err, fmt.Errorf("provider %s: region not defined", s.ID()))
}

func TestS3_Config(t *testing.T) {
	t.Parallel()

	s := &S3{
		config: &aws.Config{
			Endpoint: aws.String("https://s3.amazonaws.com"),
			Region:   aws.String("us-west-2"),
		},
	}

	expected := map[string]string{
		"Endpoint": "https://s3.amazonaws.com",
		"Region":   "us-west-2",
	}

	result := s.Config()

	test := is.New(t)
	test.Equal(result, expected)
}

func Test_nameToBucket(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "127.0.0.1",
			wantErr: true,
		},
		{
			name:    "::1",
			wantErr: true,
		},
		{
			name:    "xn--wgv71a119e.idn.icann.org",
			wantErr: true,
		},
		{
			name:    "sthree-testprefix.it",
			wantErr: true,
		},
		{
			name:    "testsuffix.-s3alias",
			wantErr: true,
		},
		{
			name:    "testsuffix.--ol-s3",
			wantErr: true,
		},
		{
			name:    "underscore_test.it",
			wantErr: true,
		},
		{
			name:    ".dot-start.com",
			wantErr: true,
		},
		{
			name:    "example.com",
			wantErr: false,
		},
		{
			name:    "example.com.",
			wantErr: true,
		},
		{
			name:    "*.example.com",
			wantErr: false,
		},
		{
			name:    "*.example.com.",
			wantErr: true,
		},
		{
			name:    "0x20eNc0d1Ng.com.",
			wantErr: true,
		},
		{
			name:    "-hypen-start.it",
			wantErr: true,
		},
		{
			name:    "hypen-end-.it",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := nameToBucket(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("nameToBucket(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}

func Test_matchFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bucket  string
		filters []string
		want    bool
	}{
		{
			name:    "match direct name",
			bucket:  "example.com",
			filters: []string{"example.com"},
			want:    true,
		},
		{
			name:    "match name with hyphen",
			bucket:  "example-com",
			filters: []string{"example.com"},
			want:    true,
		},
		{
			name:    "no match",
			bucket:  "example-com",
			filters: []string{"different.com"},
			want:    false,
		},
		{
			name:    "match regex filter",
			bucket:  "test-example-com",
			filters: []string{".+\\.example\\.com"},
			want:    true,
		},
		{
			name:    "math wildcard filter",
			bucket:  "test-example-com",
			filters: []string{"*.example.com"},
		},
		{
			name:    "empty filters",
			bucket:  "example-com",
			filters: []string{},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := is.New(t)
			got := matchFilters(&tt.bucket, tt.filters)
			test.Equal(got, tt.want)
		})
	}
}
