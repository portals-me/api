import 'cross-fetch/polyfill';
import gql from 'graphql-tag';
import AWSAppSyncClient, { AUTH_TYPE } from 'aws-appsync';
import * as queries from '../src/graphql/queries';
import * as mutations from '../src/graphql/mutations';
import * as API from '../src/API';
import aws_config from '../src/aws-exports.js';

const client = new AWSAppSyncClient({
  url: aws_config.aws_appsync_graphqlEndpoint,
  region: aws_config.aws_appsync_region,
  auth: {
    type: AUTH_TYPE.API_KEY,
    apiKey: aws_config.aws_appsync_apiKey,
  },
  disableOffline: true,
});

describe('Collection', () => {
  describe('Collection lifecycle', () => {
    const argument = {
      owner: 'foobar',
      name: 'test-foobar',
      title: 'This is a title',
      description: 'Description'
    }
    let collection;

    it('should add a collection', async () => {
      const result = await client.mutate({
        mutation: gql(mutations.addCollection),
        variables: argument as API.AddCollectionMutationVariables,
      });
      expect(result.data).toEqual(expect.anything());
      collection = result.data.addCollection;

      expect(collection.id).not.toBeNull();
      expect(collection.owner).toBe(argument.owner);
      expect(collection.name).toBe(argument.name);
      expect(collection.title).toBe(argument.title);
      expect(collection.description).toBe(argument.description);
      expect(collection.created_at).toBe(collection.updated_at);
    });

    it('should delete a collection', async () => {
      const result = await client.mutate({
        mutation: gql(mutations.deleteCollection),
        variables: {
          id: collection.id
        } as API.DeleteCollectionMutationVariables,
      });
      expect(result.data).toEqual(expect.anything());
      expect(result.data.deleteCollection).toBe(collection.id);
    });
  });
});
