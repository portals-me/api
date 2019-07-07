// tslint:disable
// this is an auto generated file. This will be overwritten

export const addSharePost = `mutation AddSharePost(
  $title: String
  $description: String
  $entity: ShareInput!
) {
  addSharePost(title: $title, description: $description, entity: $entity) {
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
export const addImagePost = `mutation AddImagePost(
  $title: String
  $description: String
  $entity: ImagesInput!
) {
  addImagePost(title: $title, description: $description, entity: $entity) {
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
