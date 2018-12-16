const uuid = require('uuid/v4');
const AWS = require('aws-sdk');
const dbc = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event, context) => {
  try {
    const method = event.httpMethod;
    const projectId = (() => {
      const latter = event.path.split('/projects')[1];
      return latter.startsWith('/') ? latter.substring(1) : latter;
    })();
    const user = event.requestContext.authorizer;

    if (method === 'GET') {
      if (projectId === '') {
        const result = await dbc.query({
          TableName: process.env.EntityTable,
          IndexName: 'owner',
          KeyConditionExpression: 'owned_by = :owned_by and begins_with(id, :id)',
          ExpressionAttributeValues: {
            ':owned_by': user.id,
            ':id': 'project',
          }
        }).promise();
        
        return {
          statusCode: 200,
          headers: {
            'Access-Control-Allow-Origin': '*',
          },
          body: JSON.stringify(result.Items),
        };
      }
    }
    
    if (method === 'POST') {
      const project = JSON.parse(event.body);
      const result = await dbc.put({
        TableName: process.env.EntityTable,
        Item: {
          id: `project##${uuid()}`,
          sort: 'detail',
          owned_by: user.id,
          title: project.title,
          description: project.description,
          cover: project.cover,
          media: [],
          comment_members: [user.id],
          comment_count: 0,
          created_at: (new Date()).getTime(),
        }
      }).promise();
      console.log(result);

      return {
        statusCode: 200,
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
        body: null,
      };
    }

    return {
      statusCode: 400,
      headers: {
        'Access-Control-Allow-Origin': '*',
      },
      body: event.body,
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
