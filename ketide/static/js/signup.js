const usernameRegex = /^[a-z0-9]{5,16}$/;
const emailRegex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,16}$/;

let fieldState = {
    email: false,
    username: false,
    password: false,
    passwordConfirmation: false,
};

let buttonSubmit = document.querySelector("button[type='submit']");

let fieldUsername = document.querySelector(".input#username");
let captionUsername = document.querySelector(".help#username");

let fieldEmail = document.querySelector(".input#email");
let captionEmail = document.querySelector(".help#email");

let fieldPassword = document.querySelector(".input#password");
let captionPassword = document.querySelector(".help#password");

let fieldPasswordConfirmation = document.querySelector(".input#password-confirmation");
let captionPasswordConfirmation = document.querySelector(".help#password-confirmation");

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
    let isValidUsername = (usernameRegex).test(fieldUsername.value);
    let fieldColorForValid = (valid) => (valid === isValidUsername)?"is-success":"is-danger";
    let message = (isValidUsername)?"Username is available":"Invalid username";

    fieldUsername.classList.remove(fieldColorForValid(false));
    fieldUsername.classList.add(fieldColorForValid(true));

    captionUsername.classList.remove(fieldColorForValid(false));
    captionUsername.classList.add(fieldColorForValid(true));

    captionUsername.innerText = message;

    fieldState.username = isValidUsername;
    buttonState();
});

fieldEmail.addEventListener("change", () => {
    let isValidEmail = (emailRegex).test(fieldEmail.value);
    let fieldColorForValid = (valid) => (valid === isValidEmail)?"is-success":"is-danger";
    let message = (isValidEmail)?"Username is available":"Invalid username";

    fieldEmail.classList.remove(fieldColorForValid(false));
    fieldEmail.classList.add(fieldColorForValid(true));

    captionEmail.classList.remove(fieldColorForValid(false));
    captionEmail.classList.add(fieldColorForValid(true));

    captionEmail.innerText = message;

    fieldState.email = isValidEmail;
    buttonState();
});

fieldPassword.addEventListener("change", () => {
    let isValidPassword = (passwordRegex).test(fieldPassword.value);
    let fieldColorForValid = (valid) => (valid === isValidPassword)?"is-success":"is-danger";
    let message = (isValidPassword)?"":"Invalid password. Password must have an uppercase, lowercase ,and numeric letter with 8-16 length.";

    fieldPassword.classList.remove(fieldColorForValid(false));
    fieldPassword.classList.add(fieldColorForValid(true));

    captionPassword.classList.remove(fieldColorForValid(false));
    captionPassword.classList.add(fieldColorForValid(true));

    captionPassword.innerText = message;

    fieldState.password = isValidPassword;
    buttonState();
});

fieldPasswordConfirmation.addEventListener("change", () => {
    let isValidPasswordConfirmation = fieldPassword.value == fieldPasswordConfirmation.value;
    let fieldColorForValid = (valid) => (valid === isValidPasswordConfirmation)?"is-success":"is-danger";
    let message = (isValidPasswordConfirmation)?"":"Password does not match.";

    fieldPasswordConfirmation.classList.remove(fieldColorForValid(false));
    fieldPasswordConfirmation.classList.add(fieldColorForValid(true));

    captionPasswordConfirmation.classList.remove(fieldColorForValid(false));
    captionPasswordConfirmation.classList.add(fieldColorForValid(true));

    captionPasswordConfirmation.innerText = message;

    fieldState.passwordConfirmation = isValidPasswordConfirmation;
    buttonState();
});

buttonSubmit.addEventListener("click", () => {
    notification.classList.add("is-hidden");
    notification.classList.remove("is-danger");
    notification.classList.remove("is-success");
    fetch('/signup', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            email: fieldEmail.value,
            username: fieldUsername.value,
            password: fieldPassword.value,
            passwordConfirmation: fieldPasswordConfirmation.value,
        }),
    })
    .then((res) => res.json())
    .then((data) => {
        notification.classList.remove("is-hidden");
        if (!data.success) {
            notification.classList.add("is-danger");
            notification.innerText = data.message;
        } else {
            notification.classList.add("is-success");
            form.classList.add("is-hidden");
            notification.innerText = "Congratulations! You account has been created!";
        }
    });
});