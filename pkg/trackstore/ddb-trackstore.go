package trackstore

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	DEFAULT_TABLE_NAME = "DoYouTrackIdStore"
	TABLE_NAME_ENV_VAR = "DDB_DO_YOU_TRACK_STORE_TABLE_NAME"
)

type DDBTrackStore struct {
	tableName string
	dynamoDb  *dynamodb.DynamoDB
}

type DDBTrackStoreRecord struct {
	TrackId   string
	DateAdded string
}

func NewDdbTrackStore() (*DDBTrackStore, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	tableName := getTableStoreTableName()

	err = createTrackStoreTableIfDoesNotExist(svc, tableName)
	if err != nil {
		return nil, err
	}

	return &DDBTrackStore{
		tableName: tableName,
		dynamoDb:  svc,
	}, nil
}

func (s *DDBTrackStore) AddTrack(trackId string) error {
	record := DDBTrackStoreRecord{
		TrackId:   trackId,
		DateAdded: time.Now().String(),
	}

	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		log.Fatalf("Got error marshalling new track store item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.tableName),
	}

	_, err = s.dynamoDb.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added track with id '" + record.TrackId + "' to table " + s.tableName)

	return nil
}

// TODO - Batch write
func (s *DDBTrackStore) AddTracks(trackIds []string) error {
	for _, trackId := range trackIds {
		if !s.HasTrack(trackId) {
			err := s.AddTrack(trackId)
			if err != nil {
				fmt.Printf("Unable to add track with ID %s due to error %s\n", trackId, err)
			}
		}
	}

	return nil
}

func (s *DDBTrackStore) HasTrack(trackId string) bool {
	result, err := s.dynamoDb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(getTableStoreTableName()),
		Key: map[string]*dynamodb.AttributeValue{
			"TrackId": {
				S: aws.String(trackId),
			},
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return false
			default:
				log.Fatalf("Unexpected error occured when calling DescribeTable: %e", err)
			}
		}
	}

	if result.Item == nil {
		return false
	}

	return true
}

func createTrackStoreTableIfDoesNotExist(svc *dynamodb.DynamoDB, tableName string) error {
	table, err := svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(getTableStoreTableName())})

	if table.Table != nil {
		fmt.Printf("Table %s already exists, skipping creation.\n", getTableStoreTableName())
		return nil
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Printf("Table %s does not exist, bootstrapping...\n", getTableStoreTableName())
				break
			default:
				log.Fatalf("Unexpected error occured when calling DescribeTable: %e", err)
			}
		}
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("TrackId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("TrackId"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		TableName:   aws.String(tableName),
	}

	_, err = svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
		return err
	}

	fmt.Println("Created the table", tableName)

	return nil
}

func getTableStoreTableName() string {
	if value, ok := os.LookupEnv(TABLE_NAME_ENV_VAR); ok {
		return value
	}
	return DEFAULT_TABLE_NAME
}
