package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"os/exec"
)

const (
	S3_REGION = "ap-southeast-2"
	S3_BUCKET = "drc2019"
	BIN_NAME  = "line-detection"
)

func main() {
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}

	// Get Lastest on Startup
	getBinary(s)
}

func getBinary(s *session.Session) {
	binaryPath := os.Getenv("GOPATH") + "/bin/"
	file, err := os.Create(binaryPath + BIN_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(s)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(BIN_NAME),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Update Perms so we can Execute The File
	perms := exec.Command("chmod", "+x", binaryPath+BIN_NAME)
	perms.Start()
}

func StartProcess() {

}

// func CheckIfCurrent(s *session.Session) error {
// 	// TODO
// }

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
