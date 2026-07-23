import { auth } from "@clerk/nextjs/server";

export async function GET(_request: Request) {
  const { isAuthenticated } = await auth();
  if (!isAuthenticated) {
    return Response.json({ error: "Unauthorized" }, { status: 401 });
  }

  return Response.json({});
}
