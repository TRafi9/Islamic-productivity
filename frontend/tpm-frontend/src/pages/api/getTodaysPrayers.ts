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

  const response = await fetch(
    // needs to be updated to something else probably
    // `http://tpm-backend:8080/api/v1/getPrayerTimes/${date}`,
    `http://localhost:8080/api/v1/getPrayerTimes/${date}`,
    {
      method: "GET",
      credentials: "include",
      headers: {
        // Include jwt cookie in the headers of the outgoing request
        Cookie: jwtCookie,
      },
    }
  );

  if (!response.ok) {
    throw new Error(`HTTP error! Status: ${response.status}`);
  }

  const data = await response.json();
  console.log(typeof data);
  console.log(data);
  res.status(200).json(data);
  return data;
}
