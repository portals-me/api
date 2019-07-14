import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import { createLambdaFunction, LambdaOptions } from "./lambda";

export const createDataSource = (
  name: string,
  options: {
    api: aws.appsync.GraphQLApi;
    dataSourceName: string;
    dataSource: Omit<
      aws.appsync.DataSourceArgs,
      "apiId" | "name" | "serviceRoleArn"
    >;
    rolePolicyDocument: pulumi.Output<string>;
  }
) => {
  const appsyncRole = new aws.iam.Role(`${name}-ds-role`, {
    assumeRolePolicy: `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "appsync.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
  }
  `
  });

  new aws.iam.RolePolicy(`${name}-ds-role-policy`, {
    policy: options.rolePolicyDocument,
    role: appsyncRole.id
  });

  return new aws.appsync.DataSource(
    `${name}-ds`,
    {
      apiId: options.api.id,
      name: options.dataSourceName,
      serviceRoleArn: appsyncRole.arn,
      ...options.dataSource
    },
    {
      dependsOn: [options.api, appsyncRole]
    }
  );
};

export const createLambdaDataSource = (
  name: string,
  options: {
    api: aws.appsync.GraphQLApi;
    dataSourceName: string;
    function: aws.lambda.Function;
    dataSource?: Omit<
      aws.appsync.DataSourceArgs,
      "apiId" | "lambdaConfig" | "serviceRoleArn" | "type" | "name"
    >;
  }
) => {
  return createDataSource(name, {
    api: options.api,
    dataSourceName: options.dataSourceName,
    dataSource: {
      lambdaConfig: {
        functionArn: options.function.arn
      },
      type: "AWS_LAMBDA"
    },
    rolePolicyDocument: pulumi.interpolate`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "lambda:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "${options.function.arn}"
      ]
    }
  ]
}
`
  });
};

export const createDynamoDBDataSource = (
  name: string,
  options: {
    api: aws.appsync.GraphQLApi;
    dataSourceName: string;
    table: aws.dynamodb.Table;
    dataSource?: Omit<
      aws.appsync.DataSourceArgs,
      "apiId" | "dynamodbConfig" | "serviceRoleArn" | "type" | "name"
    >;
  }
) => {
  return createDataSource(name, {
    api: options.api,
    dataSourceName: options.dataSourceName,
    dataSource: {
      dynamodbConfig: {
        tableName: options.table.name
      },
      type: "AMAZON_DYNAMODB"
    },
    rolePolicyDocument: pulumi.interpolate`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "dynamodb:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "${options.table.arn}",
        "${options.table.arn}/index/*"
      ]
    }
  ]
}
`
  });
};

export const createPipelineResolver = (
  name: string,
  options: {
    api: aws.appsync.GraphQLApi;
    field: string;
    type: "Query" | "Mutation";
    pipeline: Array<aws.appsync.Function>;
  }
) => {
  return new aws.appsync.Resolver(name, {
    apiId: options.api.id,
    field: options.field,
    type: options.type,
    requestTemplate: "$util.toJson($context)",
    responseTemplate: "$utils.toJson($context.prev.result)",
    kind: "PIPELINE",
    pipelineConfig: {
      functions: options.pipeline.map(f => f.functionId)
    }
  });
};

export const createLambdaResolverFunction = (
  name: string,
  options: { lambda: LambdaOptions; api: aws.appsync.GraphQLApi }
) => {
  const lambda = createLambdaFunction(name, options.lambda);
  const datasource = createLambdaDataSource(name, {
    api: options.api,
    function: lambda,
    dataSourceName: name.split("-").join("")
  });

  return new aws.appsync.Function(name, {
    apiId: options.api.id,
    dataSource: datasource.name,
    requestMappingTemplate: `{
  "version": "2017-02-28",
  "operation": "Invoke",
  "payload": $utils.toJson($context)
}`,
    responseMappingTemplate: `#if($context.error)
  $util.error($context.error.type, $context.error.message)
#end

$util.toJson($context.result)`,
    name: name.split("-").join("")
  });
};
