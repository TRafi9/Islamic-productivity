const getPieChartTitle = (stats: Object) => {
  if (stats) {
    const innerStatsKey = Object.keys(stats)[0];
    console.log(innerStatsKey);

    switch (innerStatsKey) {
      case "DailyStats":
        return "Daily stats";
      case "WeeklyStats":
        return "Weekly stats";
      case "MonthlyStats":
        return "Monthly stats";
      default:
        return "Stats";
    }
  } else {
    // Handle the case when stats is undefined or null
    return "No stats available";
  }
};

export default getPieChartTitle;
