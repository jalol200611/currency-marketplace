const user = JSON.parse(localStorage.getItem("user"));

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

/* LOAD ORDERS */

async function loadMarketOrders() {
  const response = await fetch("http://127.0.0.1:8080/markets");

  const data = await response.json();

  const ordersDiv = document.getElementById("orders");

  ordersDiv.innerHTML = "";

  data.forEach((order) => {
    let sellCurrency = "";
    let buyCurrency = "";

    if (order.sell_curr_id === 1) sellCurrency = "🇺🇸 USD";

    if (order.sell_curr_id === 2) sellCurrency = "🇪🇺 EUR";

    if (order.sell_curr_id === 3) sellCurrency = "🇷🇺 RUB";

    if (order.buy_curr_id === 1) buyCurrency = "🇺🇸 USD";

    if (order.buy_curr_id === 2) buyCurrency = "🇪🇺 EUR";

    if (order.buy_curr_id === 3) buyCurrency = "🇷🇺 RUB";

    let actionButton = `
  <button
    onclick="openPasswordModal(${order.id}, 'buy')"
  >
    Купить
  </button>
`;

    if (order.seller_id === user.id) {
      actionButton = `
    <button
      class="logout-btn"
      onclick="openPasswordModal(${order.id}, 'cancel')"
    >
      Отменить
    </button>
  `;
    }

    ordersDiv.innerHTML += `
      <tr>

        <td>${order.id}</td>

      <td>${order.seller_name}</td>

        <td>${sellCurrency}</td>

        <td>${buyCurrency}</td>

        <td>${order.amount}</td>

        <td>${order.price}</td>

        <td>${actionButton}</td>

      </tr>
    `;
  });
}
/* BUY ORDER */

async function buyMarketOrder() {
  const order_id = document.getElementById("orderId").value;

  const response = await fetch("http://127.0.0.1:8080/market/buy", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      buyer_id: user.id,
      order_id: Number(order_id),
    }),
  });

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  showNotification("✅ Ордер куплен");

  loadMarketOrders();

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
//* THEME */

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
async function cancelOrder(orderId, password) {
  const response = await fetch("http://127.0.0.1:8080/market/cancel", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      order_id: orderId,
      user_id: user.id,
      password: password,
    }),
  });

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  const data = await response.text();

  showNotification(data);

  loadMarketOrders();

  loadWallets();
}
async function buyById(orderId, password) {
  const response = await fetch("http://127.0.0.1:8080/market/buy", {
    method: "POST",

    headers: {
      "Content-Type": "application/json",
    },

    body: JSON.stringify({
      buyer_id: user.id,
      order_id: orderId,
      password: password,
    }),
  });

  if (!response.ok) {
    const error = await response.text();

    showNotification(error);

    return;
  }

  showNotification("✅ Ордер куплен");

  loadMarketOrders();

  loadWallets();
}

let currentOrderId = null;
let currentAction = null;
function openPasswordModal(orderId, action) {
  currentOrderId = orderId;

  currentAction = action;

  document.getElementById("modalPassword").value = "";

  document.getElementById("passwordModal").style.display = "flex";
}

function closePasswordModal() {
  document.getElementById("passwordModal").style.display = "none";
}
async function confirmPassword() {
  const password = document.getElementById("modalPassword").value;

  if (!password) {
    showNotification("Введите пароль");

    return;
  }

  if (currentAction === "buy") {
    await buyById(currentOrderId, password);
  }

  if (currentAction === "cancel") {
    await cancelOrder(currentOrderId, password);
  }

  closePasswordModal();
}
/* START */

loadMarketOrders();

loadWallets();

/* AUTO REFRESH */

setInterval(() => {
  loadMarketOrders();
}, 3000);
