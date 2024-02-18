import { Inter } from "next/font/google";

import React from "react";
import { Pie } from "react-chartjs-2";

import { Chart as ChartJS, ArcElement, Tooltip, Legend, Title } from "chart.js";

ChartJS.register(ArcElement, Tooltip, Legend, Title);

const inter = Inter({ subsets: ["latin"] });

function PieChartProductiveVal(props) {
  const dailyStats = props.dailyStats;
  const labels = [];
  const dataValues = [];
  console.log(dailyStats);
  if (dailyStats) {
    console.log("daily stats exist");
    Object.entries(dailyStats).forEach(([key, value]) => {
      labels.push(key);
      dataValues.push(value);
    });

    const data = {
      labels: labels,
      datasets: [
        {
          //  label: "My First Dataset",
          data: dataValues,
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
        <Pie data={data} />
      </div>
    );
  } else {
    return <></>;
  }
}
export default PieChartProductiveVal;
