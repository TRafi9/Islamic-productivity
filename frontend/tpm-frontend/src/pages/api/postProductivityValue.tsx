// use this to post productivity value to db storing data

import { NextApiResponse, NextApiRequest } from "next";

export default async function postProductivityValue(
  req: NextApiRequest,
  res: NextApiResponse
) {
  try {
    const data = JSON.parse(req.body);
    console.log("parsed json data");
    console.log(data);

    const response = await fetch(
      // needs to be updated to something else probably
      `http://localhost:8080/api/v1/userData`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data), // include date and data in the request body
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const dataRes = await response.json();
    console.log(typeof dataRes);
    console.log(dataRes);
    res.status(200).json(dataRes);
  } catch (error) {
    console.error("Error in getTodaysPrayers:", error);
    return { error: "Internal Server Error" }; // or some other error response
  }
}
