import Head from "next/head";
import styles from "@/styles/Home.module.css";
import { Inter } from "next/font/google";
import { useState } from "react";
import React, { ChangeEvent } from "react";

const inter = Inter({ subsets: ["latin"] });

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
    sanitiseEmail(event.target.value);
  };

  const handleCreateUserPassword = (event: ChangeEvent<HTMLInputElement>) => {
    // sets password and sanitises input
    setUserPassword(event.target.value);
    sanitisePassword(event.target.value);
  };

  //TODO need to verify username and password sanitization here

  // functions for sanitising input

  function sanitiseEmail(email: string) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (emailRegex.test(email)) {
      setEmailSanitiseCheck(true);
    } else {
      setEmailSanitiseCheck(false);
    }
  }

  function sanitisePassword(password: string) {
    // Ensure the password contains at least one special character
    // Ensure the password contains at least one capital letter
    const specialCharacterRegex = /[!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]/;
    const capitalLetterRegex = /[A-Z]/;

    if (
      password.length < 10 ||
      !specialCharacterRegex.test(password) ||
      !capitalLetterRegex.test(password)
    ) {
      setPasswordSanitiseCheck(false);
    } else {
      setPasswordSanitiseCheck(true);
    }
  }

  function showPasswordWarning() {
    if (passwordSanitiseCheck) {
      return <></>;
    } else {
      return (
        <p>
          Password needs to be 10+ letters, have a special character and a
          number
        </p>
      );
    }
  }
  function showEmailWarning() {
    if (emailSanitiseCheck) {
      return <></>;
    } else {
      return <p>Please enter a valid email address</p>;
    }
  }

  // here we create the struct of submissionData and pass it through to the backend
  type SubmissionData = {
    userEmail: string;
    userPassword: string;
  };
  const submitNewUser = async (data: SubmissionData) => {
    try {
      const response = await fetch("http://localhost:8080/api/v1/createUser", {
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
      console.error("Error:");
    }
  };

  async function submit() {
    const SubmissionData: SubmissionData = {
      userEmail: userEmail,
      userPassword: UserPassword,
    };
    const statusResponse = await submitNewUser(SubmissionData);
    console.log(statusResponse);
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
          <form className="register-form">
            <div className="form-group">
              {showEmailWarning()}
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
              {showPasswordWarning()}
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
              <button
                type="submit"
                className="btn btn-primary"
                onClick={(e) => {
                  e.preventDefault(); // Prevents the default form submission
                  submit();
                }}
              >
                Submit
              </button>
            </div>
          </form>
        </main>
      </div>
    </>
  );
}
