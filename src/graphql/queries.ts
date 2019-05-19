// tslint:disable
// this is an auto generated file. This will be overwritten

export const getCollection = `query GetCollection($id: ID!) {
  getCollection(id: $id) {
    id
    owner
    name
    title
    description
    created_at
    updated_at
  }
}
`;
export const listCollectionsByOwner = `query ListCollectionsByOwner($owner: String!) {
  listCollectionsByOwner(owner: $owner) {
    id
    owner
    name
    title
    description
    created_at
    updated_at
  }
}
`;
