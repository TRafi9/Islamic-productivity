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

    console.log(dataValues);
    var hasData = true;
    if (
      dataValues.length === 2 &&
      dataValues.every((element) => element === 0)
    ) {
      console.log("data values empty");
      hasData = false;
    }

    const data = {
      labels: labels,
      datasets: [
        {
          //  label: "My First Dataset",
          data: dataValues,
          backgroundColor: ["#36A2EB", "#FF6384"],
          hoverOffset: 4,
          borderWidth: 5,
          borderColor: "#fff",
        },
      ],
    };

    return (
      <div>
        <h1 className="text-center-p">{title}</h1>
        {hasData ? (
          <Pie data={data} />
        ) : (
          <div
            style={{
              width: "280px",
              height: "280px",
              border: "8px solid #888",
              borderRadius: "50%",
            }}
          ></div>
        )}
      </div>
    );
  } else {
    return <></>;
  }
}
export default PieChartProductiveVal;
