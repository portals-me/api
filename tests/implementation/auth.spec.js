const AWS = require('aws-sdk-mock');
const fs = require('fs');

process.env.IS_OFFLINE = 'true';
process.env.JwtPrivate = fs.readFileSync('token/jwtES256.key', 'utf8');

AWS.mock('CognitoIdentity', 'getId');
AWS.mock('DynamoDB.DocumentClient', 'put');
AWS.mock('DynamoDB.DocumentClient', 'get');

describe('Auth', () => {
  it('can signUp with correct token', async () => {
    AWS.remock('CognitoIdentity', 'getId', (params, callback) => {
      callback(null, {
        IdentityId: 'test-id'
      });
    });
    AWS.remock('DynamoDB.DocumentClient', 'put', (params, callback) => {
      callback(null, null);
    });

    const auth = require('../../src/functions/auth');
    auth.gverify = (token) => ({
      name: 'name',
      picture: 'picture',
    });

    const result = await auth.signUp({
      body: '',
    });

    expect(result.statusCode).toBe(201);
  });

  it('can signIn if account exists', async () => {
    AWS.remock('CognitoIdentity', 'getId', (params, callback) => {
      callback(null, {
        IdentityId: 'test-id'
      });
    });
    AWS.remock('DynamoDB.DocumentClient', 'get', (params, callback) => {
      callback(null, {
        Item: {
          name: 'name',
          display_name: 'Name',
          picture: 'picture',
          created_at: 10000,
        }
      });
    });

    const auth = require('../../src/functions/auth');

    const result = await auth.signIn({
      body: ''
    });

    expect(result.statusCode).toBe(200);
  });

  it('cannot signIn if account does not exist', async () => {
    AWS.remock('CognitoIdentity', 'getId', (params, callback) => {
      callback(null, {
        IdentityId: 'test-id'
      });
    });
    AWS.remock('DynamoDB.DocumentClient', 'get', (params, callback) => {
      callback(null, {});
    });

    const auth = require('../../src/functions/auth');

    const result = await auth.signIn({
      body: ''
    });

    expect(result.statusCode).toBe(400);
  });
});
