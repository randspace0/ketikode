const burger = document.querySelector(".navbar-burger");
const navMenu = document.querySelector(".navbar-menu");

if (burger && navMenu) {
  burger.addEventListener("click", () => {
    const open = !navMenu.classList.contains("is-active");
    burger.classList.toggle("is-active", open);
    navMenu.classList.toggle("is-active", open);
    burger.setAttribute("aria-expanded", String(open));
  });
}

const dropdownUsername = document.querySelector("#dropdown-username");
const dropdownUsernameButton = document.querySelector("#dropdown-username-button");

if (dropdownUsername && dropdownUsernameButton) {
  const setOpen = (open) => {
    dropdownUsername.classList.toggle("is-active", open);
    dropdownUsernameButton.setAttribute("aria-expanded", String(open));
  };

  dropdownUsernameButton.addEventListener("click", () => {
    setOpen(!dropdownUsername.classList.contains("is-active"));
  });

  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && dropdownUsername.classList.contains("is-active")) {
      setOpen(false);
      dropdownUsernameButton.focus();
    }
  });

  document.addEventListener("click", (e) => {
    if (!dropdownUsername.contains(e.target)) setOpen(false);
  });
}

// Mark the link for the current page.
document.querySelectorAll(".navbar-start a.navbar-item").forEach((link) => {
  if (link.getAttribute("href") === window.location.pathname) {
    link.classList.add("is-active");
    link.setAttribute("aria-current", "page");
  }
});
