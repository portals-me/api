const AWS = require('aws-xray-sdk').captureAWS(require('aws-sdk'));
const uuid = require('uuid/v4');
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
      } else {
        const project = (await dbc.get({
          TableName: process.env.EntityTable,
          Key: {
            id: `project##${projectId}`,
            sort: 'detail',
          }
        }).promise()).Item;

        const comments = (await dbc.query({
          TableName: process.env.EntityTable,
          KeyConditionExpression: 'id = :id and begins_with(sort, :sort)',
          ExpressionAttributeValues: {
            ':id': project.id,
            ':sort': 'comment',
          },
        }).promise()).Items;

        const members = (await Promise.all(
          project.comment_members.map(async (memberId) => {
            return (await dbc.get({
              TableName: process.env.EntityTable,
              Key: { id: memberId, sort: 'detail' },
            }).promise()).Item;
          })
        )).reduce((obj, item) => {
          obj[item.id] = item;
          return obj;
        }, {});
  
        return {
          statusCode: 200,
          headers: {
            'Access-Control-Allow-Origin': '*',
          },
          body: JSON.stringify(Object.assign(project, { comments, members })),
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
