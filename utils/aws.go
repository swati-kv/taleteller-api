package utils

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"taleteller/app"
)

type Service interface {
	UploadFile(bucket string, request UploadS3, isPublic bool) (string, error)
	UploadFileV2(bucket string, request UploadS3, isPublic bool) (path string, err error)
}

type service struct {
}

func NewAWSService() Service {
	return &service{}
}

// Function to upload file to s3
func (s *service) UploadFile(bucket string, request UploadS3, isPublic bool) (path string, err error) {
	// The session the S3 Uploader will use
	serviceConfig := app.InitServiceConfig()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(serviceConfig.GetAWSRegion()),
	})
	if err != nil {
		return
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	key := s.generateKey(request.FileType, request.FileFormat, request.FileName)

	inputRequest := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(request.File)),
	}
	fmt.Println("inputRequest", inputRequest)
	if isPublic {
		inputRequest.ACL = aws.String("public-read")
	}

	// Upload the file to S3.
	result, err := uploader.Upload(inputRequest)
	if err != nil {
		return path, fmt.Errorf("failed to upload file, %v", err)
	}

	return result.Location, err
}

func (s *service) UploadFileV2(bucket string, request UploadS3, isPublic bool) (path string, err error) {
	// The session the S3 Uploader will use
	serviceConfig := app.InitServiceConfig()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(serviceConfig.GetAWSRegion()),
	})
	if err != nil {
		return
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	key := s.generateKey(request.FileType, request.FileFormat, request.FileName)

	var contentType *string
	//var conEnc *string

	videoContent := "video/mp4"
	cType := "base64"
	if request.FileType == "video" {
		contentType = &videoContent
		//	conEnc = &cType
	}
	fmt.Println("contentTypecontentType", *contentType, cType)
	fmt.Println("request.File", request.File)
	inputRequest := &s3manager.UploadInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        bytes.NewReader(request.FileBytes),
		ContentType: contentType,
		//ContentEncoding: conEnc,
	}
	fmt.Println("inputRequest", inputRequest)
	if isPublic {
		inputRequest.ACL = aws.String("public-read")
	}

	// Upload the file to S3.
	result, err := uploader.Upload(inputRequest)
	if err != nil {
		return path, fmt.Errorf("failed to upload file, %v", err)
	}

	return result.Location, err
}

func (s *service) generateKey(fileType string, fileFormat string, fileName string) (key string) {
	if len(fileName) == 0 {
		fileName = uuid.New().String()
	}

	key = fileType + "/" + fileName + "." + fileFormat

	return
}
