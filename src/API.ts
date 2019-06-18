/* tslint:disable */
//  This file was automatically generated and should not be edited.

export type CoverInput = {
  color: string,
  sort: string,
};

export type EntityInput = {
  format: string,
  type: string,
  url: string,
};

export type AddCollectionMutationVariables = {
  owner: string,
  name: string,
  title?: string | null,
  description?: string | null,
  cover?: CoverInput | null,
  media?: Array< string > | null,
};

export type AddCollectionMutation = {
  addCollection:  {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    cover:  {
      __typename: "Cover",
      color: string,
      sort: string,
    } | null,
    media: Array< string > | null,
    created_at: string,
    updated_at: string,
  },
};

export type UpdateCollectionMutationVariables = {
  id: string,
  title?: string | null,
  description?: string | null,
};

export type UpdateCollectionMutation = {
  updateCollection:  {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    cover:  {
      __typename: "Cover",
      color: string,
      sort: string,
    } | null,
    media: Array< string > | null,
    created_at: string,
    updated_at: string,
  },
};

export type DeleteCollectionMutationVariables = {
  id: string,
};

export type DeleteCollectionMutation = {
  deleteCollection:  {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    cover:  {
      __typename: "Cover",
      color: string,
      sort: string,
    } | null,
    media: Array< string > | null,
    created_at: string,
    updated_at: string,
  } | null,
};

export type AddArticleMutationVariables = {
  collectionId: string,
  entity: EntityInput,
  title?: string | null,
  description?: string | null,
  owner: string,
};

export type AddArticleMutation = {
  addArticle:  {
    __typename: "Article",
    collectionId: string,
    id: string,
    entity:  {
      __typename: "Entity",
      format: string,
      type: string,
      url: string,
    },
    title: string | null,
    description: string | null,
    owner: string,
    created_at: string,
    updated_at: string,
  },
};

export type DeleteArticleMutationVariables = {
  collectionId: string,
  id: string,
};

export type DeleteArticleMutation = {
  deleteArticle:  {
    __typename: "Article",
    collectionId: string,
    id: string,
    entity:  {
      __typename: "Entity",
      format: string,
      type: string,
      url: string,
    },
    title: string | null,
    description: string | null,
    owner: string,
    created_at: string,
    updated_at: string,
  } | null,
};

export type GetCollectionQueryVariables = {
  id: string,
};

export type GetCollectionQuery = {
  getCollection:  {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    cover:  {
      __typename: "Cover",
      color: string,
      sort: string,
    } | null,
    media: Array< string > | null,
    created_at: string,
    updated_at: string,
  } | null,
};

export type ListCollectionsQueryVariables = {
  owner: string,
};

export type ListCollectionsQuery = {
  listCollections:  Array< {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    cover:  {
      __typename: "Cover",
      color: string,
      sort: string,
    } | null,
    media: Array< string > | null,
    created_at: string,
    updated_at: string,
  } > | null,
};

export type ListArticlesQueryVariables = {
  collectionId: string,
};

export type ListArticlesQuery = {
  listArticles:  Array< {
    __typename: "Article",
    collectionId: string,
    id: string,
    entity:  {
      __typename: "Entity",
      format: string,
      type: string,
      url: string,
    },
    title: string | null,
    description: string | null,
    owner: string,
    created_at: string,
    updated_at: string,
  } > | null,
};
