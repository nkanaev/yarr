export type ItemStatus = "unread" | "read" | "starred";

export interface MediaLink {
  url: string;
  type: string;
  description?: string;
}

export interface Feed {
  id: number;
  folder_id: number | null;
  title: string;
  description: string;
  link: string;
  feed_link: string;
  icon?: string | null;
}

export interface Folder {
  id: number;
  title: string;
  is_expanded: boolean;
}

export interface Item {
  id: number;
  guid: string;
  feed_id: number;
  title: string;
  link: string;
  content?: string;
  date: string;
  status: ItemStatus;
  media_links: MediaLink[];
}

export interface Settings {
  filter: string;
  feed: string;
  feed_list_width: number;
  item_list_width: number;
  sort_newest_first: boolean;
  theme_name: string;
  theme_font: string;
  theme_size: number;
  refresh_rate: number;
  language: string;
}

export interface FeedStat {
  feed_id: number;
  unread: number;
  starred: number;
}

export interface StatusResponse {
  running: number;
  stats: FeedStat[];
}

export interface ItemListResponse {
  list: Item[];
  has_more: boolean;
}

export interface FeedLink {
  url: string;
  title: string;
  title_override?: string;
}

export interface FeedCreateSuccess {
  status: "success";
  feed: Feed;
}

export interface FeedCreateMultiple {
  status: "multiple";
  choice: FeedLink[];
}

export type FeedCreateResponse = FeedCreateSuccess | FeedCreateMultiple | { status: "notfound" };

export interface CrawlResponse {
  content: string;
}

export interface FeedCreateData {
  url: string;
  folder_id?: number | null;
  title_override?: string;
}

export interface FeedUpdateData {
  title?: string;
  folder_id?: number | null;
  feed_link?: string;
}

export interface FolderCreateData {
  title: string;
}

export interface FolderUpdateData {
  title?: string;
  is_expanded?: boolean;
}

export interface ItemUpdateData {
  status?: ItemStatus;
}

export interface ItemListQuery {
  feed_id?: string;
  folder_id?: string;
  status?: string;
  search?: string;
  oldest_first?: boolean;
  after?: number;

  [key: string]: string | number | boolean | undefined;
}

export interface ItemMarkQuery {
  feed_id?: string;
  folder_id?: string;
  before?: number;

  [key: string]: string | number | undefined;
}

export interface SettingsUpdateData {
  filter?: string;
  feed?: string;
  feed_list_width?: number;
  item_list_width?: number;
  sort_newest_first?: boolean;
  theme_name?: string;
  theme_font?: string;
  theme_size?: number;
  refresh_rate?: number;
  language?: string;
}
