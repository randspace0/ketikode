document.querySelector("#button-stdin").addEventListener("click", (e) => {
  e.preventDefault();
  document.querySelector("#field-stdin").classList.toggle("is-hidden");
});