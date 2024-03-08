import Head from "next/head";
import styles from "@/styles/Home.module.css";
import { Roboto_Mono } from "next/font/google";
import { useState } from "react";
import React, { ChangeEvent } from "react";
import Router from "next/router";

import { showEmailWarning, sanitiseEmail } from "@/functions/loginFunctions";

const roboto = Roboto_Mono({
  weight: "400",
  subsets: ["latin"],
});

//TODO on this page, the user is registered, but a response is sent from server only when the function on the backend completes
// because of this, there is a delay in showRegistrationMessage() to resolve
// add a spinner/ loading animation to show the user that something is being processed while waiting for showRegistrationMessage to resolve

export default function RegisterUser() {
  const [userEmail, setCreateUserEmail] = useState<string>("");
  const [emailSanitiseCheck, setEmailSanitiseCheck] = useState<boolean>(false);

  const handleCreateUserEmailChange = (
    event: ChangeEvent<HTMLInputElement>
  ) => {
    setCreateUserEmail(event.target.value);
    sanitiseEmail(event.target.value, setEmailSanitiseCheck);
  };
  const [responseErr, setResponseErr] = useState<string>("");

  function ResendVerificationMessage() {
    if (verifiedUserEmailResponse) {
      switch (verifiedUserEmailResponse) {
        case 200:
          // setVerifyEmailView(true);
          Router.push("verify_email_view");
          return null;

        default:
          return <p>{responseErr}.</p>;
      }
    }
    return <></>;
  }

  // verification email sending form data to backend and retrieving response variables

  const [verifiedUserEmailResponse, setVerifiedUserEmailResponse] = useState<
    number | null
  >(null);

  type ReverifyEmailData = {
    userEmail: string;
  };

  const submitVerificationCheck = async (data: ReverifyEmailData) => {
    const response = await fetch("/api/postResetEmailVerification", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    // can use responseData to pull error code out to frontend
    const responseData = await response.json();
    setResponseErr(responseData["error"]);
    return response.status;
  };

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (emailSanitiseCheck) {
      const verifySubmissionData: ReverifyEmailData = {
        userEmail: userEmail,
      };

      await submitVerificationCheck(verifySubmissionData).then((response) => {
        setVerifiedUserEmailResponse(response);
      });

      // response needs to be response.status return on the backend when created, then create a check on if verifiedUserEmailResponse is 200, if so redirect to login page, else throw errors
      // setVerifiedUserEmailResponse(response);
    }
  }

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
      <div>
        <main className={`${styles.main} ${roboto.className}`}>
          <form className="register-form" onSubmit={(e) => submit(e)}>
            <h1>Verify Email</h1>
            <div className="form-group">
              {showEmailWarning(emailSanitiseCheck)}
              <label>Email address</label>
              <br />
              <input
                type="email"
                className="form-control"
                id="email"
                placeholder="name@example.com"
                onChange={handleCreateUserEmailChange}
              />
            </div>
            <div className="form-group">
              {ResendVerificationMessage()}
              <button
                type="submit"
                className="btn btn-primary"
                disabled={!emailSanitiseCheck}
              >
                Verify
              </button>
            </div>
          </form>
        </main>
      </div>
    </>
  );
}
