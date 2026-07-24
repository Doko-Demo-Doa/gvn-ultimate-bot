import { v4 as uuidv4 } from "uuid";
import { customApiClient } from "~/datasource/rest/api-client";

export async function optimizeImage(file: File) {
  const formData = new FormData();
  formData.append("file", file);

  const res = await customApiClient.postFormData(
    "/api/images/optimize",
    formData,
    {
      useBaseUrl: false,
    },
  );

  if (!res.ok) {
    const body = (await res.json().catch(() => null)) as {
      error?: string;
    } | null;
    throw new Error(body?.error || "Không thể tối ưu ảnh này");
  }

  const blob = await res.blob();
  const extension = res.headers.get("x-image-extension") || "jpg";
  const contentType = res.headers.get("content-type") || "image/jpeg";

  return new File([blob], `${uuidv4()}.${extension}`, { type: contentType });
}
