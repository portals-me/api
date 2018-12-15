const uuid = require('uuid/v4');
const AWS = require('aws-sdk');
const idp = new AWS.CognitoIdentity();
const dbc = new AWS.DynamoDB.DocumentClient();
const { OAuth2Client } = require('google-auth-library');
const gclient = new OAuth2Client(process.env.GClientId);
const jsonwebtoken = require('jsonwebtoken');

const gverify = async (token) => {
  const ticket = await gclient.verifyIdToken({
    idToken: token,
    audience: process.env.GClientId,
  });
  return ticket.getPayload();
};

exports.signUp = async (event, context) => {
  const idp_id = (await idp.getId({
    IdentityPoolId: 'ap-northeast-1:5221828e-b1d8-45e7-9361-b06057573aa9',
    Logins: {
      'accounts.google.com': event.body,
    },
  }).promise()).IdentityId;

  // Is the verification necessary?
  const gaccount = await gverify(event.body);

  await dbc.put({
    TableName: process.env.EntityTable,
    Item: {
      id: `user##${idp_id}`,
      sort: 'detail',
      created_at: (new Date()).getTime(),
      name: gaccount.name,
      display_name: gaccount.name,
      picture: gaccount.picture,
    },
  }).promise();

  return {
    statusCode: 200,
    body: 'OK',
  };
};

exports.signIn = async (event, context) => {
  const idp_id = (await idp.getId({
    IdentityPoolId: 'ap-northeast-1:5221828e-b1d8-45e7-9361-b06057573aa9',
    Logins: {
      'accounts.google.com': event.body,
    },
  }).promise()).IdentityId;
  console.log(idp_id);

  const user = (await dbc.get({
    TableName: process.env.EntityTable,
    Key: {
      id: `user##${idp_id}`,
      sort: 'detail',
    }
  }).promise()).Item;
  console.log(user);

  const token = jsonwebtoken.sign({
    id: idp_id,
    name: user.name,
    display_name: user.display_name,
    picture: user.picture,
    created_at: user.created_at,
  }, process.env.JwtPrivate, { algorithm: 'ES256' });

  return {
    statusCode: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
    },
    body: JSON.stringify({
      id_token: token,
      user: user,
    }),
  };
};
