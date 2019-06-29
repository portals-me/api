import 'cross-fetch/polyfill';
import gql from 'graphql-tag';
import AWSAppSyncClient, { AUTH_TYPE, createAppSyncLink, AWSAppSyncClientOptions } from 'aws-appsync';
import * as queries from '../src/graphql/queries';
import * as mutations from '../src/graphql/mutations';
import * as API from '../src/API';
import aws_config from '../src/aws-exports.js';
import { ApolloLink } from 'apollo-link';
import { setContext } from "apollo-link-context";
import { createHttpLink } from "apollo-link-http";

const jwt = process.env.JWT_TOKEN;

const AppSyncConfig = {
  url: aws_config.aws_appsync_graphqlEndpoint,
  region: aws_config.aws_appsync_region,
  auth: {
    type: AUTH_TYPE.API_KEY,
    apiKey: aws_config.aws_appsync_apiKey,
  },
  disableOffline: true,
} as AWSAppSyncClientOptions;

const client = new AWSAppSyncClient(AppSyncConfig, {
  link: createAppSyncLink({
    ...AppSyncConfig,
    resultsFetcherLink: ApolloLink.from([
      setContext((request, previousContext) => ({
        headers: {
          ...previousContext.headers,
          Authorization: `Bearer ${jwt}`,
        }
      })),
      createHttpLink({
        uri: AppSyncConfig.url,
      })
    ])
  } as any)
});

describe('Collection', () => {
  describe('Collection lifecycle', () => {
    const argument: API.AddCollectionMutationVariables = {
      name: 'test-foobar',
      title: 'This is a title',
      description: 'Description'
    };
    let collection;

    it('should add a collection', async () => {
      const result = await client.mutate({
        mutation: gql(mutations.addCollection),
        variables: argument,
      });
      expect(result.data).toEqual(expect.anything());
      collection = result.data.addCollection;

      expect(collection.id).not.toBeNull();
      expect(collection.owner).not.toBeNull();
      expect(collection.name).toBe(argument.name);
      expect(collection.title).toBe(argument.title);
      expect(collection.description).toBe(argument.description);
      expect(collection.created_at).toBe(collection.updated_at);
    });

    it('should add an article and delete it', async () => {
      const articleInput: API.AddArticleMutationVariables = {
        collectionId: collection.id,
        entity: {
          format: 'oembed',
          type: 'share',
          url: 'http://example.com/share_something'
        },
        title: 'foooo',
        description: 'This is the description!!!',
      };

      let article;

      {
        const result = await client.mutate({
          mutation: gql(mutations.addArticle),
          variables: articleInput,
        });
        expect(result.data).toEqual(expect.anything());
        article = result.data.addArticle;

        expect(article.id).not.toBeNull();
        expect(article.entity).toEqual(expect.objectContaining(articleInput.entity));
        expect(article.title).toEqual(articleInput.title);
        expect(article.description).toEqual(articleInput.description);
      }

      {
        const variables: API.DeleteArticleMutationVariables = {
          collectionId: article.collectionId,
          id: article.id,
        };
        const result = await client.mutate({
          mutation: gql(mutations.deleteArticle),
          variables,
        });
        expect(result.data).toEqual(expect.anything());
        expect(result.data.deleteArticle.id).toBe(article.id);
      }
    });

    it('should delete a collection', async () => {
      const variables: API.DeleteCollectionMutationVariables = {
        id: collection.id,
      };
      const result = await client.mutate({
        mutation: gql(mutations.deleteCollection),
        variables,
      });
      expect(result.data).toEqual(expect.anything());
      expect(result.data.deleteCollection.id).toBe(collection.id);
    });
  });

  /*
  it('should not add a collection for non-authorized user', async () => {
    const variables: API.AddCollectionMutationVariables = {
      name: 'test',
    };
    const promise = client.mutate({
      mutation: gql(mutations.addCollection),
      variables,
    });
    expect(promise).rejects.toThrow('Not Authorized');
  });
  */
});
