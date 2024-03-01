package messages

import (
	"context"
	"log"

	// "fmt"
	// // "log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	// "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// InboxItem represents an item in the user's inbox.
type InboxItem struct {
	Sender         string `json:"sender"`
	Reciever       string `json:"reciever"`
	S3ResourceLink string `json:"s3_resource_link" dynamodbav:"s3_resource_link"`
	Time           string `json:"time"`
	Classroom string `json:"classroom"`
}

// GetUserInbox retrieves the inbox for a given username.
func GetUserInbox(username string, classID string) ([]InboxItem, error) {
	// Create a new AWS SDK configuration.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create a DynamoDB client.
	client := dynamodb.NewFromConfig(cfg)
	keyCond := expression.Key("reciever").Equal(expression.Value(username))
	keyCond = keyCond.And(expression.Key("classroom").Equal(expression.Value(classID)))

	// Build the expression
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %v", err)
	}

	// Create the QueryInput
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("inbox"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Execute the query.
	result, err := client.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, err
	}

	// Extract the items from the result.
	inboxItems := []InboxItem{}

	for _, v := range result.Items {
		var x InboxItem
		attributevalue.UnmarshalMap(v, &x)
		log.Println(x)
		inboxItems = append(inboxItems, x)
	}

	return inboxItems, nil
}

// GetUserInbox retrieves the inbox for a given username.
func GetUserOutbox(username string, classroom string) ([]InboxItem, error) {
	// Create a new AWS SDK configuration.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create a DynamoDB client.
	client := dynamodb.NewFromConfig(cfg)
	keyCond := expression.Key("sender").Equal(expression.Value(username))
	keyCond = keyCond.And(expression.Key("classroom").Equal(expression.Value(classroom)))

	// Build the expression
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %v", err)
	}

	// Create the QueryInput
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("outbox"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	// Execute the query.
	result, err := client.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, err
	}

	// Extract the items from the result.
	inboxItems := []InboxItem{}

	for _, v := range result.Items {
		var x InboxItem
		attributevalue.UnmarshalMap(v, &x)
		log.Println(x)
		inboxItems = append(inboxItems, x)
	}

	return inboxItems, nil
}
