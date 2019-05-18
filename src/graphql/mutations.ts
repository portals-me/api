// tslint:disable
// this is an auto generated file. This will be overwritten

export const addCollection = `mutation AddCollection(
  $owner: String!
  $name: String!
  $title: String
  $description: String
) {
  addCollection(
    owner: $owner
    name: $name
    title: $title
    description: $description
  )
}
`;
export const updateCollection = `mutation UpdateCollection($id: ID!, $title: String, $description: String) {
  updateCollection(id: $id, title: $title, description: $description)
}
`;
export const updateCollectionName = `mutation UpdateCollectionName($id: ID!, $name: String!) {
  updateCollectionName(id: $id, name: $name)
}
`;
export const deleteCollection = `mutation DeleteCollection($id: ID!) {
  deleteCollection(id: $id)
}
`;
