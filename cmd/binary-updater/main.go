package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"os"
)

const (
	S3_REGION = "ap-southeast-2"
	S3_BUCKET = "drc2019"
)

func main() {
	// Create S3 Session
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}

	// Upload File
	// err = AddFileToS3(s, "main.go")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// List Files
	err = ListFiles(s)
}

func AddFileToS3(s *session.Session, fileDir string) error {
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(fileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func ListFiles(s *session.Session) error {
	svc := s3.New(s)
	response, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(S3_BUCKET),
	})
	if err != nil {
		return err
	}
	for _, item := range response.Contents {
		if *item.Key == "main.go" {
			fmt.Println("Last Modified: ", *item.LastModified)
		}
	}
	return err
}
