/* tslint:disable */
//  This file was automatically generated and should not be edited.

export type AddCollectionMutationVariables = {
  owner: string,
  name: string,
  title?: string | null,
  description?: string | null,
};

export type AddCollectionMutation = {
  addCollection: string,
};

export type UpdateCollectionMutationVariables = {
  id: string,
  title?: string | null,
  description?: string | null,
};

export type UpdateCollectionMutation = {
  updateCollection: string | null,
};

export type UpdateCollectionNameMutationVariables = {
  id: string,
  name: string,
};

export type UpdateCollectionNameMutation = {
  updateCollectionName: string | null,
};

export type DeleteCollectionMutationVariables = {
  id: string,
};

export type DeleteCollectionMutation = {
  deleteCollection: string | null,
};

export type GetCollectionQueryVariables = {
  id?: string | null,
};

export type GetCollectionQuery = {
  getCollection:  {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    created_at: string,
    updated_at: string,
  } | null,
};

export type ListCollectionsByOwnerQueryVariables = {
  owner?: string | null,
};

export type ListCollectionsByOwnerQuery = {
  listCollectionsByOwner:  Array< {
    __typename: "Collection",
    id: string,
    owner: string,
    name: string,
    title: string | null,
    description: string | null,
    created_at: string,
    updated_at: string,
  } | null > | null,
};
