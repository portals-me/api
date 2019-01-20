const fs = require('fs');

const env = {
  "GClientId": "670077302427-0r21asrffhmuhkvfq10qa8kj86cslojn.apps.googleusercontent.com",
  "EntityTable": "portals-me-entities",
  "JwtPrivate": fs.readFileSync('./token/jwtES256.key', 'utf8'),
  "JwtPublic": fs.readFileSync('./token/jwtES256.key.pub', 'utf8')
};

console.log(JSON.stringify(env));
