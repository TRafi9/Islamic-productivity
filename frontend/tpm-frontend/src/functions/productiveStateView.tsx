import Button from "react-bootstrap/Button";

function sendData(value) {}

const productiveStateView = (props) => {
    
  // set productiveState back to false after person has clicked button

  //   props.setProductiveState(true);
  //   props.setDisplayType("");

  <div>
    <Button onClick={() => sendData(1)}>Yes</Button>
    <Button onClick={() => sendData(0)}>No</Button>
  </div>;
};

export default productiveStateView;
