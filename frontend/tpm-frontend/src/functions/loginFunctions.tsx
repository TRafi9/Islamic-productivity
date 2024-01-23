// warning functions expect a usestate boolean value to be passed to display warnings

export function showEmailWarning(emailSanitiseCheck: boolean) {
  if (emailSanitiseCheck) {
    return <></>;
  } else {
    return <p>Please enter a valid email address</p>;
  }
}

export function showPasswordWarning(passwordSanitiseCheck: boolean) {
  if (passwordSanitiseCheck) {
    return <></>;
  } else {
    return (
      <p>
        Password needs to be 10+ letters, have a special character and a number
      </p>
    );
  }
}

// email and password santisation functions both expect a string and setstate value to be passed to them
export function sanitiseEmail(
  email: string,
  setEmailSanitiseCheck: React.Dispatch<React.SetStateAction<boolean>>
) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (emailRegex.test(email)) {
    setEmailSanitiseCheck(true);
  } else {
    setEmailSanitiseCheck(false);
  }
}

export function sanitisePassword(
  password: string,
  setPasswordSanitiseCheck: React.Dispatch<React.SetStateAction<boolean>>
) {
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
