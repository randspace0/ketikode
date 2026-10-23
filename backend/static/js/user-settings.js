const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,16}$/;

let fieldState = {
    password: false,
    passwordOld: false,
};

let buttonSubmit = document.querySelector("button[type='submit']");

let fieldPassword = document.querySelector(".input#password");
let captionPassword = document.querySelector(".help#password");

let fieldPasswordOld = document.querySelector(".input#password-old");
let captionPasswordOld = document.querySelector(".help#password-old");

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

function fieldClear() {
    fieldPasswordOld.value = "";
    fieldPassword.value = "";
}

fieldPassword.addEventListener("change", () => {
    let isValidPassword = (passwordRegex).test(fieldPassword.value);
    let fieldColorForValid = () => (isValidPassword)?"is-success":"is-danger";
    let message = (isValidPassword)?"":"Invalid password. Password must have an uppercase, lowercase ,and numeric letter with 8-16 length.";

    fieldPassword.classList.remove(fieldColorForValid(false));
    fieldPassword.classList.add(fieldColorForValid(true));

    captionPassword.classList.remove(fieldColorForValid(false));
    captionPassword.classList.add(fieldColorForValid(true));

    captionPassword.innerText = message;

    fieldState.password = isValidPassword;
    buttonState();
});

fieldPasswordOld.addEventListener("change", () => {
    let isValidPasswordOld = (passwordRegex).test(fieldPasswordOld.value);
    let fieldColorForValid = () => (isValidPasswordOld)?"is-success":"is-danger";
    let message = (isValidPasswordOld)?"":"Invalid password. Password must have an uppercase, lowercase ,and numeric letter with 8-16 length.";

    fieldPasswordOld.classList.remove(fieldColorForValid(false));
    fieldPasswordOld.classList.add(fieldColorForValid(true));

    captionPasswordOld.classList.remove(fieldColorForValid(false));
    captionPasswordOld.classList.add(fieldColorForValid(true));

    captionPasswordOld.innerText = message;

    fieldState.passwordOld = isValidPasswordOld;
    buttonState();
});

buttonSubmit.addEventListener("click", () => {
    notification.classList.add("is-hidden");
    notification.classList.remove("is-danger");
    notification.classList.remove("is-success");
    fetch('/settings', {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            password: fieldPassword.value,
            passwordOld: fieldPasswordOld.value,
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
          notification.innerText = "Data has been edited";
          fieldClear();
        }
    });
});