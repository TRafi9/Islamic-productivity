import Head from "next/head";
import Image from "next/image";
import { Roboto_Mono, Bebas_Neue } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { useEffect, useState } from "react";

import getTodaysPrayers from "@/functions/getTodaysPrayers";
import getNextPrayer from "@/functions/getNextPrayer";
import Countdown from "react-countdown";
import getCurrentPrayer from "@/functions/getCurrentPrayer";
import ProductiveStateView from "@/functions/productiveStateView";

import getLastPrayer from "@/functions/getLastPrayer";
import calculateTimeTillRefresh from "@/functions/calculateTimeTillRefresh";
import NavbarComponent from "@/components/NavBar";
import { Row, Col } from "react-bootstrap";
import AfterIshaView from "@/functions/afterIshaView";

const inter = Roboto_Mono({
  weight: "400",
  subsets: ["latin"],
});
const Bebas = Bebas_Neue({
  weight: "400",
  subsets: ["latin"],
});

export default function Home() {
  const [displayType, setDisplayType] = useState("countdown");
  const [checkDate, setCheckDate] = useState(null);

  function updateDate() {
    var newDate = new Date();
    var year = newDate.getFullYear();
    var month = String(newDate.getMonth() + 1).padStart(2, "0");
    var day = String(newDate.getDate()).padStart(2, "0");
    const formattedDate = `${year}-${month}-${day}`;
    return formattedDate;
  }

  var currentDate = new Date();
  var year = currentDate.getFullYear();
  var month = String(currentDate.getMonth() + 1).padStart(2, "0");
  var day = String(currentDate.getDate()).padStart(2, "0");
  let [formattedDate, setFormattedDate] = useState(`${year}-${month}-${day}`);

  var initialDelay = calculateTimeTillRefresh();
  setTimeout(() => {
    setFormattedDate(updateDate);
    setInterval(() => {
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
  const [todaysPrayers, setTodaysPrayers] = useState<PrayerData>({
    Asr: "",
    Dhuhr: "",
    Fajr: "",
    Isha: "",
    Maghrib: "",
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        // const result = await getTodaysPrayers(formattedDate);
        const result = {
          Asr: "2024-02-16T14:43:00Z",
          Dhuhr: "2024-02-16T11:40:00Z",
          Fajr: "2024-02-16T08:23:00Z",
          Isha: "2024-02-16T11:42:00Z",
          Maghrib: "2024-02-16T14:10:00Z",
        };

        if (result) {
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

  var [currentPrayer, setCurrentPrayer] = useState<ClosestPrayer | null>(null);

  interface LastPrayer {
    name: string;
    time: string;
  }
  var [lastPrayer, setLastPrayer] = useState<LastPrayer | null>(null);

  const [nextPrayerTimeActivator, setNextPrayerTimeActivator] = useState<
    number | null
  >(null);
  const [productiveState, setProductiveState] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      const nextPrayer = getNextPrayer(todaysPrayers);
      if (todaysPrayers != null) {
        const currPrayer = await getCurrentPrayer(todaysPrayers);
        setCurrentPrayer(currPrayer);
        const constLastPrayer = await getLastPrayer(todaysPrayers, currPrayer);
        setLastPrayer(constLastPrayer);
        if (nextPrayer && currPrayer && constLastPrayer) {
          setNextPrayerTime(new Date(nextPrayer.time));
          setNextPrayerName(nextPrayer.name);
          setCurrentPrayerName(currPrayer.name);
          setCurrentPrayerTime(new Date(currPrayer.time));
          setLastPrayerName(constLastPrayer.name);
          setLastPrayerTime(new Date(constLastPrayer.time));
        }
      }
      if (!nextPrayer) {
        // no nextPrayer value means that isha is the currentPrayer
        setCurrentPrayerName("Isha");
        setCurrentPrayerTime(new Date(todaysPrayers.Isha));
        setLastPrayerName("Maghrib");
        setLastPrayerTime(new Date(todaysPrayers.Maghrib));
      }
      console.log("next prayer is...");
      console.log(nextPrayer);
    };
    fetchData();
  }, [todaysPrayers, nextPrayerTimeActivator, productiveState]);

  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    if (nextPrayerTime && nextPrayerName) {
      intervalId = setInterval(
        () => updateNextPrayer(nextPrayerTime, nextPrayerName),
        1000
      );
    }

    return () => clearInterval(intervalId);
  }, [nextPrayerTime, nextPrayerName]);

  function updateNextPrayer(nextPrayerTime: Date, nextPrayerName: string) {
    const timer = new Date();

    // if (timer > nextPrayerTime && nextPrayerName == "Isha") {
    //   setNextPrayerName("AFTER ISHA");
    //   setDisplayType("after isha");
    // } else
    if (timer > nextPrayerTime) {
      setProductiveState(true);
      setNextPrayerTimeActivator(1);
    }
  }

  const countdownKey = nextPrayerTime ? nextPrayerTime.toString() : null;

  return (
    <>
      <Head>
        <link
          rel="stylesheet"
          href="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/css/bootstrap.min.css"
          integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
          crossOrigin="anonymous"
        />
      </Head>
      <main className={`${styles.main} ${inter.className}`}>
        <NavbarComponent />
        {displayType == "countdown" && productiveState == false && (
          <div>
            <Row>
              <br></br>
              <br></br>
            </Row>
            <Row>
              <p className="display-4 col-center">
                Time left till {nextPrayerName} is
              </p>
            </Row>
            <div>
              <p
                className={`${Bebas.className} text-center-p countdown-display`}
              >
                {nextPrayerTime !== null && (
                  <Countdown
                    key={countdownKey}
                    className="countdown-display"
                    date={nextPrayerTime}
                    daysInHours={true}
                  />
                )}
              </p>
            </div>
          </div>
        )}
        {displayType == "after isha" && (
          <AfterIshaView setProductiveState={setProductiveState} />
        )}
        {productiveState == true && (
          <ProductiveStateView
            setProductiveState={setProductiveState}
            currentPrayerName={currentPrayerName}
            currentPrayerTime={currentPrayerTime}
            lastPrayerName={lastPrayerName}
            lastPrayerTime={lastPrayerTime}
            nextPrayerName={nextPrayerName}
            setDisplayType={setDisplayType}
          />
        )}
      </main>
    </>
  );
}
