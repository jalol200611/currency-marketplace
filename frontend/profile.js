const user = JSON.parse(localStorage.getItem("user"));

/* USER INFO */

document.getElementById("userInfo").innerText = `ID: ${user.id}`;

document.getElementById("email").innerText = user.email;

/* PROFILE */

document.getElementById("profileId").innerText = `ID: ${user.id}`;

document.getElementById("profileName").innerText =
  `${user.first_name} ${user.last_name}`;

document.getElementById("profileEmail").innerText = `Email: ${user.email}`;

/* LOAD WALLETS */

async function loadWallets() {
  const response = await fetch("http://127.0.0.1:8080/wallets");

  const wallets = await response.json();

  const walletsDiv = document.getElementById("wallets");

  const balancesDiv = document.getElementById("balances");

  walletsDiv.innerHTML = "";

  balancesDiv.innerHTML = "";

  let totalUSD = 0;

  wallets.forEach((wallet) => {
    if (wallet.user_id === user.id) {
      let currency = "";

      // USD

      if (wallet.currency_id === 1) {
        currency = "🇺🇸 USD";

        totalUSD += wallet.balance;
      }

      // EUR

      if (wallet.currency_id === 2) {
        currency = "🇪🇺 EUR";

        totalUSD += wallet.balance * 1.17;
      }

      // RUB

      if (wallet.currency_id === 3) {
        currency = "🇷🇺 RUB";

        totalUSD += wallet.balance * 0.013;
      }

      walletsDiv.innerHTML += `
            
                <p>
                    ${currency}: ${wallet.balance.toFixed(2)}
                </p>

            `;

      balancesDiv.innerHTML += `
            
                <p>
                    ${currency}: ${wallet.balance.toFixed(2)}
                </p>

            `;
    }
  });

  document.getElementById("totalBalance").innerText = `$${totalUSD.toFixed(2)}`;
}

/* LOAD TRANSACTIONS */

async function loadTransactions() {
  const response = await fetch("http://127.0.0.1:8080/transactions");

  const transactions = await response.json();

  const transactionsDiv = document.getElementById("transactions");

  transactionsDiv.innerHTML = "";

  transactions.forEach((transaction) => {
    if (transaction.buyer_id === user.id || transaction.seller_id === user.id) {
      let sellCurrency = "";

      let buyCurrency = "";

      // SELL CURRENCY

      if (transaction.sell_curr_id === 1) {
        sellCurrency = "🇺🇸 USD";
      }

      if (transaction.sell_curr_id === 2) {
        sellCurrency = "🇪🇺 EUR";
      }

      if (transaction.sell_curr_id === 3) {
        sellCurrency = "🇷🇺 RUB";
      }

      // BUY CURRENCY

      if (transaction.buy_curr_id === 1) {
        buyCurrency = "🇺🇸 USD";
      }

      if (transaction.buy_curr_id === 2) {
        buyCurrency = "🇪🇺 EUR";
      }

      if (transaction.buy_curr_id === 3) {
        buyCurrency = "🇷🇺 RUB";
      }

      // TYPE

      let type = "";

      if (transaction.buyer_id === user.id) {
        type = "🟢 Покупка";
      }

      if (transaction.seller_id === user.id) {
        type = "🔴 Продажа";
      }
      // SELLER / BUYER
      let sellerText =
        transaction.seller_id === user.id ? "Вы" : transaction.seller_name;

      let buyerText =
        transaction.buyer_id === user.id ? "Вы" : transaction.buyer_name;

      transactionsDiv.innerHTML += `

<div class="transaction-card">

    <div class="transaction-header">

        <div class="transaction-type">

            ${type}

        </div>

        <div class="transaction-time">

            ${new Date(transaction.created_at).toLocaleString()}

        </div>

    </div>

    <div class="transaction-body">

        <div>

            <strong>
                Валютная пара
            </strong>

            <p>
                ${sellCurrency}
                →
                ${buyCurrency}
            </p>

        </div>

        <div>

            <strong>
                Количество
            </strong>

            <p>
                ${transaction.amount}
            </p>

        </div>

        <div>

            <strong>
                Цена
            </strong>

            <p>
                ${transaction.price}
            </p>

        </div>

        <div>

            <strong>
                Продавец
            </strong>

            <p>
             ${sellerText}
            </p>

        </div>

        <div>

            <strong>
                Покупатель
            </strong>

            <p>
              ${buyerText}
            </p>

        </div>

    </div>

</div>

`;
    }
  });
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

loadTransactions();
