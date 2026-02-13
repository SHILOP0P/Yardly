import type { DealMode, ItemStatus } from "@/shared/api/types";

const modeLabels: Record<DealMode, string> = {
  sale: "For sale",
  rent: "For rent",
  free: "Give away",
  sale_rent: "Sale or rent",
};

const statusLabels: Record<ItemStatus, string> = {
  active: "Active",
  in_use: "In use",
  archived: "Archived",
  deleted: "Deleted",
  transferred: "Transferred",
};

export function formatDealMode(mode: DealMode): string {
  return modeLabels[mode] ?? mode;
}

export function formatItemStatus(status: ItemStatus): string {
  return statusLabels[status] ?? status;
}
