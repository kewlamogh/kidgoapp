package messages

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func SendMessage(from string, to string, s3url string, classroom string) error {
	// put in inbox and outbox
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	sender, err := attributevalue.Marshal(from)
	if err != nil {
		return err
	}

	reciever, err := attributevalue.Marshal(to)
	if err != nil {
		return err
	}

	url, err := attributevalue.Marshal(s3url)
	if err != nil {
		return err
	}

	classID, err := attributevalue.Marshal(classroom)
	if err != nil {
		return err
	}

	timeStr := time.Now().String()
	timestamp, err := attributevalue.Marshal(timeStr) 

	// Create a DynamoDB client.
	client := dynamodb.NewFromConfig(cfg)
	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("inbox"),
		Item: map[string]types.AttributeValue{
			"sender": sender,
			"reciever": reciever,
			"s3_resource_link": url, 
			"time": timestamp, 
			"classroom": classID, 
		},
	})

	client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("outbox"),
		Item: map[string]types.AttributeValue{
			"sender": sender,
			"reciever": reciever,
			"s3_resource_link": url, 
			"time": timestamp, 
			"classroom": classID,
		},
	})

	return nil
}