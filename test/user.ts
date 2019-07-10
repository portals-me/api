import AWS from "aws-sdk";
import bcrypt from "bcrypt";

AWS.config.update({
  region: "ap-northeast-1"
});

const Dynamo = new AWS.DynamoDB.DocumentClient();

export const createUser = async (
  tableName: string,
  user: {
    id: string;
    name: string;
    password: string;
    picture: string;
    display_name: string;
  }
) => {
  await Dynamo.put({
    Item: Object.assign(user, {
      sort: "detail"
    }),
    TableName: tableName
  }).promise();

  await Dynamo.put({
    Item: {
      id: user.id,
      sort: `name-pass##${user.name}`,
      check_data: bcrypt.hashSync(user.password, 10)
    },
    TableName: tableName
  }).promise();
};

export const deleteUser = async (
  tableName: string,
  user: { id: string; name: string }
) => {
  await Dynamo.delete({
    Key: {
      id: user.id,
      sort: "detail"
    },
    TableName: tableName
  }).promise();

  await Dynamo.delete({
    Key: {
      id: user.id,
      sort: `name-pass##${user.name}`
    },
    TableName: tableName
  }).promise();
};
