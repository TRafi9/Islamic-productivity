import Button from "react-bootstrap/Button";
// need to send data to handleSendData from getCurrentPrayer + next prayer
const productiveStateView = (props: any) => {
  // set productiveState back to false after person has clicked button

  function sendData(value: boolean) {
    props.setProductiveState(false);
  }

  const handleSendData = async () => {
    try {
      const response = await fetch("/api/postProductivityValue", {
        method: "POST",
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
      console.error("Error:", error.message);
    }
  };

  <div>
    <Button onClick={() => sendData(true)}>Yes</Button>
    <Button onClick={() => sendData(false)}>No</Button>
  </div>;
};

export default productiveStateView;
