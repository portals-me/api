const fs = require('fs');

const env = {
  "GClientId": "670077302427-0r21asrffhmuhkvfq10qa8kj86cslojn.apps.googleusercontent.com",
  "EntityTable": "portals-me-entities",
  "IdentityPoolId": "ap-northeast-1:5221828e-b1d8-45e7-9361-b06057573aa9",
  "JwtPrivate": fs.readFileSync('./token/jwtES256.key', 'utf8'),
  "JwtPublic": fs.readFileSync('./token/jwtES256.key.pub', 'utf8'),
  "TwitterKey": fs.readFileSync('./token/twitter.key', 'utf8'),
};

console.log(JSON.stringify(env));
