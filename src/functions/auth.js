const AWS = require('aws-sdk');
const idp = new AWS.CognitoIdentity();

exports.signIn = async (event, context) => {
  const idp_id = await idp.getId({
    IdentityPoolId: 'ap-northeast-1:5221828e-b1d8-45e7-9361-b06057573aa9',
    Logins: {
      'accounts.google.com': event.body,
    },
  }).promise();
  const cred = await idp.getCredentialsForIdentity({
    IdentityId: idp_id.IdentityId,
    Logins: {
      'accounts.google.com': event.body,
    },
  }).promise();
  console.log(cred);

  return {
    statusCode: 200,
    body: 'OK',
  };
};
