const AWS = process.env.IS_OFFLINE === 'true'
  ? require('aws-sdk')
  : require('aws-xray-sdk').captureAWS(require('aws-sdk'));
const dbc = process.env.IS_OFFLINE === 'true'
  ? new AWS.DynamoDB.DocumentClient({ region: 'localhost', endpoint: `http://localhost:${process.env.TestPort}` })
  : new AWS.DynamoDB.DocumentClient();

exports.handler = async (event, context) => {
  try {
    const method = event.httpMethod;
    const user = event.requestContext.authorizer;

    if (method === 'POST') {
      const { collectionId, message } = JSON.parse(event.body);
      const collection = (await dbc.get({
        TableName: process.env.EntityTable,
        Key: {
          id: `collection##${collectionId}`,
          sort: 'detail'
        }
      }).promise()).Item;

      const commentIndex = collection.comment_count;
      await dbc.put({
        TableName: process.env.EntityTable,
        Item: {
          id: collection.id,
          sort: `comment##${commentIndex}`,
          owned_by: user.id,
          message: message,
          created_at: (new Date()).getTime(),
        }
      }).promise();

      await dbc.update({
        TableName: process.env.EntityTable,
        Key: {
          id: collection.id,
          sort: 'detail',
        },
        UpdateExpression: 'set comment_count = :count',
        ExpressionAttributeValues: {
          ':count': collection.comment_count + 1,
        }
      }).promise();

      if (!collection.comment_members.includes(user.id)) {
        const members = collection.comment_members;
        members.push(user.id);

        await dbc.update({
          TableName: process.env.EntityTable,
          Key: {
            id: collection.id,
            sort: 'detail',
          },
          UpdateExpression: 'set comment_members = :members',
          ExpressionAttributeValues: {
            ':members': members,
          }
        }).promise();
      }

      return {
        statusCode: 201,
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Location': `/collections/${collectionId}/comments/${commentIndex}`,
        },
        body: null,
      };
    }

    if (method === 'GET') {
      const collectionId = event.pathParameters.collectionId;

      const comments = (await dbc.query({
        TableName: process.env.EntityTable,
        KeyConditionExpression: 'id = :id and begins_with(sort, :sort)',
        ExpressionAttributeValues: {
          ':id': collectionId,
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
