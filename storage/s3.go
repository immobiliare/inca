package storage

import (
	"bytes"
	"encoding/pem"
	"fmt"
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
	if len(options) != 3 {
		return fmt.Errorf("invalid number of options for provider %s: %s", s.ID(), options)
	}

	s.config = &aws.Config{
		Credentials:      credentials.NewStaticCredentials(options[1], options[2], ""),
		Endpoint:         aws.String(options[0]),
		Region:           aws.String("eu-west-1"),
		DisableSSL:       aws.Bool(strings.HasPrefix(options[0], "http://")),
		S3ForcePathStyle: aws.Bool(true),
	}
	return nil
}

func (s *S3) Get(name string) ([]byte, []byte, error) {
	session, err := session.NewSession(s.config)
	if err != nil {
		return nil, nil, err
	}
	client := s3manager.NewDownloader(session)

	crtData := aws.NewWriteAtBuffer([]byte{})
	if _, err := client.Download(crtData, &s3.GetObjectInput{
		Bucket: bucket(name),
		Key:    &s3CrtName,
	}); err != nil {
		return nil, nil, err
	}

	keyData := aws.NewWriteAtBuffer([]byte{})
	if _, err := client.Download(keyData, &s3.GetObjectInput{
		Bucket: bucket(name),
		Key:    &s3KeyName,
	}); err != nil {
		return nil, nil, err
	}

	return crtData.Bytes(), keyData.Bytes(), nil
}

func (s *S3) Put(name string, crtData *pem.Block, keyData *pem.Block) error {
	session, err := session.NewSession(s.config)
	if err != nil {
		return err
	}
	client := s3.New(session)

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
	session, err := session.NewSession(s.config)
	if err != nil {
		return err
	}
	client := s3.New(session)

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
