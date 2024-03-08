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

export default async function getAllStats(
  req: RequestObject,
  res: NextApiResponse
) {
  const cookies = parseCookies({ req });

  // Extract the jwt cookie
  const jwtCookie = cookies.jwt;

  const response = await fetch(
    `http://tpm-backend:80/api/v1/restricted/getAllStats`,
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
    throw new Error(`HTTP error! Status: ${response.status}`);
  }

  const data = await response.json();
  console.log(typeof data);
  console.log(data);
  res.status(200).json(data);
  return data;
}
