import Head from "next/head";
import styles from "@/styles/Home.module.css";
import { Inter } from "next/font/google";
import { Roboto_Mono } from "next/font/google";
import { useState } from "react";
import Link from "next/link";
import React, { ChangeEvent } from "react";
import Router from "next/router";
import {
  sanitiseEmail,
  sanitisePassword,
  showEmailWarning,
  showPasswordWarning,
} from "@/functions/loginFunctions";
import LoadingSpinner from "@/components/LoadingSpinner";
// import "bootstrap/dist/css/bootstrap.min.css";

const roboto = Roboto_Mono({
  weight: "400",
  subsets: ["latin"],
});

export default function LoginUser() {
  const [userEmail, setUserEmail] = useState<string>("");
  const [emailSanitiseCheck, setEmailSanitiseCheck] = useState<boolean>(false);

  const [userPassword, setUserPassword] = useState<string>("");
  const [passwordSanitiseCheck, setPasswordSanitiseCheck] =
    useState<boolean>(false);

  const [submitResponseStatus, setSubmitResponseStatus] = useState<
    number | null
  >(null);
  const [loading, setLoading] = useState<boolean>(false);

  const handleUserEmailChange = (event: ChangeEvent<HTMLInputElement>) => {
    setUserEmail(event.target.value);
    sanitiseEmail(event.target.value, setEmailSanitiseCheck);
  };

  const handleUserPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    setUserPassword(event.target.value);
    sanitisePassword(event.target.value, setPasswordSanitiseCheck);
  };

  const [responseErr, setResponseErr] = useState<string>("");
  function showLoginMessage() {
    if (loading) {
      return <LoadingSpinner />;
    }

    if (submitResponseStatus !== null) {
      switch (submitResponseStatus) {
        case 200:
          Router.push("main_page");
          return <p>Login successful!</p>;
        case 406:
          return (
            <>
              <Link href={"verify_email_view"}>
                {" "}
                Please verify email before logging in
              </Link>
              <br></br>
            </>
          );
        // case 500:
        //   return <p>Error logging in, please contact the developer</p>;
        default:
          return <p>{responseErr}</p>;
      }
    }

    return null;
  }

  const loginUser = async () => {
    setLoading(true);

    const response = await fetch("/api/postLoginUser", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        userEmail,
        userPassword,
      }),
      credentials: "include",
    });

    const responseData = await response.json();
    setResponseErr(responseData["error"]);

    const authorizationHeader = response.headers.get("Authorization");

    if (authorizationHeader !== null && authorizationHeader !== "") {
      sessionStorage.setItem("jwt", authorizationHeader);
      console.log("stored jwt in session");
    }

    setSubmitResponseStatus(response.status);
    setLoading(false);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (emailSanitiseCheck && passwordSanitiseCheck) {
      await loginUser();
    }
  };

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
        <main className={`${styles.main} ${roboto.className} black-text`}>
          <div className="card">
            <div className="card-body">
              <h1 className="card-title">Login</h1>
              <form className="login-form" onSubmit={(e) => handleSubmit(e)}>
                <div className="form-group">
                  {showEmailWarning(emailSanitiseCheck)}
                  <label>Email address</label>
                  <input
                    type="email"
                    className="form-control"
                    id="email"
                    placeholder="name@gmail.com"
                    onChange={handleUserEmailChange}
                  />
                </div>
                <div className="form-group">
                  {showPasswordWarning(passwordSanitiseCheck)}
                  <label>Password</label>
                  <input
                    type="password"
                    className="form-control"
                    id="password"
                    placeholder="Password"
                    onChange={handleUserPasswordChange}
                  />
                </div>
                <div className="form-group center">
                  <div>{showLoginMessage()}</div>
                  {loading ? null : ( // Render login button only if loading is false
                    <button
                      type="submit"
                      className="btn btn-primary"
                      disabled={!emailSanitiseCheck || !passwordSanitiseCheck}
                    >
                      Login
                    </button>
                  )}
                  <Link href={"register_user"}> Click here to register!</Link>
                </div>
              </form>
            </div>
          </div>
        </main>
      </div>
    </>
  );
}
