package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	sqsClient := sqs.New(sess)
	dynamodbClient := dynamodb.New(sess)

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			<-ticker.C

			message := Message{
				Id:          generateUUID(),
				Description: "new message",
				CreatedAt:   time.Now().Add(24 * time.Hour),
			}

			_, err = dynamodbClient.PutItem(
				&dynamodb.PutItemInput{
					TableName: aws.String("messages"),
					Item: map[string]*dynamodb.AttributeValue{
						"Id": {
							S: aws.String(message.Id),
						},
						"Description": {
							S: aws.String(message.Description),
						},
						"CreatedAt": {
							S: aws.String(message.CreatedAt.Format(time.RFC3339)),
						},
					},
				},
			)
			if err != nil {
				log.Printf("failed to save message: %v", err)
			}

			if err = sendMessage(sqsClient, message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
		}
	}()

	select {}
}

type Message struct {
	Id          string    `json:"Id"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

func generateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}

	return id.String()
}

func sendMessage(svc *sqs.SQS, item Message) error {
	message := fmt.Sprintf(`{"id":"%s","description":"%s","dueDate":"%s"}`, item.Id, item.Description, item.CreatedAt)

	input := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  &message,
		QueueUrl:     aws.String("http://localhost:4566/000000000000/my-queue"),
	}

	output, err := svc.SendMessage(input)
	if err != nil {
		return err
	}
	log.Printf("message sent with id %v", output.MessageId)

	return nil
}
