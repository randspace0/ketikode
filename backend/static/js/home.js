const stdinButton = document.querySelector("#button-stdin");
const stdinField = document.querySelector("#field-stdin");

stdinButton.addEventListener("click", () => {
  const hidden = stdinField.classList.toggle("is-hidden");
  stdinButton.setAttribute("aria-expanded", String(!hidden));
  if (!hidden) stdinField.querySelector("textarea").focus();
});

// Show progress while the code runs; the page reloads when it finishes.
const runButton = document.querySelector("#button-run");
document.querySelector("#form").addEventListener("submit", () => {
  runButton.disabled = true;
  runButton.classList.add("is-loading");
});

// Coming back with the browser's back button keeps the old DOM state.
window.addEventListener("pageshow", () => {
  runButton.disabled = false;
  runButton.classList.remove("is-loading");
});
