const checkbox = document.getElementById("theme-toggle");
const root = document.documentElement;
const stored = localStorage.getItem("site-theme");

if (stored === "dark") {
  root.setAttribute("data-theme", "dark");
  checkbox.checked = true;
} else {
  root.setAttribute("data-theme", "light");
  checkbox.checked = false;
}

checkbox.addEventListener("change", () => {
  const isDark = checkbox.checked;
  root.setAttribute("data-theme", isDark ? "dark" : "light");
  localStorage.setItem("site-theme", isDark ? "dark" : "light");
});

const label = document.querySelector(".switch-ui");
label.addEventListener("keydown", (e) => {
  if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    checkbox.checked = !checkbox.checked;
    checkbox.dispatchEvent(new Event("change"));
  }
});
