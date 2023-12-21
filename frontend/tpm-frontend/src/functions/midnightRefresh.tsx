import { useEffect } from "react";

const MidnightRefresh = () => {
  useEffect(() => {
    const refreshAtMidnight = () => {
      const now: number = new Date().getTime(); // Use getTime to get a timestamp in milliseconds
      const midnight: Date = new Date(now);
      midnight.setHours(24, 0, 0, 0); // Set to midnight

      const timeUntilMidnight: number = midnight.getTime() - now;

      setTimeout(() => {
        // Reload the page when it's midnight
        window.location.reload();
      }, timeUntilMidnight);
    };

    // Call the function immediately to set up the first refresh
    refreshAtMidnight();

    // Set up recurring interval to refresh the page every 24 hours
    const refreshInterval = 24 * 60 * 60 * 1000; // 24 hours in milliseconds
    const intervalId = setInterval(refreshAtMidnight, refreshInterval);

    // Calculate the time until the next midnight and refresh accordingly
    const now: number = new Date().getTime();
    const timeUntilNextMidnight: number =
      24 * 60 * 60 * 1000 - (now % (24 * 60 * 60 * 1000));
    setTimeout(() => {
      refreshAtMidnight();
      // Refresh every 24 hours thereafter
      setInterval(refreshAtMidnight, refreshInterval);
    }, timeUntilNextMidnight);

    // Cleanup function to clear the interval when the component is unmounted
    return () => clearInterval(intervalId);
  }, []); // Empty dependency array ensures that this effect runs once on mount
};

export default MidnightRefresh;
