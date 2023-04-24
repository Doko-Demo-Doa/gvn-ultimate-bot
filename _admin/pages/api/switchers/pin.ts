import { NextApiRequest, NextApiResponse } from "next";

const baseUrl = process.env.BACKEND_BASE_URL;

const handler = async (req: NextApiRequest, res: NextApiResponse) => {
  const { method, body } = req;
  switch (method) {
    case "POST":
      try {
        const payload: { enabled: ModuleActivationStatusType } =
          JSON.parse(body);

        const r = await fetch(baseUrl + "/switchers/pin", {
          method: "POST",
          body: JSON.stringify(payload),
        });

        const respJson = await r.json();

        res.status(r.status).json(respJson);
        break;
      } catch (error) {
        res.status(500).json(error);
        break;
      }

    case "GET":
      try {
        const r = await fetch(baseUrl + "/switchers/pin");
        const respJson: ModuleActivationStatusType = await r.json();

        res.status(r.status).json(respJson);
        break;
      } catch (error) {
        res.status(500).json(error);
        break;
      }

    default:
      res.setHeader("Allow", ["POST", "GET"]);
      res.status(405).end(`Method ${method} Not Allowed`);
      break;
  }
};

export default handler;
