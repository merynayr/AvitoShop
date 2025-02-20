import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomItem, randomIntBetween, uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

// Настройки нагрузки
export let options = {
  stages: [
    // Этап 1: Увеличение нагрузки до 10 Виртуальных Пользователей
    { duration: '5s', target: 10 },

    // Этап 2: Увеличение нагрузки до 100 Виртуальных Пользователей
    { duration: '10s', target: 100 },

    // Этап 3: Увеличение нагрузки до 1000 Виртуальных Пользователей
    { duration: '15s', target: 1000 },

    // // Этап 4: Увеличение нагрузки до 10k Виртуальных Пользователей
    // { duration: '30s', target: 10000 },

    // Этап 5: Плавное снижение до нуля
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    // SLI: Время ответа должно быть менее 50 мс для 99% запросов
    http_req_duration: ['p(99)<50'],

    // Успешность ответов: 99.99% успешных запросов
    http_req_failed: ['rate<0.0001'],
  },

  rps: 1000, 
};


// API URL
const BASE_URL = 'http://localhost:8080';

// Возможные товары
const merchPrices = [
  "t-shirt", "cup", "book", "pen", "powerbank",
  "hoody", "umbrella", "socks", "wallet", "pink-hoody",
];

// Создаём случайных пользователей
function generateUser() {
  return {
    username: `user${uuidv4().substring(0, 8)}`,
    password: 'password123',
  };
}

// Авторизация
function authenticate(user) {
  const res = http.post(`${BASE_URL}/api/auth`, JSON.stringify(user), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(res, {
    'Auth successful': (r) => r.status === 200,
    'Has access token': (r) => JSON.parse(r.body).access_token !== undefined,
  });

  const { access_token, refresh_token } = JSON.parse(res.body);
  return { access_token, refresh_token };
}

// Основной тест
export default function () {
  const userA = generateUser();
  const userB = generateUser();

  const { access_token: tokenA } = authenticate(userA);
  const { access_token: tokenB } = authenticate(userB);

  const headersA = { 'Content-Type': 'application/json', "Authorization": `Bearer ${tokenA}` };
  const headersB = { 'Content-Type': 'application/json', "Authorization": `Bearer ${tokenB}` };

  // 1. Запрос информации
  const infoRes = http.get(`${BASE_URL}/api/info`, { headers: headersA });
  check(infoRes, { 'Get info success': (r) => r.status === 200 });

  // 2. Рандомный товар для покупки
  const item = randomItem(merchPrices);
  const buyRes = http.get(`${BASE_URL}/api/buy/${item}`, { headers: headersA });
  check(buyRes, { [`Buy ${item} success`]: (r) => r.status === 200 });


  // 3. Случайный перевод денег
  const amount = randomIntBetween(5, 50);
  const sendCoinRes = http.post(
    `${BASE_URL}/api/sendCoin`,
    JSON.stringify({ toUser: userB.username, amount }),
    { headers: headersA }
  );
  check(sendCoinRes, { [`Send ${amount} coins success`]: (r) => r.status === 200 });

  sleep(1);
}
