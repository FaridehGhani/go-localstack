package cloud

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewAWS() *session.Session {
	var newAWSSession *session.Session
	var err error
	var once sync.Once

	once.Do(func() {
		newAWSSession, err = session.NewSession(
			&aws.Config{
				Endpoint:    aws.String("http://localhost:4566"),
				Region:      aws.String("us-east-1"),
				Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
			},
		)
		if err != nil {
			log.Fatalf("failed to create aws session: %v", err)
		}
	})

	return newAWSSession
}
