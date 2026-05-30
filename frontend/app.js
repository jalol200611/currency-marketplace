// REGISTER

async function registerUser() {
  const firstName = document.getElementById("registerFirstName").value;

  const lastName = document.getElementById("registerLastName").value;

  const email = document.getElementById("registerEmail").value;

  const phone = document.getElementById("registerPhone").value;

  const password = document.getElementById("registerPassword").value;

  const response = await fetch("http://127.0.0.1:8080/register", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      first_name: firstName,

      last_name: lastName,

      email: email,

      phone: phone,

      password: password,
    }),
  });

  // ERROR

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  const data = await response.json();

  console.log(data);

  // SAVE USER

  localStorage.setItem("user", JSON.stringify(data));

  // SUCCESS

  showNotification("🎉 Регистрация успешна");

  // REDIRECT

  setTimeout(() => {
    window.location.href = "dashboard.html";
  }, 1000);
}

// LOGIN

async function loginUser() {
  const email = document.getElementById("loginEmail").value;

  const password = document.getElementById("loginPassword").value;

  const response = await fetch("http://127.0.0.1:8080/login", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      email: email,
      password: password,
    }),
  });

  // ERROR

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  const data = await response.json();

  console.log(data);

  // SAVE USER

  localStorage.setItem("user", JSON.stringify(data));

  // SUCCESS

  showNotification("✅ Вход выполнен");

  // REDIRECT

  setTimeout(() => {
    window.location.href = "dashboard.html";
  }, 1000);
}

/* THEME */

function toggleTheme() {
  document.body.classList.toggle("light-mode");

  localStorage.setItem("theme", document.body.classList.contains("light-mode"));
}

/* LOAD THEME */

if (localStorage.getItem("theme") === "true") {
  document.body.classList.add("light-mode");
}

/* NOTIFICATIONS */

function showNotification(message) {
  const notification = document.getElementById("notification");

  notification.innerText = message;

  notification.classList.add("show");

  setTimeout(() => {
    notification.classList.remove("show");
  }, 3000);
}
