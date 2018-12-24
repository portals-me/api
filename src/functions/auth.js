const AWS = process.env.IS_OFFLINE === 'true' ? require('aws-sdk') : require('aws-xray-sdk').captureAWS(require('aws-sdk'));
const idp = new AWS.CognitoIdentity();
const dbc = new AWS.DynamoDB.DocumentClient();
const { OAuth2Client } = require('google-auth-library');
const gclient = new OAuth2Client(process.env.GClientId);
const jsonwebtoken = require('jsonwebtoken');

const mod = (module.exports = {});

mod.gverify = async (token) => {
  const ticket = await gclient.verifyIdToken({
    idToken: token,
    audience: process.env.GClientId,
  });
  return ticket.getPayload();
};

mod.signUp = async (event, context) => {
  const { google_token, name, display_name, picture } = JSON.parse(event.body);

  const idp_id = (await idp.getId({
    IdentityPoolId: 'ap-northeast-1:5221828e-b1d8-45e7-9361-b06057573aa9',
    Logins: {
      'accounts.google.com': google_token,
    },
  }).promise()).IdentityId;

  // Is the verification necessary?
  const gaccount = await mod.gverify(google_token);

  await dbc.put({
    TableName: process.env.EntityTable,
    Item: {
      id: `user##${idp_id}`,
      sort: 'detail',
      created_at: (new Date()).getTime(),
      name: name || gaccount.name,
      display_name: display_name || gaccount.name,
      picture: picture || gaccount.picture,
    },
  }).promise();

  return {
    statusCode: 201,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Location': `/users/${idp_id}`,
    },
    body: null,
  };
};

mod.signIn = async (event, context) => {
  const idp_id = (await idp.getId({
    IdentityPoolId: 'ap-northeast-1:5221828e-b1d8-45e7-9361-b06057573aa9',
    Logins: {
      'accounts.google.com': event.body,
    },
  }).promise()).IdentityId;
  const userId = `user##${idp_id}`;

  const user = (await dbc.get({
    TableName: process.env.EntityTable,
    Key: {
      id: userId,
      sort: 'detail',
    }
  }).promise()).Item;

  if (!user) {
    return {
      statusCode: 404,
      headers: {
        'Access-Control-Allow-Origin': '*',
      },
      body: 'Specified user does not exist...',
    }
  }

  const token = jsonwebtoken.sign({
    id: userId,
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

mod.authorize = async (event, context, callback) => {
  let generatePolicy = (principalId, effect, resource, context) => ({
    principalId: principalId,
    policyDocument: {
      Version: '2012-10-17',
      Statement: [
        {
          Action: 'execute-api:Invoke',
          Effect: effect,
          Resource: resource,
        }
      ]
    },
    context: context,
  });

  try {
    console.log(event);
    const token = event.authorizationToken.split('Bearer ')[1];
    const methodArn = event.methodArn;
  
    if (!token) {
      throw new Error('Unauthorized');
    }
  
    const decoded = jsonwebtoken.verify(token, process.env.JwtPublic, { algorithm: 'ES256' });

    // skip scope check now
    const isAllowed = true;
    const effect = isAllowed ? 'Allow' : 'Deny';
    const userId = decoded.id;
    const authorizationContext = decoded;
    console.log(generatePolicy(userId, effect, methodArn, authorizationContext));

    callback(null, generatePolicy(userId, effect, methodArn, authorizationContext));
  } catch (err) {
    callback(err.message, generatePolicy('user', 'Deny', event.methodArn, null));
  }
};

mod.getMe = async (event, context) => {
  const token = event.headers.Authorization.split('Bearer ')[1];
  const decoded = jsonwebtoken.verify(token, process.env.JwtPublic, { algorithm: 'ES256' });

  return {
    statusCode: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
    },
    body: JSON.stringify(decoded),
  };
};
