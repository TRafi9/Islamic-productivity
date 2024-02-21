import Head from "next/head";
import Image from "next/image";
import { Roboto_Mono, Bebas_Neue } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { Inter } from "next/font/google";
import { useEffect, useState } from "react";

import React from "react";
import { Pie } from "react-chartjs-2";

import getTodaysPrayers from "@/functions/getTodaysPrayers";
import getNextPrayer from "@/functions/getNextPrayer";
import Countdown from "react-countdown";
import getCurrentPrayer from "@/functions/getCurrentPrayer";
import ProductiveStateView from "@/functions/productiveStateView";

import getLastPrayer from "@/functions/getLastPrayer";
import calculateTimeTillRefresh from "@/functions/calculateTimeTillRefresh";
import NavbarComponent from "@/components/NavBar";
import { Row, Col } from "react-bootstrap";
import getAllStats from "@/functions/getAllStats";
import PieChartProductiveVal from "@/components/PieChartProductiveVal";

import { Chart as ChartJS, ArcElement, Tooltip, Legend, Title } from "chart.js";
import { Doughnut } from "react-chartjs-2";

ChartJS.register(ArcElement, Tooltip, Legend, Title);

const inter = Inter({ subsets: ["latin"] });
// const [allStats, setAllStats] = useState<string | null>(null);

export default function myStats() {
  // const [dailyStats, setDailyStats] = useState(null);
  // const [allStats, setAllStats] = useState(null);
  const [dailyStats, setDailyStats] = useState<null | Object>(null);
  const [weeklyStats, setWeeklyStats] = useState<null | Object>(null);

  useEffect(() => {
    const fetchData = async () => {
      var myStats = await getAllStats();

      if (myStats) {
        myStats = JSON.parse(myStats);

        const { DailyStats, WeeklyStats } = myStats;
        const formattedDailyStats = { DailyStats };
        const formattedWeeklyStats = { WeeklyStats };
        setDailyStats(formattedDailyStats);
        setWeeklyStats(formattedWeeklyStats);
      } else {
        console.log("no stats from api");
      }
    };
    fetchData();
  }, []);

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
        <PieChartProductiveVal stats={dailyStats} />
        <PieChartProductiveVal stats={weeklyStats} />
      </main>
    </>
  );
}
