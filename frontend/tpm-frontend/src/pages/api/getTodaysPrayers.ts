import { parseCookies } from "nookies";

interface RequestObject {
  query: {
    date: string;
  };
  headers: {
    cookie?: string;
  };
}
import { NextApiResponse } from "next";

export default async function getTodaysPrayers(
  req: RequestObject,
  res: NextApiResponse
) {
  const { date } = req.query;
  const cookies = parseCookies({ req });

  // Extract the jwt cookie
  const jwtCookie = cookies.jwt;
  try {
    const response = await fetch(
      // needs to be updated to something else probably
      // `http://tpm-backend:8080/api/v1/getPrayerTimes/${date}`,
      `http://tpm-backend:80/api/v1/restricted/getPrayerTimes/${date}`,
      {
        method: "GET",
        credentials: "include",
        headers: {
          // Include jwt cookie in the headers of the outgoing request
          Authorization: jwtCookie,
        },
      }
    );

    if (!response.ok) {
      const errorData = await response.text();
      console.error("Error response body:", errorData);
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    console.log(typeof data);
    console.log(data);
    res.status(200).json(data);
    return data;
  } catch (error) {
    console.error("Error fetching prayer times:", error);
    res.status(500).json({ error: "Internal Server Error" });
  }
}
