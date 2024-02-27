import Head from "next/head";
import { Roboto_Mono } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { useEffect, useState } from "react";

import React from "react";

import NavbarComponent from "@/components/NavBar";
import { Row, Col } from "react-bootstrap";
import getAllStats from "@/functions/getAllStats";
import PieChartProductiveVal from "@/components/PieChartProductiveVal";

import { Chart as ChartJS, ArcElement, Tooltip, Legend, Title } from "chart.js";

ChartJS.register(ArcElement, Tooltip, Legend, Title);

const roboto = Roboto_Mono({
  weight: "400",
  subsets: ["latin"],
});

export default function MyStats() {
  const [dailyStats, setDailyStats] = useState<null | Object>(null);
  const [weeklyStats, setWeeklyStats] = useState<null | Object>(null);
  const [monthlyStats, setMonthlyStats] = useState<null | Object>(null);

  useEffect(() => {
    const fetchData = async () => {
      var myStats = await getAllStats();

      if (myStats) {
        myStats = JSON.parse(myStats);

        const { DailyStats, WeeklyStats, MonthlyStats } = myStats;
        const formattedDailyStats = { DailyStats };
        const formattedWeeklyStats = { WeeklyStats };
        const formattedMonthlyDtats = { MonthlyStats };
        setDailyStats(formattedDailyStats);
        setWeeklyStats(formattedWeeklyStats);
        setMonthlyStats(formattedMonthlyDtats);
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
      <main className={`${styles.main} ${roboto.className}`}>
        <NavbarComponent />

        <Row>
          <Col>
            <PieChartProductiveVal stats={dailyStats} />
          </Col>
          <Col>
            <PieChartProductiveVal stats={weeklyStats} />
          </Col>
        </Row>
        <Row>
          <PieChartProductiveVal stats={monthlyStats} />
        </Row>
      </main>
    </>
  );
}
