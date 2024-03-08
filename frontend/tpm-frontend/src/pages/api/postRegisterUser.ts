// pages/api/createUser.ts

import { NextApiRequest, NextApiResponse } from "next";

export default async function PostRegisterUser(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method === "POST") {
    const { userEmail, userPassword } = req.body;

    try {
      // Perform user creation logic here
      // You may interact with your database or any other backend service
      const response = await fetch("http://tpm-backend:80/api/v1/createUser", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ userEmail, userPassword }),
      });

      if (response.status == 200) {
        const dataRes = await response.json();
        res.status(200).json(dataRes);
      } else {
        // If response is not successful, parse JSON response to get error message and status code
        const errorResponse = await response.json();
        res.status(500).json(errorResponse);
      }
    } catch (error) {
      console.error("Error creating user:", error);
      res.status(500).json({ error: "Internal server error" });
    }
  } else {
    res.status(405).json({ error: "Method Not Allowed" });
  }
}
