const calculateTimeTillRefresh = () => {
  // This function calculates how long there is left till the application needs to refresh
  // the call to the getTodaysPrayers API
  const timeNow = new Date();
  const targetTime = new Date(timeNow);
  // time is set to 00:05
  targetTime.setHours(0, 5, 0, 0);

  let initialDelay: number = targetTime.getTime() - timeNow.getTime();
  //                            0:05                    0:20

  if (initialDelay < 0) {
    // add 24 hours to the time (this is in milliseconds)
    initialDelay += 24 * 60 * 60 * 1000;
  }
  return initialDelay;
};

// // Calculate the initial delay until 00:05
// const initialDelay = calculateInitialDelay();

// // Schedule the task to run every day at 00:05
// setInterval(dailyTask, 24 * 60 * 60 * 1000); // Repeat every 24 hours
// setTimeout(() => {
//   // Initial run after the calculated delay
//   dailyTask();

//   // Set up the repeating interval
//   setInterval(dailyTask, 24 * 60 * 60 * 1000);
// }, initialDelay);
