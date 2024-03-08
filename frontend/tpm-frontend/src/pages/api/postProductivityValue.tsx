// use this to post productivity value to db storing data

import { NextApiResponse, NextApiRequest } from "next";
import { parseCookies } from "nookies";

export default async function postProductivityValue(
  req: NextApiRequest,
  res: NextApiResponse
) {
  try {
    const data = req.body;

    const cookies = parseCookies({ req });

    // Extract the jwt cookie
    const jwtCookie = cookies.jwt;
    console.log(jwtCookie);
    const response = await fetch(
      // needs to be updated to something else probably
      `http://tpm-backend:80/api/v1/restricted/userData`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Authorization: jwtCookie,
        },
        body: JSON.stringify(data), // include date and data in the request body
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const dataRes = await response.json();

    res.status(200).json(dataRes);
  } catch (error) {
    console.error("Error in getTodaysPrayers:", error);
    return { error: "Internal Server Error" }; // or some other error response
  }
}
