import Button from "react-bootstrap/Button";
import AfterIshaView from "./afterIshaView";

// need to send data to handleSendData from getCurrentPrayer + next prayer
const ProductiveStateView = (props: any) => {
  // set productiveState back to false after person has clicked button
  const handleSendData = async (data: any) => {
    try {
      const response = await fetch(`api/postProductivityValue`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (response.ok) {
        // Request was successful
        const responseData = await response.json();
        console.log("API Response:", responseData);
      } else {
        // Handle errors
        console.error("Error:", response.statusText);
      }
    } catch (error) {
      console.error("Error:");
    }

    if (!props.nextPrayer) {
      props.setDisplayType("after isha");
      props.setProductiveState(false);
    }
  };

  async function sendData(value: boolean) {
    const data = {
      currentPrayerName: props.currentPrayerName,
      currentPrayerTime: props.currentPrayerTime,
      lastPrayerName: props.lastPrayerName,
      lastPrayerTime: props.lastPrayerTime,
      productiveValue: value,
    };
    const statusResponse = await handleSendData(data);
    console.log(statusResponse);

    if (!props.nextPrayerName) {
      props.setDisplayType("after isha");
      props.setProductiveState(false);
    }
    props.setProductiveState(false);
  }

  return (
    <div>
      <p>
        {props.currentPrayerName} & {props.lastPrayerName}
      </p>
      <Button onClick={() => sendData(true)}>Yes</Button>
      <Button onClick={() => sendData(false)}>No</Button>
    </div>
  );
};

export default ProductiveStateView;
