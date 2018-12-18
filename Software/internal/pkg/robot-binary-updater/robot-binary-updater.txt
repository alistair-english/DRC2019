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
	"time"
)

const (
	S3_REGION = "ap-southeast-2"
	S3_BUCKET = "drc2019"
	BIN_NAME  = "line-detection"
)

var hasPerfomred bool
var localDate time.Time
var processIsRunning bool

func main() {
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}
	binaryPath := os.Getenv("GOPATH") + "/bin/"
	p := exec.Command(binaryPath + BIN_NAME)

	// Get Lastest on Startup

	for {
		CheckIfCurrent(s, p)
		time.Sleep(5e+9)
	}
}

func GetBinary(s *session.Session, p *exec.Cmd) {
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

	// Restart the Process
	StartProcess(p)
}

func EndProcess(p *exec.Cmd) {
	if processIsRunning {
		fmt.Println("Killing Process")
		p.Process.Kill()
		processIsRunning = false
	}
}

func StartProcess(p *exec.Cmd) {
	if !processIsRunning {
		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
		processIsRunning = true
	}
}

func CheckIfCurrent(s *session.Session, p *exec.Cmd) {
	if hasPerfomred == false {
		localDate = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
		hasPerfomred = true
	}
	serverDate := GetDate(s)
	if serverDate.After(localDate) {
		fmt.Println("New Version Available")
		localDate = serverDate
		EndProcess(p)
		GetBinary(s, p)
	}
}

func GetDate(s *session.Session) time.Time {
	svc := s3.New(s)
	response, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(S3_BUCKET),
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range response.Contents {
		if *item.Key == BIN_NAME {
			return *item.LastModified
		}
	}
	log.Fatal("File not Found")
	return time.Now()
}
