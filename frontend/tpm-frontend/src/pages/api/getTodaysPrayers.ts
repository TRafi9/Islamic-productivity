interface RequestObject {
  query: {
    date: string;
  };
}
import { NextApiResponse } from "next";

export default async function getTodaysPrayers(
  req: RequestObject,
  res: NextApiResponse
) {
  try {
    const { date } = req.query;
    //   const token = await graphToken(bearer);
    //   const headers = {
    //     Authorization: `Bearer ${token}`,
    //   };
    console.log("date in api");
    console.log(date);
    const response = await fetch(
      // needs to be updated to something else probably
      `http://tpm-backend:8080/api/v1/getPrayerTimes/${date}`,
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

// export default async function directReportsAPI(req, res) {
//   try {
//     const { bearer } = req.query;
//     const token = await graphToken(bearer);

//     const headers = {
//       Authorization: `Bearer ${token}`,
//     };

//     const response = await fetch(
//       "https://graph.microsoft.com/v1.0/me/directReports",
//       { method: "GET", headers }
//     );

//     if (!response.ok) {
//       throw new Error(`Request failed with status ${response.status}`);
//     }

//     const graphResponse = await response.json();
//     res.status(200).json(graphResponse);
//   } catch (error) {
//     console.error(error);
//     res.status(500).send("Error processing the request");
//   }
// }
