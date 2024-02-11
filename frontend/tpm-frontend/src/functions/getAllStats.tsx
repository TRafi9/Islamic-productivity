const getAllStats = async () => {
  try {
    const response = await fetch(`api/getAllStats`, {
      method: "GET",
      credentials: "include",
    });
    console.log("awaiting response...");

    const data = await response.json();

    return data;
  } catch (error) {
    console.log("error calling api in getAllStats api : ", error);
  }
};

export default getAllStats;
