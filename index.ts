import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";
import { createLambdaFunction } from "./infrastructure/lambda";

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
