/* tslint:disable */
//  This file was automatically generated and should not be edited.

export type ShareInput = {
  format: string,
  url: string,
};

export type ImagesInput = {
  images: Array< ImageInput >,
};

export type ImageInput = {
  filetype: string,
  s3path: string,
};

export type AddSharePostMutationVariables = {
  title?: string | null,
  description?: string | null,
  entity: ShareInput,
};

export type AddSharePostMutation = {
  addSharePost:  {
    __typename: "PostSummary",
    id: string,
    title: string | null,
    description: string | null,
    updated_at: string,
    created_at: string,
    owner: string,
    entity_type: string,
    entity: ( {
        __typename: "Share",
        format: string,
        url: string,
      } | {
        __typename: "Images",
      }
    ),
  },
};

export type AddImagePostMutationVariables = {
  title?: string | null,
  description?: string | null,
  entity: ImagesInput,
};

export type AddImagePostMutation = {
  addImagePost:  {
    __typename: "PostSummary",
    id: string,
    title: string | null,
    description: string | null,
    updated_at: string,
    created_at: string,
    owner: string,
    entity_type: string,
    entity: ( {
        __typename: "Share",
        format: string,
        url: string,
      } | {
        __typename: "Images",
      }
    ),
  },
};

export type GetUserByNameQueryVariables = {
  name: string,
};

export type GetUserByNameQuery = {
  getUserByName:  {
    __typename: "User",
    id: string,
    name: string,
    picture: string | null,
    display_name: string | null,
  } | null,
};

export type GetPostSummaryQueryVariables = {
  id: string,
};

export type GetPostSummaryQuery = {
  getPostSummary:  {
    __typename: "PostSummary",
    id: string,
    title: string | null,
    description: string | null,
    updated_at: string,
    created_at: string,
    owner: string,
    entity_type: string,
    entity: ( {
        __typename: "Share",
        format: string,
        url: string,
      } | {
        __typename: "Images",
      }
    ),
  },
};

export type ListPostSummaryQueryVariables = {
  owner?: string | null,
};

export type ListPostSummaryQuery = {
  listPostSummary:  Array< {
    __typename: "PostSummary",
    id: string,
    title: string | null,
    description: string | null,
    updated_at: string,
    created_at: string,
    owner: string,
    entity_type: string,
    entity: ( {
        __typename: "Share",
        format: string,
        url: string,
      } | {
        __typename: "Images",
      }
    ),
  } >,
};
