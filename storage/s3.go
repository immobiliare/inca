package storage

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/util"
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
	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	crtData := bytes.NewBuffer(nil)
	if data, err := client.GetObject(&s3.GetObjectInput{
		Bucket: nameToBucket(name),
		Key:    &s3CrtName,
	}); err != nil {
		return nil, nil, err
	} else {
		if _, err := io.Copy(crtData, data.Body); err != nil {
			return nil, nil, err
		}
		data.Body.Close()
	}

	keyData := bytes.NewBuffer(nil)
	if data, err := client.GetObject(&s3.GetObjectInput{
		Bucket: nameToBucket(name),
		Key:    &s3KeyName,
	}); err != nil {
		return nil, nil, err
	} else {
		if _, err := io.Copy(keyData, data.Body); err != nil {
			return nil, nil, err
		}
		data.Body.Close()
	}

	return crtData.Bytes(), keyData.Bytes(), nil
}

func (s *S3) Put(name string, crtData, keyData []byte) error {
	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	if _, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: nameToBucket(name)}); err != nil &&
		strings.HasPrefix(err.Error(), s3.ErrCodeBucketAlreadyExists) &&
		strings.HasPrefix(err.Error(), s3.ErrCodeBucketAlreadyOwnedByYou) {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: nameToBucket(name),
		Key:    &s3CrtName,
		Body:   bytes.NewReader(crtData),
	}); err != nil {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: nameToBucket(name),
		Key:    &s3KeyName,
		Body:   bytes.NewReader(keyData),
	}); err != nil {
		return err
	}

	return nil
}

func (s *S3) Del(name string) error {
	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	if err := s3manager.NewBatchDeleteWithClient(client).Delete(
		aws.BackgroundContext(),
		s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
			Bucket: nameToBucket(name),
		})); err != nil {
		return err
	}

	if _, err := client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: nameToBucket(name),
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
		if !(pki.IsValidCN(bucketToName(bucket.Name)) &&
			util.RegexesMatch(bucketToName(bucket.Name), filters...)) {
			continue
		}

		crt, _, err := s.Get(*bucket.Name)
		if err != nil {
			return nil, err
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

func nameToBucket(name string) *string {
	bucket := strings.ReplaceAll(name, ".", "-")
	return &bucket
}

func bucketToName(bucket *string) string {
	return strings.ReplaceAll(*bucket, "-", ".")
}
