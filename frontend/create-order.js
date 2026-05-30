const user = JSON.parse(localStorage.getItem("user"));

/* USER INFO */

document.getElementById("userInfo").innerText = `ID: ${user.id}`;

/* LOAD WALLETS */

async function loadWallets() {
  const response = await fetch("http://127.0.0.1:8080/wallets");

  const wallets = await response.json();

  const balancesDiv = document.getElementById("balances");

  balancesDiv.innerHTML = "";

  wallets.forEach((wallet) => {
    if (wallet.user_id === user.id) {
      let currency = "";

      if (wallet.currency_id === 1) {
        currency = "🇺🇸 USD";
      }

      if (wallet.currency_id === 2) {
        currency = "🇪🇺 EUR";
      }

      if (wallet.currency_id === 3) {
        currency = "🇷🇺 RUB";
      }

      balancesDiv.innerHTML += `

                <p>
                    ${currency}: ${wallet.balance.toFixed(2)}
                </p>

            `;
    }
  });
}

/* CREATE ORDER */

async function createMarketOrder() {
  const sell_curr_id = document.getElementById("sellCurrency").value;

  const buy_curr_id = document.getElementById("buyCurrency").value;

  const amount = document.getElementById("marketAmount").value;

  const price = document.getElementById("marketPrice").value;

  const response = await fetch("http://127.0.0.1:8080/market/order", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      seller_id: user.id,

      sell_curr_id: Number(sell_curr_id),

      buy_curr_id: Number(buy_curr_id),

      amount: Number(amount),

      price: Number(price),
    }),
  });

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  const data = await response.text();

  showNotification(data);
}

/* LOGOUT */

function logout() {
  const confirmLogout = confirm("Вы хотите выйти?");

  if (confirmLogout) {
    localStorage.removeItem("user");

    window.location.href = "index.html";
  }
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
/* START */

loadWallets();
