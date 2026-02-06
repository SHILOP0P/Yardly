export type ItemStatus = "active" | "in_use" | "archived" | "deleted" | "transferred";
export type DealMode = "sale" | "rent" | "free" | "sale_rent";

export type ItemImage = {
  id: number;
  item_id: number;
  url: string;
  sort_order: number;
  created_at: string;
};

export type Item = {
  id: number;
  owner_id: number;
  title: string;
  status: ItemStatus;
  mode: DealMode;
  images?: ItemImage[];
};

export type Tokens = {
  access_token: string;
  refresh_token: string;
  token_type?: string;
};

export type Me = {
  id: number;
  email: string;
  role?: string;
  profile?: { first_name?: string; last_name?: string };
};
