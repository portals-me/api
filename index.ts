import * as fs from "fs";
import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";
import { createLambdaFunction } from "./infrastructure/lambda";
import {
  createLambdaDataSource,
  createDynamoDBDataSource
} from "./infrastructure/appsync";

const config = {
  service: new pulumi.Config().name,
  stage: pulumi.getStack()
};

const parameter = {
  jwtPrivate: aws.ssm
    .getParameter({
      name: `${config.service}-${config.stage}-jwt-private`,
      withDecryption: true
    })
    .then(result => result.value)
};

const lambdaRole = (() => {
  const role = new aws.iam.Role("lambda-role", {
    assumeRolePolicy: aws.iam
      .getPolicyDocument({
        statements: [
          {
            actions: ["sts:AssumeRole"],
            principals: [
              {
                identifiers: ["lambda.amazonaws.com"],
                type: "Service"
              }
            ]
          }
        ]
      })
      .then(result => result.json)
  });
  new aws.iam.RolePolicyAttachment("lambda-role-lambdafull", {
    role: role,
    policyArn: aws.iam.AWSLambdaFullAccess
  });

  return role;
})();

const authorizer = createLambdaFunction("authorizer", {
  filepath: "authorizer",
  handlerName: `${config.service}-${config.stage}-authorizer`,
  role: lambdaRole,
  lambdaOptions: {
    environment: {
      variables: {
        jwtPrivate: parameter.jwtPrivate
      }
    }
  }
});

const postTable = new aws.dynamodb.Table("post", {
  billingMode: "PAY_PER_REQUEST",
  name: `${config.service}-${config.stage}-post`,
  attributes: [
    {
      name: "id",
      type: "S"
    },
    {
      name: "sort",
      type: "S"
    },
    {
      name: "updated_at",
      type: "N"
    },
    {
      name: "owner",
      type: "S"
    }
  ],
  hashKey: "id",
  rangeKey: "sort",
  globalSecondaryIndexes: [
    {
      name: "owner",
      hashKey: "owner",
      rangeKey: "updated_at",
      projectionType: "ALL"
    }
  ]
});

const accountReplicaTable = new aws.dynamodb.Table("account-replica", {
  billingMode: "PAY_PER_REQUEST",
  name: `${config.service}-${config.stage}-account-replica`,
  attributes: [
    {
      name: "id",
      type: "S"
    },
    {
      name: "sort",
      type: "S"
    },
    {
      name: "name",
      type: "S"
    }
  ],
  hashKey: "id",
  rangeKey: "sort",
  globalSecondaryIndexes: [
    {
      name: "name",
      hashKey: "name",
      projectionType: "ALL"
    }
  ]
});

const putAccountEvent = createLambdaFunction("put-account-table", {
  filepath: "put-account-table",
  handlerName: `${config.service}-${config.stage}-put-account-table`,
  role: lambdaRole,
  lambdaOptions: {
    environment: {
      variables: {
        accountTableName: accountReplicaTable.name
      }
    }
  }
});

const putAccountEventSubscription = new aws.sns.TopicSubscription(
  "put-account-table-event-subscription",
  {
    protocol: "lambda",
    endpoint: putAccountEvent.arn,
    topic: "arn:aws:sns:ap-northeast-1:941528793676:portals-me-account-dev-account-table-event-topic" as any
  }
);

const putAccountEventPermission = new aws.lambda.Permission(
  "put-account-table-event-permission",
  {
    function: putAccountEvent.name,
    action: "lambda:InvokeFunction",
    principal: "sns.amazonaws.com"
  }
);

const graphqlApi = new aws.appsync.GraphQLApi("graphql-api", {
  name: `${config.service}-${config.stage}-api`,
  authenticationType: "API_KEY",
  schema: fs.readFileSync("./schema.graphql").toString(),
  logConfig: {
    cloudwatchLogsRoleArn:
      "arn:aws:iam::941528793676:role/AppSyncToCWLServiceRole",
    fieldLogLevel: "ALL"
  }
});

// API Key expires in one year
const graphqlApiKey = new aws.appsync.ApiKey(
  "graphql-api-key",
  {
    apiId: graphqlApi.id,
    expires: new Date(
      new Date().getTime() + 1000 * 60 * 60 * 24 * 365
    ).toISOString()
  },
  {
    dependsOn: [graphqlApi]
  }
);

const authorizerDS = createLambdaDataSource("authorizer", {
  api: graphqlApi,
  function: authorizer,
  dataSourceName: "authorizer"
});
const postDS = createDynamoDBDataSource("post", {
  api: graphqlApi,
  table: postTable,
  dataSourceName: "post"
});
const accountDS = createDynamoDBDataSource("account", {
  api: graphqlApi,
  table: accountReplicaTable,
  dataSourceName: "account"
});

const getUserByName = new aws.appsync.Resolver("getUserByName", {
  apiId: graphqlApi.id,
  dataSource: accountDS.name,
  field: "getUserByName",
  type: "Query",
  requestTemplate: fs.readFileSync("./vtl/user/GetUserByName.vtl").toString(),
  responseTemplate: fs
    .readFileSync("./vtl/user/GetUserByNameResponse.vtl")
    .toString()
});

const authorizerFunctionResolver = new aws.appsync.Function(
  "authorizer-function",
  {
    apiId: graphqlApi.id,
    dataSource: authorizerDS.name,
    requestMappingTemplate: fs
      .readFileSync("./vtl/AuthorizerRequest.vtl")
      .toString(),
    responseMappingTemplate: fs
      .readFileSync("./vtl/AuthorizerResponse.vtl")
      .toString(),
    name: "authorizer"
  }
);

