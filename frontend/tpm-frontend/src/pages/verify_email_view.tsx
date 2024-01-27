import Head from "next/head";
import styles from "@/styles/Home.module.css";
import { Inter } from "next/font/google";
import { SetStateAction, useState, Dispatch } from "react";
import React, { ChangeEvent } from "react";
import Router from "next/router";
import {
  showEmailWarning,
  showPasswordWarning,
  sanitiseEmail,
  sanitisePassword,
} from "@/functions/loginFunctions";

const inter = Inter({ subsets: ["latin"] });

//TODO on this page, the user is registered, but a response is sent from server only when the function on the backend completes
// because of this, there is a delay in showRegistrationMessage() to resolve
// add a spinner/ loading animation to show the user that something is being processed while waiting for showRegistrationMessage to resolve

export default function RegisterUser() {
  const [userEmail, setCreateUserEmail] = useState<string>("");
  const [emailSanitiseCheck, setEmailSanitiseCheck] = useState<boolean>(false);

  const [UserPassword, setUserPassword] = useState<string>("");

  const [passwordSanitiseCheck, setPasswordSanitiseCheck] =
    useState<boolean>(false);
  const handleCreateUserEmailChange = (
    event: ChangeEvent<HTMLInputElement>
  ) => {
    setCreateUserEmail(event.target.value);
    sanitiseEmail(event.target.value, setEmailSanitiseCheck);
  };

  const handleCreateUserPassword = (event: ChangeEvent<HTMLInputElement>) => {
    // sets password and sanitises input
    setUserPassword(event.target.value);
    sanitisePassword(event.target.value, setPasswordSanitiseCheck);
  };

  const [emailVerificationCode, setEmailVerificationCode] = useState<
    number | null
  >(null);
  // the function below updates the value of emailVerificationCode when the user types it in
  const handleEmailVerificationCode = (
    event: ChangeEvent<HTMLInputElement>
  ) => {
    const numericValue = parseInt(event.target.value, 10);

    // Check if the conversion is successful and it's a 6-digit integer
    if (!isNaN(numericValue) && String(numericValue).length === 6) {
      setEmailVerificationCode(numericValue);
    }
  };

  const [verifyEmailView, setVerifyEmailView] = useState<boolean>(false);

  //TODO need to verify username and password sanitization here

  // functions for sanitising input

  function VerificationMessage() {
    if (verifiedUserEmailResponse) {
      switch (verifiedUserEmailResponse) {
        case 200:
          // setVerifyEmailView(true);
          Router.push("login");
          return null;
        case 400:
          return (
            <p>Error creating user, please contact developer: status: 400.</p>
          );
        default:
          return <p>Error: Unknown status code received.</p>;
      }
    }
    return <></>;
  }

  // verification email sending form data to backend and retrieving response variables

  const [verifiedUserEmailResponse, setVerifiedUserEmailResponse] = useState<
    number | null
  >();
  type VerifyEmailData = {
    userEmail: string;
    verificationCode: number | null;
  };

  const submitVerificationCheck = async (data: VerifyEmailData) => {
    const response = await fetch(
      "http://localhost:8080/api/v1/userVerification",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      }
    );

    return response.status;
  };

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (emailSanitiseCheck) {
      const verifySubmissionData: VerifyEmailData = {
        userEmail: userEmail,
        verificationCode: emailVerificationCode,
      };

      const response = await submitVerificationCheck(verifySubmissionData);

      // response needs to be response.status return on the backend when created, then create a check on if verifiedUserEmailResponse is 200, if so redirect to login page, else throw errors
      setVerifiedUserEmailResponse(response);
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
        <main className={`${styles.main} ${inter.className}`}>
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
              <label>Verification Code</label>
              <input
                type="text"
                className="form-control"
                id="verification code"
                placeholder="Please enter your 6 digit verification code here"
                onChange={handleEmailVerificationCode}
              />
            </div>
            <div className="form-group">
              {VerificationMessage()}
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
