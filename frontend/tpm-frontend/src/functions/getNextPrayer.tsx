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

// this function returns next prayer in closestPrayer interface
const getNextPrayer = (todaysPrayers: PrayerData): ClosestPrayer | null => {
  const currTime = new Date();
  // TODO for testing as need to move it 1 day up, remove later
  // currTime.setDate(currTime.getDate() + 1);

  const filteredPrayerObj: Record<string, ClosestPrayer> = Object.entries(
    todaysPrayers
  ).reduce((acc, [key, value]) => {
    const prayerTime = new Date(value);

    if (prayerTime > currTime) {
      const difference = prayerTime.getTime() - currTime.getTime();
      acc[key] = { name: key, time: value, difference };
    }

    return acc;
  }, {} as Record<string, ClosestPrayer>);

  // If there are no upcoming prayers, return null
  if (Object.keys(filteredPrayerObj).length === 0) {
    console.log(
      "lmao- get next prayer returned null, comparing today to tomorrow"
    );
    return null;
  }

  const closestPrayer = Object.values(filteredPrayerObj).reduce(
    (closest, current) => {
      return current.difference < closest.difference ? current : closest;
    },
    { name: "", time: "", difference: Infinity }
  );

  return closestPrayer;
};

export default getNextPrayer;
