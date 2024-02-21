import { Inter } from "next/font/google";

import React from "react";
import { Pie } from "react-chartjs-2";

import { Chart as ChartJS, ArcElement, Tooltip, Legend, Title } from "chart.js";
import getPieChartTitle from "@/functions/getPieChartTitle";

ChartJS.register(ArcElement, Tooltip, Legend, Title);

const inter = Inter({ subsets: ["latin"] });

function PieChartProductiveVal(props) {
  const stats = props.stats;
  const labels = [];
  const dataValues = [];

  if (stats) {
    const title = getPieChartTitle(stats);
    const innerStatsKey = Object.keys(stats)[0];
    console.log(innerStatsKey);
    const innerObjectStats = stats[innerStatsKey];
    console.log("daily stats exist");
    Object.entries(innerObjectStats).forEach(([key, value]) => {
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

    return (
      <div>
        <h1 className="text-center-p">{title}</h1>
        <Pie data={data} />
      </div>
    );
  } else {
    return <></>;
  }
}
export default PieChartProductiveVal;
