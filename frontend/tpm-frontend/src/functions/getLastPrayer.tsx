interface PrayerData {
  Asr: string;
  Dhuhr: string;
  Fajr: string;
  Isha: string;
  Maghrib: string;
}

interface CurrentPrayer {
  name: string;
  time: string;
  difference: number;
}

interface LastPrayer {
  name: string;
  time: string;
}

const getLastPrayer = (
  todaysPrayers: PrayerData,
  currentPrayer: CurrentPrayer | null
): LastPrayer | null => {
  var lastPrayer: LastPrayer = {
    name: "",
    time: "",
  };
  if (currentPrayer) {
    var CurrentPrayerName: string = currentPrayer.name;

    switch (CurrentPrayerName) {
      case "Fajr":
        break;
      case "Dhuhr":
        lastPrayer = {
          name: "Fajr",
          time: todaysPrayers.Fajr,
        };
        break;
      case "Asr":
        lastPrayer = {
          name: "Dhuhr",
          time: todaysPrayers.Dhuhr,
        };
        break;
      case "Maghrib":
        lastPrayer = {
          name: "Asr",
          time: todaysPrayers.Asr,
        };
        break;
      case "Isha":
        lastPrayer = {
          name: "Maghrib",
          time: todaysPrayers.Maghrib,
        };
        break;
    }

    return lastPrayer;
  } else {
    console.log(
      "current prayer missing, so cannot use switch case to get last prayer"
    );
    return null;
  }
};

export default getLastPrayer;
