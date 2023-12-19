// use this to post productivity value to db storing data

interface RequestObject {
  query: {
    date: string;
  };
}
import { NextApiResponse } from "next";

export default async function postProductivityValue(
  req: RequestObject,
  res: NextApiResponse
) {
  try {
    const { date } = req.query;
    console.log("date in api");
    console.log(date);
    const response = await fetch(
      // needs to be updated to something else probably
      `http://localhost:8080/api/v1/sendUserInput/${date}`,
      { method: "GET" }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    console.log(typeof data);
    console.log(data);
    res.status(200).json(data);
    return data;
  } catch (error) {
    console.error("Error in getTodaysPrayers:", error);
    return { error: "Internal Server Error" }; // or some other error response
  }
}
