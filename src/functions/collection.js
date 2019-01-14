const AWS = process.env.IS_OFFLINE === 'true' ? require('aws-sdk') : require('aws-xray-sdk').captureAWS(require('aws-sdk'));
const uuid = require('uuid/v4');
const dbc = process.env.IS_OFFLINE === 'true'
  ? new AWS.DynamoDB.DocumentClient({ region: 'localhost', endpoint: `http://localhost:${process.env.TestPort}` })
  : new AWS.DynamoDB.DocumentClient();

exports.handler = async (event, context) => {
  try {
    const method = event.httpMethod;
    const collectionId = event.pathParameters ? event.pathParameters.collectionId : null;
    const user = event.requestContext.authorizer;

    if (method === 'GET') {
      if (!collectionId) {
        const result = await dbc.query({
          TableName: process.env.EntityTable,
          IndexName: 'owner',
          KeyConditionExpression: 'owned_by = :owned_by and begins_with(id, :id)',
          ExpressionAttributeValues: {
            ':owned_by': user.id,
            ':id': 'collection',
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
        const collection = (await dbc.get({
          TableName: process.env.EntityTable,
          Key: {
            id: `collection##${collectionId}`,
            sort: 'detail',
          }
        }).promise()).Item;

        const comments = (await dbc.query({
          TableName: process.env.EntityTable,
          KeyConditionExpression: 'id = :id and begins_with(sort, :sort)',
          ExpressionAttributeValues: {
            ':id': collection.id,
            ':sort': 'comment',
          },
        }).promise()).Items;

        const members = (await Promise.all(
          collection.comment_members.map(async (memberId) => {
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
          body: JSON.stringify(Object.assign(collection, { comments, members })),
        };
      }
    }
    
    if (method === 'POST') {
      const collection = JSON.parse(event.body);
      const collectionId = uuid();
      await dbc.put({
        TableName: process.env.EntityTable,
        Item: {
          id: `collection##${collectionId}`,
          sort: 'detail',
          owned_by: user.id,
          title: collection.title,
          description: collection.description,
          cover: collection.cover,
          media: [],
          comment_members: [user.id],
          comment_count: 0,
          created_at: (new Date()).getTime(),
        }
      }).promise();

      return {
        statusCode: 201,
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Location': `/collections/${collectionId}`,
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
