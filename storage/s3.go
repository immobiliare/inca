package storage

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.rete.farm/sistemi/inca/pki"
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

func (s *S3) Tune(options ...string) error {
	if len(options) != 4 {
		return fmt.Errorf("invalid number of options for provider %s: %s", s.ID(), options)
	}

	s.config = aws.NewConfig().
		WithEndpoint(options[0]).
		WithDisableSSL(strings.HasPrefix(options[0], "http://")).
		WithRegion(options[3]).
		WithS3ForcePathStyle(true).
		WithCredentials(credentials.NewStaticCredentials(
			options[1],
			options[2],
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
		Bucket: bucket(name),
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
		Bucket: bucket(name),
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

func (s *S3) Put(name string, crtData *pem.Block, keyData *pem.Block) error {
	client := s3.New(
		session.Must(session.NewSession()),
		s.config,
	)

	if _, err := client.CreateBucket(&s3.CreateBucketInput{Bucket: bucket(name)}); err != nil &&
		strings.HasPrefix(err.Error(), s3.ErrCodeBucketAlreadyExists) &&
		strings.HasPrefix(err.Error(), s3.ErrCodeBucketAlreadyOwnedByYou) {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: bucket(name),
		Key:    &s3CrtName,
		Body:   bytes.NewReader(pki.ExportBytes(crtData)),
	}); err != nil {
		return err
	}

	if _, err := client.PutObject(&s3.PutObjectInput{
		Bucket: bucket(name),
		Key:    &s3KeyName,
		Body:   bytes.NewReader(pki.ExportBytes(keyData)),
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
			Bucket: bucket(name),
		})); err != nil {
		return err
	}

	if _, err := client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: bucket(name),
	}); err != nil {
		return err
	}

	return nil
}

func bucket(name string) *string {
	bucket := strings.ReplaceAll(name, ".", "-")
	return &bucket
}
