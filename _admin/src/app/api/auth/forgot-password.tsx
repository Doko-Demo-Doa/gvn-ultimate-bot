import type { NextApiRequest, NextApiResponse } from "next";

type IResponseData = {
  code: number;
  message: string;
  data?: {
    token: string;
    email: string;
  };
};

type IRequestBody = {
  email: string;
};

import { Resend } from "resend";
import ForgotPasswordEmailComponent from "~/app/_components/forgot-pw-email.component";
const resend = new Resend(process.env.RESEND_API_KEY);

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<IResponseData>
) {
  const body: IRequestBody = req.body;

  if (!body.email) {
    return res.status(400).json({ code: 400, message: "Email is required" });
  }

  const resp = await fetch(
    `${process.env.NEXT_PUBLIC_BASE_API_URL}/api/forgot-password`
  );
  const data: IResponseData = await resp.json();

  await resend.emails.send({
    from: "Admin <info@aniviet.com>",
    to: body.email,
    subject: "Reset your password",
    react: <ForgotPasswordEmailComponent token={data.data?.token!} />,
  });

  if (req.method !== "POST") {
    return res.status(405).json({ code: 405, message: "Method Not Allowed" });
  }

  res.status(200).json({ code: 200, message: "Success" });
}
