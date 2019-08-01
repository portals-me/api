package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
)

func AsDynamoDBAttributeValue(old events.DynamoDBAttributeValue) *dynamodb.AttributeValue {
	if old.DataType() == events.DataTypeBoolean {
		return &dynamodb.AttributeValue{
			BOOL: aws.Bool(old.Boolean()),
		}
	} else if old.DataType() == events.DataTypeNumber {
		number := old.Number()

		return &dynamodb.AttributeValue{
			N: &number,
		}
	} else if old.DataType() == events.DataTypeString {
		return &dynamodb.AttributeValue{
			S: aws.String(old.String()),
		}
	} else if old.DataType() == events.DataTypeList {
		list := old.List()
		newList := make([]*dynamodb.AttributeValue, len(list))

		for i, v := range list {
			newList[i] = AsDynamoDBAttributeValue(v)
		}

		return &dynamodb.AttributeValue{
			L: newList,
		}
	} else if old.DataType() == events.DataTypeMap {
		kv := old.Map()
		newKv := make(map[string]*dynamodb.AttributeValue)

		for k, v := range kv {
			newKv[k] = AsDynamoDBAttributeValue(v)
		}

		return &dynamodb.AttributeValue{
			M: newKv,
		}
	}

	return nil
}

func AsDynamoDBAttributeValues(old map[string]events.DynamoDBAttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	new := map[string]*dynamodb.AttributeValue{}

	for key, value := range old {
		new[key] = AsDynamoDBAttributeValue(value)
	}

	return new, nil
}
