const cdk = require('@aws-cdk/cdk');
const s3 = require('@aws-cdk/aws-s3');
const lambda = require('@aws-cdk/aws-lambda');
const apigateway = require('@aws-cdk/aws-apigateway');
const iam = require('@aws-cdk/aws-iam');
const dynamodb = require('@aws-cdk/aws-dynamodb');

const service = 'portals-me';
const stage = 'dev';

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

    const entityTableResource = /** @type {dynamodb.cloudformation.TableResource} */ (entityTable.findChild('Resource'));
    entityTableResource.propertyOverrides.billingMode = 'PAY_PER_REQUEST';
    delete entityTableResource.properties.provisionedThroughput;

    const api = new apigateway.RestApi(this, 'RestApi', {
      restApiName: `${service}-${stage}`,
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
    api.root.addResource('signUp').addMethod('POST', new apigateway.LambdaIntegration(signUpHandler));

    const signInHandler = new lambda.Function(this, 'SignInHandler', {
      runtime: lambda.Runtime.NodeJS810,
      code: lambda.Code.directory('src/functions'),
      handler: 'auth.signIn',
      role: role,
      environment: {
        EntityTable: entityTable.findChild('Resource').ref,
      },
    });
    api.root.addResource('signIn').addMethod('POST', new apigateway.LambdaIntegration(signInHandler));

    entityTable.grantFullAccess(role);
  }
}

class App extends cdk.App {
  constructor () {
    super();
    new MainStack(this);
  }
}

new App().run();
