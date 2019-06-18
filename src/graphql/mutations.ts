// tslint:disable
// this is an auto generated file. This will be overwritten

export const addCollection = `mutation AddCollection(
  $name: String!
  $title: String
  $description: String
  $cover: CoverInput
  $media: [String!]
) {
  addCollection(
    name: $name
    title: $title
    description: $description
    cover: $cover
    media: $media
  ) {
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
export const updateCollection = `mutation UpdateCollection($id: ID!, $title: String, $description: String) {
  updateCollection(id: $id, title: $title, description: $description) {
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
export const deleteCollection = `mutation DeleteCollection($id: ID!) {
  deleteCollection(id: $id) {
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
export const addArticle = `mutation AddArticle(
  $collectionId: String!
  $entity: EntityInput!
  $title: String
  $description: String
) {
  addArticle(
    collectionId: $collectionId
    entity: $entity
    title: $title
    description: $description
  ) {
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
export const deleteArticle = `mutation DeleteArticle($collectionId: String!, $id: String!) {
  deleteArticle(collectionId: $collectionId, id: $id) {
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
