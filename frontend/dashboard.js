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

/* TOP UP */

async function topUp() {
  const user_id = user.id;

  // CARD

  const cardNumber = document.getElementById("cardNumber").value;

  const cardName = document.getElementById("cardName").value;

  const cardDate = document.getElementById("cardDate").value;

  const cardCvv = document.getElementById("cardCvv").value;

  // VALIDATION

  if (!cardNumber || !cardName || !cardDate || !cardCvv) {
    showNotification("Заполните данные карты");

    return;
  }

  // BALANCE

  const currency_id = document.getElementById("topupCurrencyId").value;

  const amount = document.getElementById("topupAmount").value;

  const response = await fetch("http://127.0.0.1:8080/topUp", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      user_id: Number(user_id),

      currency_id: Number(currency_id),

      amount: Number(amount),
    }),
  });

  const data = await response.text();

  showNotification(data);

  // CLEAR INPUTS

  document.getElementById("cardNumber").value = "";

  document.getElementById("cardName").value = "";

  document.getElementById("cardDate").value = "";

  document.getElementById("cardCvv").value = "";

  document.getElementById("topupAmount").value = "";

  loadWallets();
}

/* CONVERT */

async function convertCurrency() {
  const user_id = user.id;

  const from = document.getElementById("fromCurrency").value;

  const to = document.getElementById("toCurrency").value;

  const amount = document.getElementById("convertAmount").value;

  const response = await fetch("http://127.0.0.1:8080/convert", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      user_id: Number(user_id),

      from: Number(from),

      to: Number(to),

      amount: Number(amount),
    }),
  });

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  const data = await response.json();

  showNotification(`💱 Конвертация успешна: ${data.result}`);
  document.getElementById("convertAmount").value = "";

  loadWallets();
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
/* CARD PREVIEW */

const cardNumberInput = document.getElementById("cardNumber");

const cardNameInput = document.getElementById("cardName");

const cardDateInput = document.getElementById("cardDate");

/* CARD NUMBER */

cardNumberInput.addEventListener("input", () => {
  let value = cardNumberInput.value;

  value = value.replace(/\s/g, "");

  value = value.replace(/(.{4})/g, "$1 ");

  document.getElementById("previewCardNumber").innerText =
    value || "0000 0000 0000 0000";
});

/* CARD NAME */

cardNameInput.addEventListener("input", () => {
  document.getElementById("previewCardName").innerText =
    cardNameInput.value.toUpperCase() || "YOUR NAME";
});

async function loadRates() {
  const response = await fetch("http://127.0.0.1:8080/rates");

  const rates = await response.json();

  const ratesDiv = document.getElementById("rates");

  ratesDiv.innerHTML = "";

  rates.forEach((rate) => {
    let pair = "";

    if (rate.from_id === 1 && rate.to_id === 2) {
      pair = "🇺🇸 USD / 🇪🇺 EUR";
    }

    if (rate.from_id === 1 && rate.to_id === 3) {
      pair = "🇺🇸 USD / 🇷🇺 RUB";
    }

    if (pair === "") {
      return;
    }

    ratesDiv.innerHTML += `
  <div class="rate-card">

    <div class="rate-title">
      ${pair}
    </div>

    <div class="rate-price">
      ${rate.coefficient.toFixed(4)}
    </div>

    <div class="rate-change positive">
      ↑ LIVE
    </div>

    <div class="rate-chart">
      <span></span>
    </div>

  </div>
`;
  });
}

/* CARD DATE */

cardDateInput.addEventListener("input", () => {
  document.getElementById("previewCardDate").innerText =
    cardDateInput.value || "00/00";
});
/* START */
loadRates();

setInterval(loadRates, 5000);
loadWallets();
