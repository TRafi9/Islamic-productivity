const getTodaysPrayers = async (date: string, jwt: string | null) => {
  if (date) {
    try {
      const response = await fetch(`api/getTodaysPrayers?date=${date}`, {
        method: "GET",
        credentials: "include",
      });
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
