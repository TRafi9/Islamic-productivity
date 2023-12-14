const getTodaysPrayers = async (date: string) => {
  if (date) {
    try {
      // const graphResponse = await instance.acquireTokenSilent(request);
      // const token = `Bearer ${graphResponse.accessToken}`;
      const response = await fetch(
        // `api/getTodaysPrayers?bearer=${token}&id=${id}`,
        `api/getTodaysPrayers?date=${date}`,
        {
          method: "GET",
        }
      );
      console.log("awaiting response...");

      const data = await response.json();
      console.log(data);
      return data;
    } catch (error) {
      console.log("error calling api in getTodaysPrayers : ", error);
    }
  } else {
    ("getTodaysPrayers failed");
  }
};

export default getTodaysPrayers;
