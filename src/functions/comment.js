const AWS = require('aws-sdk');
const dbc = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event, context) => {
  try {
    const method = event.httpMethod;
    const user = event.requestContext.authorizer;

    if (method === 'POST') {
      const { projectId, message } = JSON.parse(event.body);
      const project = (await dbc.get({
        TableName: process.env.EntityTable,
        Key: {
          id: `project##${projectId}`,
          sort: 'detail'
        }
      }).promise()).Item;

      await dbc.put({
        TableName: process.env.EntityTable,
        Item: {
          id: project.id,
          sort: `comment##${project.comment_count}`,
          owned_by: user.id,
          message: message,
          created_at: (new Date()).getTime(),
        }
      }).promise();

      await dbc.update({
        TableName: process.env.EntityTable,
        Key: {
          id: project.id,
          sort: 'detail',
        },
        UpdateExpression: 'set comment_count = :count',
        ExpressionAttributeValues: {
          ':count': project.comment_count + 1,
        }
      }).promise();

      if (!project.comment_members.includes(user.id)) {
        const members = project.comment_members;
        members.push(user.id);

        await dbc.update({
          TableName: process.env.EntityTable,
          Key: {
            id: project.id,
            sort: 'detail',
          },
          UpdateExpression: 'set comment_members = :members',
          ExpressionAttributeValues: {
            ':members': members,
          }
        }).promise();
      }

      return {
        statusCode: 200,
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
        body: "OK",
      };
    }

    if (method === 'GET') {
      const projectId = event.path.split('/projects/')[1].split('/comments')[0];

      const comments = (await dbc.query({
        TableName: process.env.EntityTable,
        KeyConditionExpression: 'id = :id and begins_with(sort, :sort)',
        ExpressionAttributeValues: {
          ':id': projectId,
          ':sort': 'comment',
        },
      }).promise()).Items;

      return {
        statusCode: 200,
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
        body: JSON.stringify(comments),
      }
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
