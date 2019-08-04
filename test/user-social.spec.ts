import AWS from "aws-sdk";
import uuid from "uuid/v4";
import axios from "axios";
import promiseRetry from "promise-retry";
import { createUser, deleteUser } from "./user";

AWS.config.update({
  region: "ap-northeast-1"
});

jest.setTimeout(30000);

const Dynamo = new AWS.DynamoDB.DocumentClient();

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

const Alice = {
  id: uuid(),
  name: `alice_${uuid().replace(/\-/g, "_")}`,
  password: uuid(),
  display_name: "Alice",
  picture: "pic"
};

let AliceJWT;

const Bob = {
  id: uuid(),
  name: `bob_${uuid().replace(/\-/g, "_")}`,
  password: uuid(),
  display_name: "Bob",
  picture: "pic"
};

beforeAll(async () => {
  await createUser(accountEnv.tableName, Alice);
  await createUser(accountEnv.tableName, Bob);

  AliceJWT = (await axios.post(`${accountEnv.restApi}/signin`, {
    auth_type: "password",
    data: {
      user_name: Alice.name,
      password: Alice.password
    }
  })).data;
});

afterAll(async () => {
  await deleteUser(accountEnv.tableName, Alice);
  await deleteUser(accountEnv.tableName, Bob);
});

describe("Scenario: user follow and unfollow", () => {
  it("should push to replica table", async () => {
    const result = await promiseRetry(
      async (retry, number) => {
        const result = await Dynamo.get({
          TableName: apiEnv.accountReplicaTableName,
          Key: {
            id: Alice.id,
            sort: "detail"
          }
        }).promise();

        if (!result.Item) {
          retry(result);
        } else {
          return result;
        }
      },
      {
        retries: 3,
        minTimeout: 100
      }
    );

    expect(result.Item).toBeTruthy();
    expect(result.Item.id).toBe(Alice.id);
  });

  it("should put social record", async () => {
    const result = await Dynamo.get({
      TableName: apiEnv.accountReplicaTableName,
      Key: {
        id: Alice.id,
        sort: "social"
      }
    }).promise();

    expect(result.Item.followers).toBe(0);
    expect(result.Item.followings).toBe(0);
  });

  it("should show UserMore information", async () => {
    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `query Q {
          getUserMoreByName(name: "${Alice.name}") {
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
          Authorization: `Bearer ${AliceJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    expect(result.data.data.getUserMoreByName).toBeTruthy();

    const userMore = result.data.data.getUserMoreByName;

    expect(userMore.id).toBe(Alice.id);
    expect(userMore.name).toBe(Alice.name);
    expect(userMore.display_name).toBe(Alice.display_name);
    expect(userMore.picture).toBe(Alice.picture);
    expect(userMore.is_following).toBe(false);
    expect(userMore.followings).toBe(0);
    expect(userMore.followers).toBe(0);
  });

  it("should follow", async () => {
    await axios.post(
      apiEnv.appsync.url,
      {
        query: `mutation M {
            followUser(targetId: "${Bob.id}") { id }
          }`
      },
      {
        headers: {
          Authorization: `Bearer ${AliceJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    const user = await promiseRetry(
      async (retry, number) => {
        const result = await axios.post(
          apiEnv.appsync.url,
          {
            query: `query Q {
            getUserMoreByName(name: "${Bob.name}") {
              id
              is_following
              followings
              followers
            }
          }`
          },
          {
            headers: {
              Authorization: `Bearer ${AliceJWT}`,
              "x-api-key": apiEnv.appsync.apiKey
            }
          }
        );

        expect(result.data.errors).not.toBeTruthy();
        const user = result.data.data.getUserMoreByName;

        if (user.followers == 0) {
          retry(user);
        } else {
          return user;
        }
      },
      {
        retries: 3,
        minTimeout: 100
      }
    );

    expect(user.id).toBe(Bob.id);
    expect(user.is_following).toBe(true);
    expect(user.followings).toBe(0);
    expect(user.followers).toBeLessThanOrEqual(1);
  });

  it("shoud not follow twice", async () => {
    await axios.post(
      apiEnv.appsync.url,
      {
        query: `mutation M {
          followUser(targetId: "${Bob.id}") { id }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${AliceJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `query Q {
          getUserMoreByName(name: "${Bob.name}") {
            id
            is_following
            followings
            followers
          }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${AliceJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );
    expect(result.data.errors).not.toBeTruthy();
    const user = result.data.data.getUserMoreByName;

    expect(user.id).toBe(Bob.id);
    expect(user.is_following).toBe(true);
    expect(user.followings).toBe(0);
    expect(user.followers).toBeLessThanOrEqual(1);
  });

  it("should unfollow", async () => {
    const result = await axios.post(
      apiEnv.appsync.url,
      {
        query: `mutation M {
          unfollowUser(targetId: "${Bob.id}") { id }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${AliceJWT}`,
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
          unfollowUser(targetId: "${Bob.id}") { id }
        }`
      },
      {
        headers: {
          Authorization: `Bearer ${AliceJWT}`,
          "x-api-key": apiEnv.appsync.apiKey
        }
      }
    );

    expect(result.data.errors).toBeTruthy();
    expect(result.data.data).not.toBeTruthy();
  });
});
