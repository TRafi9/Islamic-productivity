import Head from "next/head";
import Image from "next/image";
import { Inter } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { useEffect, useState } from "react";

import getTodaysPrayers from "@/functions/getTodaysPrayers";
import getNextPrayer from "@/functions/getNextPrayer";
import Countdown from "react-countdown";
import getCurrentPrayer from "@/functions/getCurrentPrayer";
import ProductiveStateView from "@/functions/productiveStateView";

import getLastPrayer from "@/functions/getLastPrayer";
import calculateTimeTillRefresh from "@/functions/calculateTimeTillRefresh";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  const [displayType, setDisplayType] = useState("countdown");
  // checkDate is used as a value to check if the currentDate has been changed
  const [checkDate, setCheckDate] = useState("");

  // add cron job to reexecute current date setup and useeffect setup, needs to run every 24 hours
  // update currentDate every 24 hours
  // in the same loop check if the formatted current date isnt the same
  function updateDate() {
    console.log("running update Date, which returns a formattedDate val!");
    var newDate = new Date();
    var year = newDate.getFullYear();
    var month = String(newDate.getMonth() + 1).padStart(2, "0"); // Months are zero-based
    var day = String(newDate.getDate()).padStart(2, "0");
    // const formattedDate = `${year}-${month}-${day}`;
    const formattedDate = `2024-01-14`;
    return formattedDate;
  }

  var currentDate = new Date();
  // // TODO REMOVE for testing as need to move it 1 day up, remove later
  // currentDate.setDate(currentDate.getDate() + 1);
  var year = currentDate.getFullYear();
  var month = String(currentDate.getMonth() + 1).padStart(2, "0"); // Months are zero-based
  var day = String(currentDate.getDate()).padStart(2, "0");
  // Create the formatted date string to match api call date type
  let [formattedDate, setFormattedDate] = useState(`${year}-${month}-${day}`);

  // initial delay before running the refresh via setInterval
  var initialDelay = calculateTimeTillRefresh();
  console.log("initial delay is, ", initialDelay);

  //test this code to see if it works, adjust the time in the calculateTimeTillRefresh code to test
  setTimeout(() => {
    setFormattedDate(updateDate);
    setInterval(() => {
      console.log("interval running, setting formatted date");
      setFormattedDate(updateDate);
    }, 24 * 60 * 60 * 1000);
  }, initialDelay);

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
  //TODO IMPORTANT need to update formattedDate daily/hourly to run this constantly
  useEffect(() => {
    console.log("use effect triggered from formattedDate");
    const fetchData = async () => {
      try {
        const result = await getTodaysPrayers(formattedDate);

        if (result) {
          console.log("results for prayers today...");
          console.log(result);
          setTodaysPrayers(result);
        } else {
          console.log("Results undefined couldnt get todays prayers");
        }
      } catch (error) {
        console.error("Error fetching todays prayers", error);
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

  interface ClosestPrayer {
    name: string;
    time: string;
    difference: number;
  }

  // used to update currentPrayer values when date/time changes
  var [currentPrayer, setCurrentPrayer] = useState<ClosestPrayer | null>(null);

  interface LastPrayer {
    name: string;
    time: string;
  }
  var [lastPrayer, setLastPrayer] = useState<LastPrayer | null>(null);

  // used in a use effect to trigger a rerun of the getNextPrayer function, it runs when the time passes that of the next prayer
  const [nextPrayerTimeActivator, setNextPrayerTimeActivator] = useState<
    number | null
  >(null);
  const [productiveState, setProductiveState] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      if (todaysPrayers != null) {
        console.log(todaysPrayers);
        const nextPrayer = getNextPrayer(todaysPrayers);
        // currPrayer and currentPrayer are different because currPrayer
        // is used in the if statement below, as setting parts of usestate
        // may not be resolved before if statement runs below if it called if(currentPrayer),
        // same for constLastPrayer
        const currPrayer = await getCurrentPrayer(todaysPrayers);
        console.log("curr prayer set!");
        console.log(currPrayer);
        setCurrentPrayer(currPrayer);
        const constLastPrayer = await getLastPrayer(todaysPrayers, currPrayer);
        setLastPrayer(constLastPrayer);
        if (nextPrayer && currPrayer && constLastPrayer) {
          console.log("promises resolved, setting times & names");
          setNextPrayerTime(new Date(nextPrayer.time));
          setNextPrayerName(nextPrayer.name);
          setCurrentPrayerName(currPrayer.name);
          setCurrentPrayerTime(new Date(currPrayer.time));
          setLastPrayerName(constLastPrayer.name);
          setLastPrayerTime(new Date(constLastPrayer.time));
        }
      }
    };
    fetchData();
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
