import cdk = require('@aws-cdk/cdk');
import s3 = require('@aws-cdk/aws-s3');

const service = 'portals-me';
const stage = 'dev';

class MainStack extends cdk.Stack {
  constructor(parent: cdk.App, props?: cdk.StackProps) {
    super(parent, `${service}-${stage}`, props);

    new s3.Bucket(this, 'StaticFileStore', {
      bucketName: `${service}-web`,
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
