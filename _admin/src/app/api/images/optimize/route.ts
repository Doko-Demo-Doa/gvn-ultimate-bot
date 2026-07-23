import { auth } from "@clerk/nextjs/server";
import sharp from "sharp";

const MAX_IMAGE_DIMENSION = 1600;
const IMAGE_QUALITY = 86;

export async function POST(request: Request) {
  const { isAuthenticated } = await auth();

  if (!isAuthenticated) {
    return Response.json({ error: "Unauthorized" }, { status: 401 });
  }

  const formData = await request.formData();
  const file = formData.get("file");

  if (!(file instanceof File)) {
    return Response.json({ error: "Image file is required" }, { status: 400 });
  }

  if (!file.type.startsWith("image/")) {
    return Response.json(
      { error: "Only images are supported" },
      { status: 400 },
    );
  }

  const input = Buffer.from(await file.arrayBuffer());
  const optimized = await sharp(input)
    .rotate()
    .resize({
      width: MAX_IMAGE_DIMENSION,
      height: MAX_IMAGE_DIMENSION,
      fit: "inside",
      withoutEnlargement: true,
    })
    .jpeg({
      quality: IMAGE_QUALITY,
      mozjpeg: true,
    })
    .toBuffer();

  return new Response(optimized, {
    headers: {
      "Content-Disposition": 'inline; filename="embed-image.jpg"',
      "Content-Type": "image/jpeg",
    },
  });
}
