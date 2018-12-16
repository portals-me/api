const AWS = require('aws-sdk');
const dbc = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event, context) => {
  try {
    const method = event.httpMethod;
    const name = event.path.startsWith('/') ? event.path.substring(1) : event.path;

    if (method === 'GET') {
      if (event.path === '/') {
        const result = await dbc.query({
          TableName: process.env.EntityTable,
          IndexName: 'owner',
          KeyConditionExpression: 'owned_by = :owned_by and begins_with(id, :id)',
          ExpressionAttributeValues: {
            ':owned_by': userId,
            ':id': 'project',
          }
        });
        
        return {
          statusCode: 200,
          headers: {
            'Access-Control-Allow-Origin': '*',
          },
          body: JSON.stringify(result.Items),
        };
      }

      if (name) {
      }
    }

    return {
      statusCode: 200,
      headers: {
        'Access-Control-Allow-Origin': '*',
      },
      body: body,
    };
  } catch(error) {
    const body = error.stack || JSON.stringify(error, null, 2);

    return {
      statusCode: 400,
      headers: {
        'Access-Control-Allow-Origin': '*',
      },
      body: body,
    };
  }
};
