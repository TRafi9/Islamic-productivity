import Head from "next/head";
import Image from "next/image";
import { Inter } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { useEffect, useState } from "react";
import { format } from "path";
import { text } from "stream/consumers";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  //TODO link api calls from api page to this
  const [checkDate, setCheckDate] = useState("");
  // update currentDate every 24 hours
  // in the same loop check if the formatted current date isnt the same
  var currentDate = new Date();
  var year = currentDate.getFullYear();
  var month = String(currentDate.getMonth() + 1).padStart(2, "0"); // Months are zero-based
  var day = String(currentDate.getDate()).padStart(2, "0");

  // Create the formatted date string to match api call date type
  const formattedDate = `${year}-${month}-${day}`;

  interface PrayerData {
    Asr: string;
    Dhuhr: string;
    Fajr: string;
    Isha: string;
    Maghrib: string;
  }
  // create the var which is the same structure as the response
  const [todaysPrayers, setTodaysPrayers] = useState<PrayerData>({
    Asr: "",
    Dhuhr: "",
    Fajr: "",
    Isha: "",
    Maghrib: "",
  });
  // check if its first load, or the day has changed, if so call the API to get new results in todaysPrayers
  // need to update formattedDate daily/hourly to run this constantly
  useEffect(() => {
    if (formattedDate !== checkDate || checkDate == null) {
      setCheckDate(formattedDate);
      console.log("formattedDate");
      console.log(formattedDate);
      getTodaysPrayers(formattedDate).then((result) => {
        console.log("Checking if result exists:", result);

        if (result) {
          console.log("Before SetState:", todaysPrayers);
          setTodaysPrayers(result);
          console.log("After SetState:", todaysPrayers);
          console.log("Asr value:", result.Asr);
        } else {
          console.log("Result is undefined or null");
        }
      });
    }
  }, [formattedDate]);

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
        console.log("responseeee");
        console.log(data);
        return data;
      } catch (error) {
        console.log("error calling api in getTodaysPrayers : ", error);
      }
    } else {
      //TODO update this call
      ("getGranularDirectReport function did not recieve an id - undefined");
    }
  };

  function findClosestAheadTime(prayerTimes: PrayerData): string {
    const now = new Date();

    // Convert the time strings to Date objects
    const timeObjects = Object.entries(prayerTimes).map(([prayer, time]) => ({
      prayer,
      time: new Date(time),
    }));

    // Filter out times that are in the past
    const futureTimes = timeObjects.filter(({ time }) => time > now);

    // If there are no future times, consider the first time as the closest ahead
    if (futureTimes.length === 0) {
      return timeObjects[0].prayer;
    }

    // Find the closest ahead time
    const closestAheadTime = futureTimes.reduce((closest, current) => {
      const closestDifference = closest.time.getTime() - now.getTime();
      const currentDifference = current.time.getTime() - now.getTime();

      return currentDifference < closestDifference ? current : closest;
    });

    return closestAheadTime.prayer;
  }

  const closestAheadTime = findClosestAheadTime(todaysPrayers);
  const closestTimeValue =
    todaysPrayers[closestAheadTime as keyof typeof todaysPrayers];

  // Step 1: Create state variable for countdown
  const [countdown, setCountdown] = useState<number | null>(null);
  useEffect(() => {
    // Convert closestTimeValue to a Date object
    const closestTimeDate = closestTimeValue
      ? new Date(closestTimeValue)
      : null;

    // Step 2: Update countdown every second
    const interval = setInterval(() => {
      // Step 3: Calculate time left until closestTimeDate
      if (closestTimeDate) {
        const now = new Date();
        const timeLeftInSeconds = Math.floor(
          (closestTimeDate.getTime() - now.getTime()) / 1000
        );

        // Step 4: Update countdown state variable
        setCountdown(timeLeftInSeconds);

        // Optionally: You can clear the interval if the countdown reaches zero
        if (timeLeftInSeconds <= 0) {
          clearInterval(interval);
        }
      }
    }, 1000);

    // Clear the interval when the component unmounts
    return () => clearInterval(interval);
  }, [closestTimeValue]);

  return (
    <>
      <Head>
        <title>The productive muslim</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={`${styles.main} ${inter.className}`}>
        <div>
          {/* Step 5: Display the countdown in your JSX */}
          <h1>The productive muslim! {todaysPrayers.Asr}</h1>
          {closestTimeValue && (
            <h2>
              Time left until {closestAheadTime} prayer:{" "}
              {formatCountdown(countdown)}
            </h2>
          )}

          <h2> Measuring your productivity one step at a time</h2>
        </div>
      </main>
    </>
  );
}

// Helper function to format seconds into HH:MM:SS
function formatCountdown(seconds: number | null): string {
  if (seconds === null) {
    return "";
  }

  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;

  return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(
    2,
    "0"
  )}:${String(remainingSeconds).padStart(2, "0")}`;
}
