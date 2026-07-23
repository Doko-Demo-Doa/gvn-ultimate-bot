"use client";

import { useMutation } from "@tanstack/react-query";
import { useUploadThing } from "~/utils/uploadthing";

async function optimizeEmbedImage(file: File) {
  const body = new FormData();
  body.append("file", file);

  const res = await fetch("/api/images/optimize", {
    method: "POST",
    body,
  });

  if (!res.ok) {
    const body = (await res.json().catch(() => null)) as {
      error?: string;
    } | null;
    throw new Error(body?.error || "Cannot optimize selected image");
  }

  const blob = await res.blob();
  const fileName = file.name.replace(/\.[^.]+$/, "") || "embed-image";
  return new File([blob], `${fileName}.jpg`, { type: "image/jpeg" });
}

export function useEmbedImageUpload() {
  const { startUpload, isUploading: isUploadThingUploading } =
    useUploadThing("imageUploader");

  const mutation = useMutation({
    mutationFn: async (file: File) => {
      const optimizedFile = await optimizeEmbedImage(file);
      const uploaded = await startUpload([optimizedFile]);
      const uploadedUrl = uploaded?.[0]?.ufsUrl;

      if (!uploadedUrl) {
        throw new Error("Upload completed without a file URL");
      }

      return uploadedUrl;
    },
    retry: false,
  });

  return {
    ...mutation,
    isUploading: mutation.isPending || isUploadThingUploading,
  };
}
