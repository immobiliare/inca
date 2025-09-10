package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/immobiliare/inca/pki"
	"github.com/immobiliare/inca/util"
	"github.com/rs/zerolog/log"
)

var (
	s3CrtName = "crt.pem"
	s3KeyName = "key.pem"
)

type S3 struct {
	Storage
	config *aws.Config
}

func (s S3) ID() string {
	return "S3"
}

func (s *S3) Tune(options map[string]interface{}) error {
	endpoint, ok := options["endpoint"]
	if !ok {
		return fmt.Errorf("provider %s: endpoint not defined", s.ID())
	}

	access, ok := options["access"]
	if !ok {
		return fmt.Errorf("provider %s: access not defined", s.ID())
	}

	secret, ok := options["secret"]
	if !ok {
		return fmt.Errorf("provider %s: secret not defined", s.ID())
	}

	region, ok := options["region"]
	if !ok {
		return fmt.Errorf("provider %s: region not defined", s.ID())
	}

	s.config = aws.NewConfig().
		WithEndpoint(endpoint.(string)).
		WithDisableSSL(strings.HasPrefix(endpoint.(string), "http://")).
		WithRegion(region.(string)).
		WithS3ForcePathStyle(true).
		WithCredentials(credentials.NewStaticCredentials(
			access.(string),
			secret.(string),
			"",
		))
	return nil
}

func (s *S3) Get(name string) ([]byte, []byte, error) {
	bucketName, err := nameToBucket(name)
	if err != nil {
		return nil, nil, err
	}

	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	crtData := bytes.NewBuffer(nil)
	if data, err := client.GetObject(&s3.GetObjectInput{
		Bucket: bucketName,
		Key:    &s3CrtName,
	}); err != nil {
		return nil, nil, err
	} else {
		if _, err := io.Copy(crtData, data.Body); err != nil {
			return nil, nil, errors.Join(err, data.Body.Close())
		}
		if err := data.Body.Close(); err != nil {
			return nil, nil, err
		}
	}

	keyData := bytes.NewBuffer(nil)
	if data, err := client.GetObject(&s3.GetObjectInput{
		Bucket: bucketName,
		Key:    &s3KeyName,
	}); err != nil {
		return nil, nil, err
	} else {
		if _, err := io.Copy(keyData, data.Body); err != nil {
			return nil, nil, errors.Join(err, data.Body.Close())
		}
		return crtData.Bytes(), keyData.Bytes(), data.Body.Close()
	}
}

func (s *S3) Put(name string, crtData, keyData []byte) error {
	bucketName, err := nameToBucket(name)
	if err != nil {
		return err
	}

	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	if _, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: bucketName}); err != nil &&
		strings.HasPrefix(err.Error(), s3.ErrCodeBucketAlreadyExists) &&
		strings.HasPrefix(err.Error(), s3.ErrCodeBucketAlreadyOwnedByYou) {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    &s3CrtName,
		Body:   bytes.NewReader(crtData),
	}); err != nil {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    &s3KeyName,
		Body:   bytes.NewReader(keyData),
	}); err != nil {
		return err
	}

	return nil
}

func (s *S3) Del(name string) error {
	bucketName, err := nameToBucket(name)
	if err != nil {
		return err
	}

	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	if err := s3manager.NewBatchDeleteWithClient(client).Delete(
		aws.BackgroundContext(),
		s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
			Bucket: bucketName,
		})); err != nil {
		return err
	}

	// Bucket is left intact - no longer deleting the bucket
	return nil
}

func (s *S3) Renew(name string, crtData, keyData []byte) error {
	bucketName, err := nameToBucket(name)
	if err != nil {
		return err
	}

	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    &s3CrtName,
		Body:   bytes.NewReader(crtData),
	}); err != nil {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    &s3KeyName,
		Body:   bytes.NewReader(keyData),
	}); err != nil {
		return err
	}

	return nil
}

func (s *S3) Find(filters ...string) ([][]byte, error) {
	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	results := [][]byte{}
	for _, bucket := range buckets.Buckets {
		if !matchFilters(bucket.Name, filters) {
			continue
		}

		crt, _, err := s.Get(*bucket.Name)
		if err != nil {
			log.Error().Err(err).Msg("storage/s3: skip empty buckets with missing certificates")
			// Skip empty buckets or buckets with missing certificates
			// This can happen after certificate deletion when bucket is left intact
			continue
		}

		results = append(results, crt)
	}

	return results, nil
}

func (s *S3) Config() map[string]string {
	return map[string]string{
		"Endpoint": *s.config.Endpoint,
		"Region":   *s.config.Region,
	}
}

func nameToBucket(name string) (*string, error) {
	bucket := strings.TrimPrefix(name, "*.")
	if !validateBucketName(bucket) {
		return nil, errors.New("unsupported CN (protocol violation)")
	}

	// Buckets used with Amazon S3 Transfer Acceleration can't have dots (.) in their names
	bucket = strings.ReplaceAll(bucket, ".", "-")
	return &bucket, nil
}

func matchFilters(bucket *string, filters []string) bool {
	name := strings.ReplaceAll(*bucket, "-", ".")
	wildcardName := "*." + name

	return pki.IsValidCN(name) && util.RegexesMatch(name, filters...) ||
		pki.IsValidCN(wildcardName) && util.RegexesMatch(wildcardName, filters...)
}

func validateBucketName(name string) bool {

	// Bucket names cannot be formatted as IP addresses
	if net.ParseIP(name) != nil {
		return false
	}

	// Bucket names must not start with the prefix xn--
	if strings.HasPrefix(name, "xn--") {
		return false
	}

	// Bucket names must not start with the prefix sthree-
	if strings.HasPrefix(name, "sthree-") {
		return false
	}

	// Bucket names must not end with the suffix -s3alias
	if strings.HasSuffix(name, "-s3alias") {
		return false
	}

	// Bucket names must not end with the suffix --ol-s3
	if strings.HasSuffix(name, "--ol-s3") {
		return false
	}

	// Bucket names can be between 3 and 63 characters long
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	// Bucket names must not contain uppercase characters
	if name != strings.ToLower(name) {
		return false
	}

	// Bucket names must not contain underscores
	if strings.Contains(name, "_") {
		return false
	}

	// Bucket names must be a series of one or more labels
	for _, label := range strings.Split(name, ".") {

		// Adjacent labels are separated by a single period (.)
		if len(label) < 1 {
			return false
		}

		// Each label can contain lowercase letters, numbers, and hyphens
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return false
			}
		}

		// Each label must start and end with a lowercase letter or a number
		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
	}

	return true
}
