
export interface Collection {
  comment_count: number,
  comment_members: Array<string>,
  cover: Map<string, string>,
  created_at: number,
  description: string,
  id: string,
  media: Map<string, string>,
  owner: string,
  title: string,
}
