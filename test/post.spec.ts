import AWS from "aws-sdk";
import bcrypt from "bcrypt";
import uuid from "uuid/v4";
import axios from "axios";
import * as API from "../src/API";
import * as mutations from "../src/graphql/mutations";
import FormData from "form-data";
import fs from "fs";

AWS.config.update({
  region: "ap-northeast-1"
});

const Dynamo = new AWS.DynamoDB.DocumentClient();
const S3 = new AWS.S3();

const apiEnv: {
  appsync: {
    url: string;
    apiKey: string;
  };
  userStorageBucket: string;
  postTableName: string;
} = JSON.parse(process.env.API_ENV);

const accountEnv: {
  restApi: string;
  tableName: string;
} = JSON.parse(process.env.ACCOUNT_ENV);

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
    const filename = "package.json";

    afterAll(async () => {
      await S3.deleteObject({
        Bucket: apiEnv.userStorageBucket,
        Key: `${adminUser.id}/${filename}`
      }).promise();
    });

    it("should upload a file", async () => {
      const [url] = (await axios.post(
        apiEnv.appsync.url,
        {
          query: `mutation GenerateUploadURL { generateUploadURL(keys: [${JSON.stringify(
            filename
          )}]) }`
        },
        {
          headers: {
            Authorization: `Bearer ${adminUserJWT}`,
            "x-api-key": apiEnv.appsync.apiKey
          }
        }
      )).data.data.generateUploadURL;
      expect(url).toBeTruthy();

      const form = new FormData();
      form.append(filename, fs.readFileSync(filename));

      const result = await axios.put(url, form, {
        headers: {
          "Content-Length": form.getLengthSync(),
          ...form.getHeaders()
        }
      });
      expect(result.status).toBe(200);
    });

    it("should create an image post", async () => {
      const result = await axios.post(
        apiEnv.appsync.url,
        {
          query: `mutation AddImagePost { addImagePost(
          title: "test post"
          description: "description"
          entity: {
            images: [
              {
                filetype: "text/json",
                s3path: "${filename}"
              }
            ]
          }
        ) { id } }`
        },
        {
          headers: {
            Authorization: `Bearer ${adminUserJWT}`,
            "x-api-key": apiEnv.appsync.apiKey
          }
        }
      );
      expect(result.data.data.addImagePost.id).toBeTruthy();

      await Dynamo.delete({
        TableName: apiEnv.postTableName,
        Key: {
          id: result.data.data.addImagePost.id,
          sort: "summary"
        }
      }).promise();
    });
  });
});
