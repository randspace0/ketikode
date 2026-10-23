const usernameRegex = /^[a-z0-9]{5,16}$/;
const emailRegex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

let fieldState = {};

// TODO
// iconUsername is not working when I tried to change the icon inside the event listener

let buttonSubmit = document.querySelector("button[type='submit']");

// let iconUsername = document.querySelector(".fas#username");
let fieldUsername = document.querySelector(".input#username");
let captionUsername = document.querySelector(".help#username");

let fieldEmail = document.querySelector(".input#email");
let captionEmail = document.querySelector(".help#email");

buttonSubmit.disabled = true;

function buttonState() {
    for (let key in Object.keys(fieldState)) {
        if (!fieldState[key]) {
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

    fieldState["username"] = isValidUsername;
    buttonState();
    
    // iconUsername.classList.remove("fa-exclamation-triangle");
    // iconUsername.classList.add("fa-check");
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

    fieldState["email"] = isValidEmail;
    buttonState();
});

// TODO
// Implement AJAX signup

buttonSubmit.addEventListener("click", () => {

});