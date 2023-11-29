export default async function getTodayPrayers(req) {
  const { date } = req.query;
  //   const token = await graphToken(bearer);
  //   const headers = {
  //     Authorization: `Bearer ${token}`,
  //   };

  const response = await fetch(
    // needs to be updated to something else probably
    "localhost:8080/api/v1/getPrayerTimes/2023-11-11",
    // { method: "GET", headers }
    { method: "GET" }
  );
}

export default async function directReportsAPI(req, res) {
  try {
    const { bearer } = req.query;
    const token = await graphToken(bearer);

    const headers = {
      Authorization: `Bearer ${token}`,
    };

    const response = await fetch(
      "https://graph.microsoft.com/v1.0/me/directReports",
      { method: "GET", headers }
    );

    if (!response.ok) {
      throw new Error(`Request failed with status ${response.status}`);
    }

    const graphResponse = await response.json();
    res.status(200).json(graphResponse);
  } catch (error) {
    console.error(error);
    res.status(500).send("Error processing the request");
  }
}
