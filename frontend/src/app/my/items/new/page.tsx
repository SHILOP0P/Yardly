"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
import { itemsApi } from "@/shared/api/endpoints/items";
import type { DealMode } from "@/shared/api/types";
import { useHydrateSession } from "@/shared/auth/useHydrateSession";
import { formatDealMode } from "@/features/item/presentation";

type CreateResult = {
  id: number;
  uploadedImages: number;
  failedImages: number;
};

const modes: DealMode[] = ["rent", "sale", "free", "sale_rent"];

export default function NewItemPage() {
  useHydrateSession();
  const router = useRouter();

  const [title, setTitle] = useState("");
  const [mode, setMode] = useState<DealMode>("rent");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [deposit, setDeposit] = useState("");
  const [location, setLocation] = useState("");
  const [category, setCategory] = useState("");
  const [imageFiles, setImageFiles] = useState<File[]>([]);
  const filesCount = useMemo(() => imageFiles.length, [imageFiles]);

  const mut = useMutation({
    mutationFn: async (): Promise<CreateResult> => {
      const created = await itemsApi.create({
        title: title.trim(),
        mode,
        description: description.trim(),
        price: price ? Number(price) : 0,
        deposit: deposit ? Number(deposit) : 0,
        location: location.trim(),
        category: category.trim(),
      });

      if (imageFiles.length === 0) {
        return { id: created.id, uploadedImages: 0, failedImages: 0 };
      }

      const uploaded = await Promise.allSettled(imageFiles.map((file) => itemsApi.images.add(created.id, file)));
      const uploadedImages = uploaded.filter((r) => r.status === "fulfilled").length;
      const failedImages = uploaded.length - uploadedImages;

      return { id: created.id, uploadedImages, failedImages };
    },
    onSuccess: (res) => {
      router.push(`/items/${res.id}`);
    },
  });

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Create product</h1>

      <div className="border rounded-xl p-4 space-y-2 text-sm">
        <div className="font-medium">What can be filled from current backend model</div>
        <div>- title: editable</div>
        <div>- mode: editable</div>
        <div>- description, price, deposit, location, category: editable</div>
        <div>- images: optional, uploaded as files after create</div>
        <div>- status: set by backend as active</div>
        <div>- owner_id: set by backend from current user</div>
      </div>

      <div className="space-y-2">
        <div className="text-sm opacity-70">Title</div>
        <input
          className="border rounded-lg px-3 py-2 w-full"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Example: Bosch drill"
        />
      </div>

      <div className="space-y-2">
        <div className="text-sm opacity-70">Mode</div>
        <select className="border rounded-lg px-3 py-2" value={mode} onChange={(e) => setMode(e.target.value as DealMode)}>
          {modes.map((modeValue) => (
            <option key={modeValue} value={modeValue}>
              {formatDealMode(modeValue)} ({modeValue})
            </option>
          ))}
        </select>
      </div>

      <div className="space-y-2">
        <div className="text-sm opacity-70">Description</div>
        <textarea
          className="border rounded-lg px-3 py-2 w-full min-h-24"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Condition, details, kit, etc."
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div className="space-y-2">
          <div className="text-sm opacity-70">Price</div>
          <input
            className="border rounded-lg px-3 py-2 w-full"
            type="number"
            min={0}
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            placeholder="0"
          />
        </div>
        <div className="space-y-2">
          <div className="text-sm opacity-70">Deposit</div>
          <input
            className="border rounded-lg px-3 py-2 w-full"
            type="number"
            min={0}
            value={deposit}
            onChange={(e) => setDeposit(e.target.value)}
            placeholder="0"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div className="space-y-2">
          <div className="text-sm opacity-70">Location</div>
          <input
            className="border rounded-lg px-3 py-2 w-full"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            placeholder="City / area"
          />
        </div>
        <div className="space-y-2">
          <div className="text-sm opacity-70">Category</div>
          <input
            className="border rounded-lg px-3 py-2 w-full"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            placeholder="Tools, transport, electronics..."
          />
        </div>
      </div>

      <div className="space-y-2">
        <div className="text-sm opacity-70">Image files (optional)</div>
        <input
          className="border rounded-lg px-3 py-2 w-full"
          type="file"
          accept=".png,.jpg,.jpeg,image/png,image/jpeg"
          multiple
          onChange={(e) => setImageFiles(Array.from(e.target.files ?? []))}
        />
        <div className="text-xs opacity-70">Selected files: {filesCount}</div>
      </div>

      <div className="flex flex-wrap gap-2">
        <button
          className="border rounded-lg px-3 py-2"
          onClick={() => mut.mutate()}
          disabled={mut.isPending || !title.trim()}
        >
          {mut.isPending ? "Creating..." : "Create product"}
        </button>
        <Link className="border rounded-lg px-3 py-2" href="/my/items">
          Back to my products
        </Link>
      </div>

      {mut.error ? (
        <div className="text-sm text-red-600">{String((mut.error as Error).message)}</div>
      ) : null}
    </div>
  );
}
