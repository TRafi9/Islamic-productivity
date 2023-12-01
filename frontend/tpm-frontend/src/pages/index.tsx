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

  // Create the formatted date string
  const formattedDate = `${year}-${month}-${day}`;

  useEffect(() => {
    if (formattedDate !== checkDate || checkDate == null) {
      setCheckDate(formattedDate);
      console.log("formattedDate");
      console.log(formattedDate);
      getTodaysPrayers(formattedDate).then((result) => {
        console.log(result);
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
        const textResponse = await response.text();
        console.log("API Response:", textResponse);

        // const data = JSON.parse(textResponse);
        // const data = await response.json();

        return textResponse;
      } catch (error) {
        console.log("error calling api in getTodaysPrayers : ", error);
      }
    } else {
      ("getGranularDirectReport function did not recieve an id - undefined");
    }
  };

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
          <h1>The productive muslim</h1>
          <h2> Measuring your productivity one step at a time</h2>
        </div>
      </main>
    </>
  );
}
