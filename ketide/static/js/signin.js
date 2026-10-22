let fieldState = {
  username: false,
  password: false,
};

let buttonSubmit = document.querySelector("button[type='submit']");

let fieldUsername = document.querySelector(".input#username");

let fieldPassword = document.querySelector(".input#password");

let notification = document.querySelector("#notification");
let form = document.querySelector("#form");

buttonSubmit.disabled = true;

function buttonState() {
  let objArray = Object.values(fieldState);
  for (let i in objArray) {
      if (!objArray[i]) {
          buttonSubmit.disabled = true;
          return;
      }
  }
  buttonSubmit.disabled = false;
}

fieldUsername.addEventListener("change", () => {
  let isValidUsername = fieldUsername.value != "";
  let fieldColorForValid = (valid) => (valid === isValidUsername)?"is-success":"is-danger";

  fieldUsername.classList.remove(fieldColorForValid(false));
  fieldUsername.classList.add(fieldColorForValid(true));

  fieldState.username = isValidUsername;
  buttonState();
});

fieldPassword.addEventListener("change", () => {
  let isValidPassword = fieldPassword.value != "";
  let fieldColorForValid = (valid) => (valid === isValidPassword)?"is-success":"is-danger";

  fieldPassword.classList.remove(fieldColorForValid(false));
  fieldPassword.classList.add(fieldColorForValid(true));

  fieldState.password = isValidPassword;
  buttonState();
});

buttonSubmit.addEventListener("click", () => {
  notification.classList.add("is-hidden");
  notification.classList.remove("is-danger");
  notification.classList.remove("is-success");
  fetch('/signin', {
      method: 'POST',
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      },
      body: JSON.stringify({
          username: fieldUsername.value,
          password: fieldPassword.value,
      }),
  })
  .then((res) => res.json())
  .then((data) => {
      notification.classList.remove("is-hidden");
      if (!data.success) {
          notification.classList.add("is-danger");
          notification.innerText = data.message;
      } else {
          form.classList.add("is-hidden");
          notification.classList.add("is-success");
          notification.innerText = "Welcome " + fieldUsername.value + " !";
          window.location.replace("/myrecent");
      }
  });
});