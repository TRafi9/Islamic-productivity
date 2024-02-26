import Head from "next/head";
import styles from "@/styles/Home.module.css";
import { Roboto_Mono } from "next/font/google";
import { useState } from "react";

import Router from "next/router";
import React, { ChangeEvent } from "react";
import {
  showEmailWarning,
  showPasswordWarning,
  sanitiseEmail,
  sanitisePassword,
} from "@/functions/loginFunctions";
import LoadingSpinner from "@/components/LoadingSpinner";

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

  const [loading, setLoading] = useState<boolean>(false);

  //TODO need to verify username and password sanitization here

  // functions for sanitising input
  const [responseErr, setResponseErr] = useState<string>("");
  function showRegistrationMessage() {
    if (loading) {
      return <LoadingSpinner />;
    }
    if (submitResponseStatus) {
      switch (submitResponseStatus) {
        case 200:
          Router.push("verify_email_view");
          return null;
        // case 208:
        //   return <p> Error creating user, email already in use</p>;
        // case 400:
        //   return (
        //     <p>Error creating user, please contact developer: status: 400.</p>
        //   );
        // case 500:
        //   return (
        //     <p>Error creating user, please contact developer: status: 500.</p>
        //   );
        default:
          return <p>{responseErr}</p>;
      }
    }
    return <></>;
  }

  // here we create the struct of submissionData and pass it through to the backend
  type SubmissionData = {
    userEmail: string;
    userPassword: string;
  };
  const submitNewUser = async (data: SubmissionData) => {
    setLoading(true);
    const response = await fetch("http://localhost:8080/api/v1/createUser", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    const responseData = await response.json();
    setResponseErr(responseData["error"]);

    return response.status;
  };
  // once user submites, this updates with response code, which triggers showRegistrationMessage() to display something
  const [submitResponseStatus, setSubmitResponseStatus] = useState<
    number | null
  >(null);

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    //TODO, rewrite maybe? no need for submitNewUser to be its own function, can directly call it in here, like in submitVerification function
    e.preventDefault();
    if (emailSanitiseCheck && passwordSanitiseCheck) {
      const SubmissionData: SubmissionData = {
        userEmail: userEmail,
        userPassword: UserPassword,
      };
      const response = await submitNewUser(SubmissionData);

      setSubmitResponseStatus(response);
      setLoading(false);
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
            <h1> User Registration</h1>
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
              {showPasswordWarning(passwordSanitiseCheck)}
              <label htmlFor="exampleInputPassword1">Password</label>
              <input
                type="password"
                className="form-control"
                id="password"
                placeholder="Password"
                onChange={handleCreateUserPassword}
              />
            </div>
            <div className="form-group">
              {showRegistrationMessage()}
              {loading ? null : (
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={!emailSanitiseCheck || !passwordSanitiseCheck}
                >
                  Register
                </button>
              )}
            </div>
          </form>
        </main>
      </div>
    </>
  );
}
