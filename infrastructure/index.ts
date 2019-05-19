import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as fs from 'fs';
import './auth';

const config = {
  service: new pulumi.Config().name,
  stage: pulumi.getStack(),
};

const ApiRole = new aws.iam.Role(`api-role`, {
  assumeRolePolicy: aws.iam.getPolicyDocument({
    statements: [{
      actions: ['sts:AssumeRole'],
      principals: [{
        identifiers: ['appsync.amazonaws.com'],
        type: 'Service'
      }]
    }]
  }).then(result => result.json),
});

const GraphQLApi = new aws.appsync.GraphQLApi(`api`, {
  authenticationType: 'API_KEY',
  logConfig: {
    cloudwatchLogsRoleArn: ApiRole.arn,
    fieldLogLevel: 'ALL',
  },
  schema: fs.readFileSync('./schema.graphql').toString(),
});
const ApiKey = new aws.appsync.ApiKey(`api-key`, {
  apiId: GraphQLApi.id,
});

const CollectionTable = new aws.dynamodb.Table(`appsync-collection`, {
  attributes: [
    {
      name: 'id',
      type: 'S'
    },
    {
      name: 'sort-id',
      type: 'S'
    }
  ],
  hashKey: 'id',
  rangeKey: 'sort-id',
  billingMode: 'PAY_PER_REQUEST',
});
const CollectionTableAccessPolicy = new aws.iam.Policy('grant-dynamodb-policy', {
  name: `${config.service}-${config.stage}-grant-dynamodb-policy`,
  policy: aws.iam.getPolicyDocument({
    statements: [{
      actions: ['dynamodb:*'],
      resources: ['*'],
    }]
  }).then(result => result.json),
});

new aws.iam.RolePolicyAttachment(`role-policy-cwl`, {
  policyArn: "arn:aws:iam::aws:policy/service-role/AWSAppSyncPushToCloudWatchLogs",
  role: ApiRole.name,
});
new aws.iam.RolePolicyAttachment(`role-policy-ddb`, {
  policyArn: CollectionTableAccessPolicy.arn,
  role: ApiRole.name,
});

const snakeCase = (str: string) => str.replace(/-/g, '_');
const CollectionTableDataSource = new aws.appsync.DataSource(`datasource-collection`, {
  apiId: GraphQLApi.id,
  dynamodbConfig: {
    tableName: CollectionTable.name,
    region: 'ap-northeast-1',
  },
  serviceRoleArn: ApiRole.arn,
  type: 'AMAZON_DYNAMODB',
  name: snakeCase(`${config.service}-${config.stage}-collection`),
});
const GetCollectionResolver = new aws.appsync.Resolver(`appsync-get-collection`, {
  apiId: GraphQLApi.id,
  dataSource: CollectionTableDataSource.name,
  field: 'getCollection',
  requestTemplate: `{
    "version": "2017-02-28",
    "operation": "GetItem",
    "key": {
      "id": { "S": "\${context.arguments.id}" },
      "sort-id": { "S": "detail" }
    }
  }`,
  responseTemplate: `$utils.toJson($context.result)`,
  type: 'Query',
});

const AddCollectionResolver = new aws.appsync.Resolver(`appsync-add-collection`, {
  apiId: GraphQLApi.id,
  dataSource: CollectionTableDataSource.name,
  field: 'addCollection',
  requestTemplate: `{
    "version": "2017-02-28",
    "operation": "PutItem",
    "key": {
      "id": { "S": "\${util.autoId()}" },
      "sort-id": { "S": "detail" }
    },
    "attributeValues": {
      "owner": { "S": "\${context.arguments.owner}" },
      "name": { "S": "\${context.arguments.name}" },
      "title": { "S": "\${context.arguments.title}" },
      "description": { "S": "\${context.arguments.description}" },
      "created_at": { "N": $util.time.nowEpochSeconds() },
      "updated_at": { "N": $util.time.nowEpochSeconds() }
    }
  }`,
  responseTemplate: `$utils.toJson($context.result)`,
  type: 'Mutation',
});

const UpdateCollectionResolver = new aws.appsync.Resolver(`appsync-update-collection`, {
  apiId: GraphQLApi.id,
  dataSource: CollectionTableDataSource.name,
  field: 'updateCollection',
  requestTemplate: `{
    "version": "2017-02-28",
    "operation": "PutItem",
    "key": {
      "id": { "S": "\${context.arguments.name}" },
      "sort-id": { "S": "detail" }
    },
    "attributeValues": {
      "owner": { "S": "\${context.arguments.owner}" },
      "title": { "S": "\${context.arguments.title}" },
      "description": { "S": "\${context.arguments.description}" },
      "updated_at": { "N": $util.time.nowEpochSeconds() }
    }
  }`,
  responseTemplate: `$utils.toJson($context.result)`,
  type: 'Mutation',
});

const UpdateCollectionNameResolver = new aws.appsync.Resolver(`appsync-update-collection-name`, {
  apiId: GraphQLApi.id,
  dataSource: CollectionTableDataSource.name,
  field: 'updateCollectionName',
  requestTemplate: `{
    "version": "2017-02-28",
    "operation": "PutItem",
    "key": {
      "id": { "S": "\${context.arguments.id}" },
      "sort-id": { "S": "detail" }
    },
    "attributeValues": {
      "name": { "S": "\${context.arguments.name}" },
      "updated_at": { "N": $util.time.nowEpochSeconds() }
    }
  }`,
  responseTemplate: `$utils.toJson($context.result)`,
  type: 'Mutation',
});

const DeleteCollectionResolver = new aws.appsync.Resolver(`appsync-delete-collection-name`, {
  apiId: GraphQLApi.id,
  dataSource: CollectionTableDataSource.name,
  field: 'deleteCollection',
  requestTemplate: `{
    "version": "2017-02-28",
    "operation": "DeleteItem",
    "key": {
      "id": { "S": "\${context.arguments.id}" },
      "sort-id": { "S": "detail" }
    }
  }`,
  responseTemplate: `$utils.toJson($context.result)`,
  type: 'Mutation',
});
