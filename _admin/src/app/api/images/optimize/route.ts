import { auth } from "@clerk/nextjs/server";
import sharp from "sharp";

export const runtime = "nodejs";

const MAX_IMAGE_DIMENSION = 1600;
const IMAGE_QUALITY = 86;

type UploadedImage = Blob & {
  name?: string;
};

function getUploadedImage(
  value: FormDataEntryValue | null,
): UploadedImage | null {
  if (
    typeof value === "object" &&
    value !== null &&
    "arrayBuffer" in value &&
    "type" in value &&
    "size" in value
  ) {
    return value as UploadedImage;
  }

  return null;
}

function isSupportedImage(file: UploadedImage) {
  return (
    file.type === "image/png" ||
    file.type === "image/jpeg" ||
    file.type === "image/webp"
  );
}

export async function POST(request: Request) {
  const { isAuthenticated } = await auth();

  if (!isAuthenticated) {
    return Response.json({ error: "Unauthorized" }, { status: 401 });
  }

  const formData = await request.formData();
  const file = getUploadedImage(formData.get("file"));

  if (!file) {
    return Response.json({ error: "Image file is required" }, { status: 400 });
  }

  if (!isSupportedImage(file)) {
    return Response.json(
      { error: "Only PNG, JPEG, and WebP images are supported" },
      { status: 400 },
    );
  }

  const input = Buffer.from(await file.arrayBuffer());
  const image = sharp(input).rotate();
  const metadata = await image.metadata();
  const resized = image.resize({
    width: MAX_IMAGE_DIMENSION,
    height: MAX_IMAGE_DIMENSION,
    fit: "inside",
    withoutEnlargement: true,
  });

  const shouldOutputPng =
    metadata.format === "png" || metadata.format === "webp";
  const optimized = shouldOutputPng
    ? await resized
        .png({ compressionLevel: 9, quality: IMAGE_QUALITY })
        .toBuffer()
    : await resized
        .jpeg({
          quality: IMAGE_QUALITY,
          mozjpeg: true,
        })
        .toBuffer();
  const extension = shouldOutputPng ? "png" : "jpg";

  return new Response(optimized, {
    headers: {
      "Content-Disposition": `inline; filename="optimized.${extension}"`,
      "Content-Type": shouldOutputPng ? "image/png" : "image/jpeg",
      "X-Image-Extension": extension,
    },
  });
}
