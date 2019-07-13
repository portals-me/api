import AWS from "aws-sdk";
import uuid from "uuid/v4";
import axios from "axios";
import { createUser, deleteUser } from "./user";

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
  accountReplicaTableName: string;
} = JSON.parse(process.env.API_ENV);

const accountEnv: {
  restApi: string;
  tableName: string;
} = JSON.parse(process.env.ACCOUNT_ENV);

const adminUser = {
  id: uuid(),
  name: `admin_${uuid().replace(/\-/g, "_")}`,
  password: uuid(),
  display_name: "Admin",
  picture: "pic"
};

let adminUserJWT;

const guestUser = {
  id: uuid(),
  name: `guest_${uuid().replace(/\-/g, "_")}`
};

beforeAll(async () => {
  const userRecord = await createUser(accountEnv.tableName, adminUser);
  await Dynamo.put({
    Item: userRecord,
    TableName: apiEnv.accountReplicaTableName
  }).promise();

  adminUserJWT = (await axios.post(`${accountEnv.restApi}/signin`, {
    auth_type: "password",
    data: {
      user_name: adminUser.name,
      password: adminUser.password
    }
  })).data;

  await Dynamo.put({
    Item: Object.assign(guestUser, {
      sort: "detail"
    }),
    TableName: apiEnv.accountReplicaTableName
  }).promise();
});

afterAll(async () => {
  await deleteUser(accountEnv.tableName, adminUser);
  await Dynamo.delete({
    Key: {
      id: adminUser.id,
      sort: "detail"
    },
    TableName: apiEnv.accountReplicaTableName
  }).promise();
  await Dynamo.delete({
    Key: {
      id: adminUser.id,
      sort: "social"
    },
    TableName: apiEnv.accountReplicaTableName
  }).promise();

  await Dynamo.delete({
    Key: {
      id: guestUser.id,
      sort: "detail"
    },
    TableName: apiEnv.accountReplicaTableName
  }).promise();
  await Dynamo.delete({
    Key: {
      id: guestUser.id,
      sort: "social"
    },
    TableName: apiEnv.accountReplicaTableName
  }).promise();
});

describe("User", () => {
  it("should show UserMore information", async () => {
    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `query Q {
          getUserMoreByName(name: "${adminUser.name}") {
            id
            name
            display_name
            picture
            is_following
            followings
            followers
          }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${adminUserJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    expect(result.data.data.getUserMoreByName).toBeTruthy();

    const userMore = result.data.data.getUserMoreByName;

    expect(userMore.id).toBe(adminUser.id);
    expect(userMore.name).toBe(adminUser.name);
    expect(userMore.display_name).toBe(adminUser.display_name);
    expect(userMore.picture).toBe(adminUser.picture);
    expect(userMore.is_following).toBe(false);
    expect(userMore.followings).toBe(0);
    expect(userMore.followers).toBe(0);
  });

  it("should follow and show correct user information", async () => {
    {
      const result = await axios.post(
        apiEnv.appsync.url,
        {
          query: `mutation M {
            followUser(targetId: "${guestUser.id}") { id }
          }`
        },
        {
          headers: {
            Authorization: `Bearer ${adminUserJWT}`,
            "x-api-key": apiEnv.appsync.apiKey
          }
        }
      );

      expect(result.data.errors).not.toBeTruthy();
      expect(result.data.data.followUser.id).toBeTruthy();
    }

    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `query Q {
          getUserMoreByName(name: "${guestUser.name}") {
            id
            is_following
            followings
            followers
          }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${adminUserJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );
    expect(result.data.errors).not.toBeTruthy();
    const user = result.data.data.getUserMoreByName;

    expect(user.id).toBe(guestUser.id);
    expect(user.is_following).toBe(true);
    expect(user.followings).toBe(0);
    // Because of the eventually consistency, it's hard to check this.
    // expect(user.followers).toBe(1);
  });

  it("shoud not follow twice", async () => {
    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `mutation M {
          followUser(targetId: "${guestUser.id}") { id }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${adminUserJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    expect(result.data.data).not.toBeTruthy();
    expect(result.data.errors).toBeTruthy();
  });

  it("should unfollow", async () => {
    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `mutation M {
          unfollowUser(targetId: "${guestUser.id}") { id }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${adminUserJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    expect(result.data.errors).not.toBeTruthy();
    expect(result.data.data.unfollowUser.id).toBeTruthy();
  });

  it("should not unfollow twice", async () => {
    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `mutation M {
          unfollowUser(targetId: "${guestUser.id}") { id }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${adminUserJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    expect(result.data.errors).not.toBeTruthy();
    expect(result.data.data.unfollowUser.id).toBeTruthy();
  });
});
