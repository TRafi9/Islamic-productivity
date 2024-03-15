const getPieChartTitle = (stats: Object) => {
  if (stats) {
    const innerStatsKey = Object.keys(stats)[0];
    console.log(innerStatsKey);

    switch (innerStatsKey) {
      case "DailyStats":
        return "Daily";
      case "WeeklyStats":
        return "Weekly";
      case "MonthlyStats":
        return "Monthly";
      default:
        return "Stats";
    }
  } else {
    // Handle the case when stats is undefined or null
    return "No stats available";
  }
};

export default getPieChartTitle;
