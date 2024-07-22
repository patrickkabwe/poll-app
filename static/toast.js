const closeToastBtn = document.getElementById("closeToastBtn");
const toast = document.getElementById("toast");
const toastMessage = document.getElementById("toastMessage");
const toastIcon = document.getElementById("toastIcon");

function addIcon(ttype) {
  let icon = "";
  switch (ttype) {
    case "success":
      icon = `
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" class="h-5 w-5 text-current success">
                <path fill="currentColor" d="m10.6 13.8l-2.15-2.15q-.275-.275-.7-.275t-.7.275t-.275.7t.275.7L9.9 15.9q.3.3.7.3t.7-.3l5.65-5.65q.275-.275.275-.7t-.275-.7t-.7-.275t-.7.275zM12 22q-2.075 0-3.9-.788t-3.175-2.137T2.788 15.9T2 12t.788-3.9t2.137-3.175T8.1 2.788T12 2t3.9.788t3.175 2.137T21.213 8.1T22 12t-.788 3.9t-2.137 3.175t-3.175 2.138T12 22"></path>
            </svg>
        `;
      break;
    case "error":
      icon = `
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" class="h-5 w-5 text-current error">
                <path fill="currentColor" d="M12 17q.425 0 .713-.288T13 16t-.288-.712T12 15t-.712.288T11 16t.288.713T12 17m0-4q.425 0 .713-.288T13 12V8q0-.425-.288-.712T12 7t-.712.288T11 8v4q0 .425.288.713T12 13m0 9q-2.075 0-3.9-.788t-3.175-2.137T2.788 15.9T2 12t.788-3.9t2.137-3.175T8.1 2.788T12 2t3.9.788t3.175 2.137T21.213 8.1T22 12t-.788 3.9t-2.137 3.175t-3.175 2.138T12 22"/>
            </svg>
        `;
      break;
    default:
      icon = "";
  }
  toastIcon.innerHTML = icon;
}

function closeToast() {
  const toast = document.getElementById("toast");
  toast.classList.add("opacity-0");
  toast.classList.remove("opacity-100");
  toast.classList.add("translate-y-0");
  toast.classList.remove("translate-y-16");
}

document.body.addEventListener("showToast", (e) => {
  toast.classList.remove("opacity-0");
  toast.classList.add("opacity-100");
  toast.classList.remove("translate-y-0");
  toast.classList.add("translate-y-16");
  addIcon(e.detail.type);
  toastMessage.innerHTML = e.detail.message;
  const classes = [e.detail.type, "toast"];
  toastMessage.classList.add(...classes);
  toast.classList.add(...classes);
  clearTimeout(window.toastTimeout);
  window.toastTimeout = setTimeout(closeToast, 3000);
});

closeToastBtn.addEventListener("click", closeToast);
