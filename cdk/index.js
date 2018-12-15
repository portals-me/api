const cdk = require('@aws-cdk/cdk');
const s3 = require('@aws-cdk/aws-s3');
const lambda = require('@aws-cdk/aws-lambda');
const apigateway = require('@aws-cdk/aws-apigateway');
const iam = require('@aws-cdk/aws-iam');

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

    const signInHandler = new lambda.Function(this, 'SignInHandler', {
      runtime: lambda.Runtime.NodeJS810,
      code: lambda.Code.directory('src/functions'),
      handler: 'auth.signIn',
      role: role,
    });

    const api = new apigateway.RestApi(this, 'RestApi', {
      restApiName: `${service}-${stage}`,
    });

    api.root.addResource('signIn').addMethod('POST', new apigateway.LambdaIntegration(signInHandler));
  }
}

class App extends cdk.App {
  constructor () {
    super();
    new MainStack(this);
  }
}

new App().run();
