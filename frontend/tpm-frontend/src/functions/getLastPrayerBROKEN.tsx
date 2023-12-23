interface PrayerData {
  Asr: string;
  Dhuhr: string;
  Fajr: string;
  Isha: string;
  Maghrib: string;
}

interface ClosestPrayer {
  name: string;
  time: string;
  difference: number;
}

const getLastPrayer = (todaysPrayers: PrayerData): ClosestPrayer => {
  const currTime = new Date();

  const filteredPrayerArray: ClosestPrayer[] = Object.entries(
    todaysPrayers
  ).reduce((acc, [key, value]) => {
    const prayerTime = new Date(value);

    if (prayerTime <= currTime) {
      const difference = currTime.getTime() - prayerTime.getTime();
      acc.push({ name: key, time: value, difference });
    }

    return acc;
  }, [] as ClosestPrayer[]);

  // Sort the array based on the difference in ascending order
  filteredPrayerArray.sort((a, b) => a.difference - b.difference);
  console.log(filteredPrayerArray);
  // Return the last prayer that isnt the current prayer going on
  return filteredPrayerArray.slice(-2)[0];
};

export default getLastPrayerOLD;
