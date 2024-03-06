// pages/api/createUser.ts

import { NextApiRequest, NextApiResponse } from "next";

export default async function PostRegisterUser(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method === "POST") {
    const { data } = req.body;

    try {
      // Perform user creation logic here
      // You may interact with your database or any other backend service
      const response = await fetch(
        "http://tpm-backend:8080/api/v1/createUser",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        }
      );

      const dataRes = await response.json();

      res.status(200).json(dataRes);
    } catch (error) {
      console.error("Error creating user:", error);
      res.status(500).json({ error: "Internal server error" });
    }
  } else {
    res.status(405).json({ error: "Method Not Allowed" });
  }
}
