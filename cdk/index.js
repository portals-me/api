const cdk = require('@aws-cdk/cdk');
const s3 = require('@aws-cdk/aws-s3');
const lambda = require('@aws-cdk/aws-lambda');
const apigateway = require('@aws-cdk/aws-apigateway');
const iam = require('@aws-cdk/aws-iam');
const dynamodb = require('@aws-cdk/aws-dynamodb');
const fs = require('fs');

const service = 'portals-me';
const stage = 'dev';

/**
 * @param api {apigateway.IRestApiResource}
 */
function addCorsOptions(api) {
  const options = api.addMethod('OPTIONS', new apigateway.MockIntegration({
    integrationResponses: [{
      statusCode: '200',
      responseParameters: {
        'method.response.header.Access-Control-Allow-Headers': "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent'",
        'method.response.header.Access-Control-Allow-Origin': "'*'",
        'method.response.header.Access-Control-Allow-Credentials': "'false'",
        'method.response.header.Access-Control-Allow-Methods': "'OPTIONS,GET,POST,PUT,DELETE'",
      },
    }],
    passthroughBehavior: apigateway.PassthroughBehavior.Never,
    requestTemplates: {
      'application/json': '{"statusCode": 200}'
    },
  }));
  const methodResource = /** @type {apigateway.cloudformation.MethodResource} */ (options.findChild('Resource'));
  methodResource.propertyOverrides.methodResponses = [{
    statusCode: '200',
    responseModels: {
      'application/json': 'Empty',
    },
    responseParameters: {
      'method.response.header.Access-Control-Allow-Headers': true,
      'method.response.header.Access-Control-Allow-Origin': true,
      'method.response.header.Access-Control-Allow-Credentials': true,
      'method.response.header.Access-Control-Allow-Methods': true,
    },
  }];

  return api;
}

class MainStack extends cdk.Stack {
  constructor(parent, props) {
    super(parent, `${service}-${stage}`, props);

    new s3.Bucket(this, 'Hosting', {
      bucketName: `${service}`,
    });

    const role = new iam.Role(this, 'cognitoRoleForLambda', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
    });
    role.addToPolicy(
      new iam.PolicyStatement()
        .addAllResources()
        .addActions([
          'logs:CreateLogGroup',
          'logs:CreateLogStream',
          'logs:PutLogEvents',
          'cognito-identity:*',
        ])
    );

    const entityTable = new dynamodb.Table(this, 'EntityTable', {
      partitionKey: {
        name: 'id',
        type: dynamodb.AttributeType.String
      },
      sortKey: {
        name: 'sort',
        type: dynamodb.AttributeType.String
      },
    });
    entityTable.grantFullAccess(role);

    const entityTableResource = /** @type {dynamodb.cloudformation.TableResource} */ (entityTable.findChild('Resource'));
    entityTableResource.propertyOverrides.billingMode = 'PAY_PER_REQUEST';
    delete entityTableResource.properties.provisionedThroughput;

    entityTableResource.propertyOverrides.globalSecondaryIndexes = [
      {
        indexName: 'owner',
        keySchema: [
          {
            attributeName: 'owned_by',
            keyType: 'HASH',
          },
          {
            attributeName: 'id',
            keyType: 'RANGE'
          }
        ],
        projection: {
          projectionType: dynamodb.ProjectionType.All
        },
      }
    ];
    entityTableResource.properties.attributeDefinitions.push({
      attributeName: 'owned_by',
      attributeType: dynamodb.AttributeType.String,
    });
    role.addToPolicy(
      new iam.PolicyStatement()
        .addResource(entityTable.tableArn + "/index/*")
        .addAction('dynamodb:*')
    );

    const api = new apigateway.RestApi(this, 'RestApi', {
      restApiName: `${service}-${stage}`,
      deployOptions: {
        stageName: 'dev',
        dataTraceEnabled: true,
        loggingLevel: apigateway.MethodLoggingLevel.Error,
      }
    });

    const signUpHandler = new lambda.Function(this, 'SignUpHandler', {
      runtime: lambda.Runtime.NodeJS810,
      code: lambda.Code.directory('src/functions'),
      handler: 'auth.signUp',
      role: role,
      environment: {
        GClientId: '670077302427-0r21asrffhmuhkvfq10qa8kj86cslojn.apps.googleusercontent.com',
        EntityTable: entityTable.findChild('Resource').ref,
      },
    });
    addCorsOptions(api.root.addResource('signUp')).addMethod('POST', new apigateway.LambdaIntegration(signUpHandler));

    const signInHandler = new lambda.Function(this, 'SignInHandler', {
      runtime: lambda.Runtime.NodeJS810,
      code: lambda.Code.directory('src/functions'),
      handler: 'auth.signIn',
      role: role,
      environment: {
        EntityTable: entityTable.findChild('Resource').ref,
        JwtPrivate: fs.readFileSync('./token/jwtES256.key', 'utf8'),
      },
    });
    addCorsOptions(api.root.addResource('signIn')).addMethod('POST', new apigateway.LambdaIntegration(signInHandler));

    const authorizerHandler = new lambda.Function(this, 'Authorizer', {
      runtime: lambda.Runtime.NodeJS810,
      code: lambda.Code.directory('src/functions'),
      handler: 'auth.authorize',
      role: role,
      environment: {
        JwtPublic: fs.readFileSync('./token/jwtES256.key.pub', 'utf8'),
      },
    });

    authorizerHandler.addPermission('AuthorizerPermission', {
      action: 'lambda:InvokeFunction',
      principal: new iam.ServicePrincipal('apigateway.amazonaws.com'),
    });

    const authorizerResource = new apigateway.cloudformation.AuthorizerResource(this, 'LambdaTokenAuthorizer', {
      restApiId: api.restApiId,
      identitySource: 'method.request.header.Authorization',
      type: 'TOKEN',
      authorizerUri: `arn:aws:apigateway:ap-northeast-1:lambda:path/2015-03-31/functions/${authorizerHandler.functionArn}/invocations`,
      name: 'LambdaTokenAuthorizer',
    });

    const projectHandler = new lambda.Function(this, 'ProjectHandler', {
      runtime: lambda.Runtime.NodeJS810,
      code: lambda.Code.directory('src/functions'),
      handler: 'project.handler',
      role: role,
      environment: {
        EntityTable: entityTable.findChild('Resource').ref,
      },
    });

    const apiProject = addCorsOptions(api.root.addResource('projects'));
    apiProject
      .addMethod('GET', new apigateway.LambdaIntegration(projectHandler), {
        authorizationType: apigateway.AuthorizationType.Custom,
        authorizerId: authorizerResource.ref,
      });
    apiProject
      .addMethod('POST', new apigateway.LambdaIntegration(projectHandler), {
        authorizationType: apigateway.AuthorizationType.Custom,
        authorizerId: authorizerResource.ref,
      });
  }
}

class App extends cdk.App {
  constructor () {
    super();
    new MainStack(this);
  }
}

new App().run();
