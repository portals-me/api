// tslint:disable
// this is an auto generated file. This will be overwritten

export const getUserByName = `query GetUserByName($name: String!) {
  getUserByName(name: $name) {
    id
    name
    picture
    display_name
  }
}
`;
export const getPostSummary = `query GetPostSummary($id: ID!) {
  getPostSummary(id: $id) {
    id
    title
    description
    updated_at
    created_at
    owner
    entity_type
    entity {
      ... on Share {
        format
        url
      }
    }
  }
}
`;
export const listPostSummary = `query ListPostSummary($owner: String) {
  listPostSummary(owner: $owner) {
    id
    title
    description
    updated_at
    created_at
    owner
    entity_type
    entity {
      ... on Share {
        format
        url
      }
    }
  }
}
`;