const replaceOwner = createLambdaFunction("replace-owner", {
  filepath: "replace-owner",
  handlerName: `${config.service}-${config.stage}-replace-owner`,
  role: lambdaRole,
  lambdaOptions: {
    environment: {
      variables: {
        accountTableName: accountReplicaTable.name
      }
    }
  }
});

const replaceOwnerDS = createLambdaDataSource("replace-owner", {
  api: graphqlApi,
  function: replaceOwner,
  dataSourceName: "replaceOwner"
});

const replaceOwnerFunction = new aws.appsync.Function(
  "replace-owner-function",
  {
    apiId: graphqlApi.id,
    dataSource: replaceOwnerDS.name,
    requestMappingTemplate: fs
      .readFileSync("./vtl/ReplaceOwnerRequest.vtl")
      .toString(),
    responseMappingTemplate: fs
      .readFileSync("./vtl/ReplaceOwnerResponse.vtl")
      .toString(),
    name: "replaceOwner"
  }
);

const listPostSummary = (() => {
  const listPostSummaryFunction = new aws.appsync.Function("listPostSummary", {
    apiId: graphqlApi.id,
    dataSource: postDS.name,
    requestMappingTemplate: fs
      .readFileSync("./vtl/post/ListPostSummary.vtl")
      .toString(),
    responseMappingTemplate: fs
      .readFileSync("./vtl/post/PostSummaryItems.vtl")
      .toString(),
    name: "listPostSummary"
  });

  return new aws.appsync.Resolver(
    "listPostSummary",
    {
      apiId: graphqlApi.id,
      field: "listPostSummary",
      type: "Query",
      requestTemplate: fs.readFileSync("./vtl/ContextRequest.vtl").toString(),
      responseTemplate: fs.readFileSync("./vtl/PrevResult.vtl").toString(),
      kind: "PIPELINE",
      pipelineConfig: {
        functions: [
          authorizerFunctionResolver.functionId,
          listPostSummaryFunction.functionId,
          replaceOwnerFunction.functionId
        ]
      }
    },
    {
      dependsOn: [
        authorizerFunctionResolver,
        listPostSummaryFunction,
        replaceOwnerFunction
      ]
    }
  );
})();

const addSharePost = (() => {
  const addSharePostFunction = new aws.appsync.Function("addSharePost", {
    apiId: graphqlApi.id,
    dataSource: postDS.name,
    requestMappingTemplate: fs
      .readFileSync("./vtl/post/AddSharePost.vtl")
      .toString(),
    responseMappingTemplate: fs
      .readFileSync("./vtl/post/PostSummary.vtl")
      .toString()
  });

  return new aws.appsync.Resolver(
    "addSharePost",
    {
      apiId: graphqlApi.id,
      field: "addSharePost",
      type: "Mutation",
      requestTemplate: fs.readFileSync("./vtl/ContextRequest.vtl").toString(),
      responseTemplate: fs.readFileSync("./vtl/PrevResult.vtl").toString(),
      kind: "PIPELINE",
      pipelineConfig: {
        functions: [
          authorizerFunctionResolver.functionId,
          addSharePostFunction.functionId
        ]
      }
    },
    {
      dependsOn: [authorizerFunctionResolver, addSharePostFunction]
    }
  );
})();

const userStorage = new aws.s3.Bucket("user-storage", {
  bucketPrefix: `${config.service}-${config.stage}-user-storage`
});

const generateUploadURL = (() => {
  const generateUploadURLDataSource = createLambdaDataSource(
    "generate-upload-url-ds",
    {
      api: graphqlApi,
      function: createLambdaFunction("generate-upload-url", {
        filepath: "generate-upload-url",
        handlerName: `${config.service}-${config.stage}-generate-upload-url`,
        role: lambdaRole,
        lambdaOptions: {
          environment: {
            variables: {
              storageBucket: userStorage.bucket
            }
          }
        }
      }),
      dataSourceName: "generateUploadURL"
    }
  );

  const generateUploadURLFunction = new aws.appsync.Function(
    "generateUploadURL",
    {
      apiId: graphqlApi.id,
      dataSource: generateUploadURLDataSource.name,
      requestMappingTemplate: fs
        .readFileSync("./vtl/AuthorizerRequest.vtl")
        .toString(),
      responseMappingTemplate: fs
        .readFileSync("./vtl/AuthorizerResponse.vtl")
        .toString(),
      name: "generateUploadURL"
    }
  );

  return new aws.appsync.Resolver(
    "generateUploadURL",
    {
      apiId: graphqlApi.id,
      field: "generateUploadURL",
      type: "Mutation",
      requestTemplate: fs.readFileSync("./vtl/ContextRequest.vtl").toString(),
      responseTemplate: fs.readFileSync("./vtl/PrevResult.vtl").toString(),
      kind: "PIPELINE",
      pipelineConfig: {
        functions: [
          authorizerFunctionResolver.functionId,
          generateUploadURLFunction.functionId
        ]
      }
    },
    {
      dependsOn: [authorizerFunctionResolver, generateUploadURLFunction]
    }
  );
})();

export const output = {
  appsync: {
    url: graphqlApi.uris["GRAPHQL"],
    apiKey: graphqlApiKey.key
  }
};
