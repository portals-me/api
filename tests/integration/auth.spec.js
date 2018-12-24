const axios = require('axios');
const api = require('../../src/api');
const jwt = require('jsonwebtoken');
const fs = require('fs');

describe('Authorizer', () => {
  it('should pass with correct jwt', async () => {
    const testJWT = jwt.sign({
      id: 'user##test-user'
    }, fs.readFileSync('./token/jwtES256.key', 'utf8'), { algorithm: 'ES256' });
    const sdk = api.genSDK('http://localhost:3000', testJWT, axios);

    const result = await sdk.user.me();
    expect(result.status).not.toBe(400);
  });
});
