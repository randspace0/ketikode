let deleteButton = document.querySelector("#delete");

// Only the owner gets a delete button.
if (deleteButton) deleteButton.addEventListener("click", () => {
  if (confirm("Are you sure you want to delete this?")) {
    fetch('/code/delete', {
      method: 'DELETE',
      headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
      },
      body: JSON.stringify({
          id: code.id,
      }),
    })
    .then((res) => res.json())
    .then((data) => {
        if (!data.success) {
          alert("Error: " + data.message);
        } else {
          window.location.replace("/myrecent");
        }
    });
  }
});