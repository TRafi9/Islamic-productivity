import Head from "next/head";
import Image from "next/image";
import { Roboto_Mono, Bebas_Neue } from "next/font/google";
import styles from "@/styles/Home.module.css";
import { Inter } from "next/font/google";
import { useEffect } from "react";

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

import { Chart as ChartJS, ArcElement, Tooltip, Legend, Title } from "chart.js";
import { Doughnut } from "react-chartjs-2";

ChartJS.register(ArcElement, Tooltip, Legend, Title);

const inter = Inter({ subsets: ["latin"] });
// const [allStats, setAllStats] = useState<string | null>(null);

const MyPieChart = () => {
  const data = {
    labels: ["Productive", "Unproductive"],
    datasets: [
      {
        label: "My First Dataset",
        data: [12, 19],
        backgroundColor: ["#36A2EB", "#FF6384"], // Use custom colors
        hoverOffset: 4,
        borderWidth: 1, // Add border width
        borderColor: "#fff", // Add border color
      },
    ],
  };

  const options = {
    plugins: {
      title: {
        display: true,
        // text: "Your total",
        font: {
          size: 16,
        },
      },
      legend: {
        display: true,
        position: "bottom",
      },
      tooltip: {
        enabled: true,
      },
    },
    animation: {
      animateRotate: true, // Enable rotation animation
      animateScale: true, // Enable scaling animation
    },
    rotation: -0.5 * Math.PI, // Rotate the chart to a specific angle
    responsive: true, // Make the chart responsive
  };

  return (
    <div>
      <h2>Todays productivity</h2>
      <Pie data={data} options={options} />
    </div>
  );
};

export default function myStats() {
  useEffect(() => {
    const fetchData = async () => {
      const myStats = await getAllStats();

      if (myStats) {
        console.log(myStats);
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
        <MyPieChart></MyPieChart>
        <div>
          <h1>my stats </h1>
        </div>
      </main>
    </>
  );
}
