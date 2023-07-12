package main

import (
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	sess, err := session.NewSession(
		&aws.Config{
			Endpoint:    aws.String("http://localhost:4566"),
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
		},
	)
	if err != nil {
		log.Fatalf("failed to create aws session: %v", err)
	}

	svc := sqs.New(sess)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err = receiveMessages(svc); err != nil {
				log.Printf("failed to receive message: %v", err)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

func receiveMessages(svc *sqs.SQS) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String("http://localhost:4566/000000000000/my-queue"),
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(0),
	}

	result, err := svc.ReceiveMessage(input)
	if err != nil {
		return err
	}

	for _, msg := range result.Messages {
		log.Printf("received message: %s", *msg.Body)

		// Delete the message from the queue
		_, err = svc.DeleteMessage(
			&sqs.DeleteMessageInput{
				QueueUrl:      aws.String("http://localhost:4566/000000000000/my-queue"),
				ReceiptHandle: msg.ReceiptHandle,
			},
		)
		if err != nil {
			log.Printf("failed to delete message from SQS: %v", err)
		}
	}

	return nil
}
