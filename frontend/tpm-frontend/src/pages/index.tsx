import Head from "next/head";
import Image from "next/image";
import { Inter } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { useEffect, useState } from "react";
import { format } from "path";
import { text } from "stream/consumers";
import { Button } from "react-bootstrap";
import getTodaysPrayers from "@/functions/getTodaysPrayers";
import getNextPrayer from "@/functions/getNextPrayer";
import next from "next";
import Countdown from "react-countdown";
import getCurrentPrayer from "@/functions/getCurrentPrayer";
import ProductiveStateView from "@/functions/productiveStateView";

import getLastPrayer from "@/functions/getLastPrayer";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  const [displayType, setDisplayType] = useState("countdown");
  // checkDate is used as a value to check if the currentDate has been changed
  const [checkDate, setCheckDate] = useState("");

  // add cron job to reexecute current date setup and useeffect setup, needs to run every 24 hours
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

  const [prayersLeftInDay, setPrayersLeftInDay] = useState<Record<
    string,
    string
  > | null>(null);
  // check if its first load, or the day has changed, if so call the API to get new results in todaysPrayers
  //TODO IMPORTANT need to update formattedDate daily/hourly to run this constantly
  useEffect(() => {
    const fetchData = async () => {
      if (formattedDate !== checkDate || checkDate == null) {
        setCheckDate(formattedDate);
        try {
          const result = await getTodaysPrayers(formattedDate);

          if (result) {
            setTodaysPrayers({
              Asr: "2023-12-27T15:42:00Z",
              Dhuhr: "2023-12-27T13:54:00Z",
              Fajr: "2023-12-27T06:56:00Z",
              Isha: "2023-12-27T21:32:00Z",
              Maghrib: "2023-12-27T19:35:00Z",
            });
          } else {
            console.log("Results undefined couldnt get todays prayers");
          }
        } catch (error) {
          console.error("Error fetching todays prayers", error);
        }
      }
    };
    fetchData();
  }, [formattedDate]);

  const [nextPrayerName, setNextPrayerName] = useState<string | null>(null);
  const [nextPrayerTime, setNextPrayerTime] = useState<Date | null>(null);
  const [currentPrayerName, setCurrentPrayerName] = useState<string | null>(
    null
  );
  const [currentPrayerTime, setCurrentPrayerTime] = useState<Date | null>(null);
  const [lastPrayerName, setLastPrayerName] = useState<string | null>(null);
  const [lastPrayerTime, setLastPrayerTime] = useState<Date | null>(null);

  // used in a use effect to trigger a rerun of the getNextPrayer function, it runs when the time passes that of the next prayer
  const [nextPrayerTimeActivator, setNextPrayerTimeActivator] = useState<
    number | null
  >(null);
  const [productiveState, setProductiveState] = useState(false);

  useEffect(() => {
    if (todaysPrayers != null) {
      const nextPrayer = getNextPrayer(todaysPrayers);
      const currentPrayer = getCurrentPrayer(todaysPrayers);
      const lastPrayer = getLastPrayer(todaysPrayers, currentPrayer);
      console.log("last prayer...");
      console.log(lastPrayer);
      console.log("current prayer...");
      console.log(currentPrayer);

      if (nextPrayer && currentPrayer && lastPrayer) {
        setNextPrayerTime(new Date(nextPrayer.time));
        setNextPrayerName(nextPrayer.name);
        setCurrentPrayerName(currentPrayer.name);
        setCurrentPrayerTime(new Date(currentPrayer.time));
        setLastPrayerName(lastPrayer.name);
        setLastPrayerTime(new Date(lastPrayer.time));
      }
    }
  }, [todaysPrayers, nextPrayerTimeActivator]);

  // if a nextPrayerTime exists (should do after first load), start timer to see when it goes past nextPrayerTime
  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    if (nextPrayerTime && nextPrayerName) {
      // Pass a function reference, not an invocation
      intervalId = setInterval(
        () => updateNextPrayer(nextPrayerTime, nextPrayerName),
        1000
      );
    }

    // Cleanup function to clear the interval when the component unmounts
    return () => clearInterval(intervalId);
  }, [nextPrayerTime, nextPrayerName]);

  // if timer in this function goes past nextPrayerTime, it will know and will hit the activator which will rerun the useeffect to call a new prayer time
  function updateNextPrayer(nextPrayerTime: Date, nextPrayerName: string) {
    const timer = new Date();
    console.log("timer running...");
    if (timer > nextPrayerTime && nextPrayerName == "Isha") {
      setNextPrayerName("AFTER ISHA");
      setDisplayType("after isha");
      // timer is past prayer time, show productive state
    } else if (timer > nextPrayerTime) {
      setProductiveState(true);

      setNextPrayerTimeActivator(1);
    }
  }

  const countdownKey = nextPrayerTime ? nextPrayerTime.toString() : null;

  return (
    <>
      <Head>
        <title>The productive muslim</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={`${styles.main} ${inter.className}`}>
        {displayType == "countdown" && productiveState == false && (
          <div>
            <p>
              {" "}
              the next Prayer is {nextPrayerName} at {String(nextPrayerTime)}
            </p>
            <br></br>
            <p> Time left till {nextPrayerName} is</p>
            <p>
              {nextPrayerTime !== null && (
                <Countdown key={countdownKey} date={nextPrayerTime} />
              )}
            </p>
          </div>
        )}
        {displayType == "after isha" && (
          <div>
            <p> after isha, come back tomorrow</p>
          </div>
        )}
        {productiveState == true && (
          <ProductiveStateView
            setProductiveState={setProductiveState}
            currentPrayerName={currentPrayerName}
            currentPrayerTime={currentPrayerTime}
            lastPrayerName={lastPrayerName}
            lastPrayerTime={lastPrayerTime}
          />
        )}
      </main>
    </>
  );
}
