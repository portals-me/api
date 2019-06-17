// tslint:disable
// this is an auto generated file. This will be overwritten

export const getCollection = `query GetCollection($id: ID!) {
  getCollection(id: $id) {
    id
    owner
    name
    title
    description
    cover {
      color
      sort
    }
    media
    created_at
    updated_at
  }
}
`;
export const listCollections = `query ListCollections($owner: String!) {
  listCollections(owner: $owner) {
    id
    owner
    name
    title
    description
    cover {
      color
      sort
    }
    media
    created_at
    updated_at
  }
}
`;
export const listArticles = `query ListArticles($collection: String!) {
  listArticles(collection: $collection) {
    collectionId
    id
    entity {
      format
      type
      url
    }
    title
    description
    owner
    created_at
    updated_at
  }
}
`;
