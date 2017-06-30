package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	awsSession *session.Session
	uploader   *s3manager.Uploader
)

type AWSOption func(a *AWS)

type AWS struct {
	awsSession *session.Session
	uploader   *s3manager.Uploader
	bucket     *string
}

func (a *AWS) Options(options ...AWSOption) {
	for _, option := range options {
		option(a)
	}
}

func WithBucket(s string) func(a *AWS) {
	return func(a *AWS) {
		a.bucket = aws.String(s)
	}
}

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
