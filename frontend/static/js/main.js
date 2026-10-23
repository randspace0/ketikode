document.querySelectorAll(".navbar-burger").forEach((element) => {
    element.addEventListener("click", () => {
        let className = "is-active";
        let activate = !element.classList.contains(className);
        
        if (activate) {
            element.classList.add(className);
        } else {
            element.classList.remove(className);
        }

        document.querySelectorAll(".navbar-menu").forEach((element) => {
            if (activate) {
                element.classList.add(className);
            } else {
                element.classList.remove(className);
            }
        });
    });
});