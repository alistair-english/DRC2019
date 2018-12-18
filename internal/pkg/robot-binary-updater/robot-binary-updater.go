package main

import (
	// "fmt"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	S3_REGION = "ap-southeast-2"
	S3_BUCKET = "drc2019"
	BIN_NAME  = "test"
)

var started bool
var process bool
var localFileDate time.Time
var currentFileDate time.Time
var p *exec.Cmd

func main() {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
	})
	if err != nil {
		log.Fatal(err)
	}
	binaryPath := os.Getenv("GOPATH") + "/bin/"

	for {
		// Check if New Version is Available
		if !isCurrent(s) {
			if !process {
				pullBinary(s)
				p = exec.Command(binaryPath + BIN_NAME)
				if err := p.Start(); err != nil {
					log.Fatal(err)
				}
				process = true
			} else {
				p.Process.Kill()
				pullBinary(s)
				p = exec.Command(binaryPath + BIN_NAME)
				if err := p.Start(); err != nil {
					log.Fatal(err)
				}
			}
		}
		time.Sleep(5e+9)
	}
}

func isCurrent(s *session.Session) bool {
	// Check if this is the first occurance
	if !started {
		localFileDate = time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC)
		started = true
	}
	svc := s3.New(s)
	response, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(S3_BUCKET),
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range response.Contents {
		if *item.Key == BIN_NAME {
			currentFileDate = *item.LastModified
		}
	}
	if currentFileDate.After(localFileDate) {
		localFileDate = currentFileDate
		fmt.Println("New Version Available")
		return false
	} else {
		return true
	}
}

func pullBinary(s *session.Session) {
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
