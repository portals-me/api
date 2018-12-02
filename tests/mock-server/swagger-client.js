const express = require('express');
const Swagger = require('swagger-client');
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

const main = async () => {
  const client = await Swagger({
    spec: yaml.safeLoad(fs.readFileSync(path.resolve(__dirname, '../openapi.yml'), 'utf8')),
    requestInterceptor: (req) => {
      console.log('request');
      console.log(req);

      throw new Error('error');
    },
  });
  console.log(client.spec);
  console.log(client.apis);

  {
    const res = await client.apis.user.get_users_me();
    console.log(res.body);
  }
};

main();
