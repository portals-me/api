import AWS from "aws-sdk";
import bcrypt from "bcrypt";
import uuid from "uuid/v4";
import axios from "axios";
import * as API from "../src/API";
import * as mutations from "../src/graphql/mutations";

AWS.config.update({
  region: "ap-northeast-1"
});

const apiEnv: {
  appsync: {
    url: string;
    apiKey: string;
  };
} = JSON.parse(process.env.API_ENV);

const accountEnv: {
  restApi: string;
  tableName: string;
} = JSON.parse(process.env.ACCOUNT_ENV);

const Dynamo = new AWS.DynamoDB.DocumentClient();

const createUser = async (user: {
  id: string;
  name: string;
  password: string;
  picture: string;
  display_name: string;
}) => {
  await Dynamo.put({
    Item: Object.assign(user, {
      sort: "detail"
    }),
    TableName: accountEnv.tableName
  }).promise();

  await Dynamo.put({
    Item: {
      id: user.id,
      sort: `name-pass##${user.name}`,
      check_data: bcrypt.hashSync(user.password, 10)
    },
    TableName: accountEnv.tableName
  }).promise();
};

const deleteUser = async (user: { id: string; name: string }) => {
  await Dynamo.delete({
    Key: {
      id: user.id,
      sort: "detail"
    },
    TableName: accountEnv.tableName
  }).promise();

  await Dynamo.delete({
    Key: {
      id: user.id,
      sort: `name-pass##${user.name}`
    },
    TableName: accountEnv.tableName
  }).promise();
};

const adminUser = {
  id: uuid(),
  name: `admin_${uuid().replace(/\-/g, "_")}`,
  password: uuid(),
  display_name: "Admin",
  picture: "pic"
};

let adminUserJWT;

beforeAll(async () => {
  await createUser(adminUser);

  adminUserJWT = (await axios.post(`${accountEnv.restApi}/signin`, {
    auth_type: "password",
    data: {
      user_name: adminUser.name,
      password: adminUser.password
    }
  })).data;
});

afterAll(async () => {
  await deleteUser(adminUser);
});

describe("Post", () => {
  describe("Image", () => {
    it("should do smth", async () => {
      const r = await axios.post(
        `${apiEnv.appsync.url}`,
        {
          query: `mutation GenerateUploadURL { generateUploadURL(keys: ["foooo"]) }`
        },
        {
          headers: {
            Authorization: `Bearer ${adminUserJWT}`,
            "x-api-key": apiEnv.appsync.apiKey
          }
        }
      );

      console.log(JSON.stringify(r.data));
    });
  });
});
