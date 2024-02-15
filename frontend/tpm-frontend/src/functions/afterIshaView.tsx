import { useEffect } from "react";

const AfterIshaView = (props: any) => {
  useEffect(() => {
    props.setProductiveState(false);
  });

  return (
    <div className="container mt-5">
      <div className="row">
        <div className="col">
          <p className="lead">
            Thanks for submitting, please come back tomorrow!
          </p>
        </div>
      </div>
    </div>
  );
};

export default AfterIshaView;
