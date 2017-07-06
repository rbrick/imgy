package storage

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AWSOption represents an option which can be used to configure
// an AWS instance
type AWSOption func(a *AWS)

// AWS is a wrapper for common AWS functions
type AWS struct {
	awsSession *session.Session
	uploader   *s3manager.Uploader
	bucket     *string
}

// Options sets all the options for the struct
func (a *AWS) Options(options ...AWSOption) {
	for _, option := range options {
		option(a)
	}
}

// WithBucket is an Option for the AWS struct
// This sets the bucket that will be used for uploading files
func WithBucket(s string) func(a *AWS) {
	return func(a *AWS) {
		a.bucket = aws.String(s)
	}
}

// InitAWS initiates a new AWS helper. Optionally can include options
func InitAWS(config *aws.Config, options ...AWSOption) (*AWS, error) {
	sess, err := session.NewSession(config)

	if err != nil {
		return nil, err
	}

	aws := &AWS{
		awsSession: sess,
		uploader:   s3manager.NewUploader(sess),
	}

	aws.Options(options...)

	return aws, nil
}

// Upload uploads a file to the Amazon S3 Bucket
func (a *AWS) Upload(key, contentType string, data []byte) (*s3manager.UploadOutput, error) {
	input := &s3manager.UploadInput{
		Bucket:      a.bucket,
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	}

	res, err := a.uploader.Upload(input)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get returns an image from AWS
// func (a *AWS) Get(key string) (*imgy.Image, error) {
// 	return nil, nil
// }

// func main() {
// 	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	uploader := s3manager.NewUploader(sess)

// 	file, err := os.Open("rm.gif")

// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	uploadInput := s3manager.UploadInput{
// 		Bucket:      aws.String("imgy-s3"),
// 		Key:         aws.String("perfect-solution-3.gif"),
// 		Body:        file,
// 		ContentType: aws.String("image/gif"),
// 	}

// 	result, err := uploader.Upload(&uploadInput)

// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	log.Println("Upload ID:", result.UploadID, "Location:", result.Location)
// }
