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
  description: string;
  price: number;
  deposit: number;
  location: string;
  category: string;
  images?: ItemImage[];
};

export type FavoriteItem = {
  item_id: number;
  title: string;
  status: string;
  mode: string;
  owner_id: number;
  favorited_at: string;
};

export type Tokens = {
  access_token: string;
  refresh_token?: string;
  token_type?: string;
};

export type RegisterResponse = {
  id: number;
  email: string;
  role: UserRole;
  first_name: string;
  last_name?: string | null;
};

export type UserRole = "user" | "admin" | "superadmin";

export type Me = {
  id: number;
  email: string;
  role: UserRole;
  first_name?: string;
  last_name?: string | null;
  profile?: { first_name?: string; last_name?: string };
};

export type AdminListResp<T, K extends string> = {
  [P in K]: T[];
} & {
  limit: number;
  offset: number;
};

export type AdminUser = {
  id: number;
  email: string;
  role: UserRole;
  banned_at?: string | null;
  ban_expires_at?: string | null;
  ban_reason?: string | null;
  created_at: string;
  updated_at: string;
};

export type AdminBooking = {
  id: number;
  item_id: number;
  requester_id: number;
  owner_id: number;
  type: string;
  status: string;
  start?: string | null;
  end?: string | null;
  handover_deadline?: string | null;
  handover_confirmed_by_owner_at?: string | null;
  handover_confirmed_by_requester_at?: string | null;
  return_confirmed_by_owner_at?: string | null;
  return_confirmed_by_requester_at?: string | null;
  created_at: string;
};

export type AdminBookingEvent = {
  id: number;
  booking_id: number;
  actor_user_id?: number | null;
  action: string;
  from_status?: string | null;
  to_status?: string | null;
  meta?: unknown;
  created_at: string;
};

export type AdminItem = {
  id: number;
  owner_id: number;
  title: string;
  status: ItemStatus;
  mode: DealMode;
  blocked_at?: string | null;
  block_reason?: string | null;
};

export type AdminEvent = {
  id: number;
  actor_user_id: number;
  entity_type: string;
  entity_id: number;
  action: string;
  reason?: string | null;
  meta?: unknown;
  created_at: string;
};
