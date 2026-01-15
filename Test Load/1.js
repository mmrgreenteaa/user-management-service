import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

// Настройки нагрузки
export const options = {
  stages: [
    { duration: '10s', target: 5 }, // Разгон до 5 пользователей
    { duration: '5s', target: 0 },  // Остановка
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% запросов быстрее 500мс
  },
};

// Базовые заголовки (Content-Type)
const params = {
  headers: {
    'Content-Type': 'application/json',
  },
};

export default function () {
  // --- ШАГ 1: Генерация уникальных данных ---
  
  // Используем текущее время и ID вузера (__VU) для уникальности
  const uniqueId = `${Date.now()}_${__VU}_${__ITER}`;
  
  // Формируем payload на основе вашего примера
  const payloadData = {
    login: `user_${uniqueId}`,           // Пример: user_16788822_1_1
    password: "passwoвrdd123",           // Пароль можно оставить одинаковым
    email: `user_${uniqueId}@example.com` // Пример: user_16788822_1_1@example.com
  };
  
  const payload = JSON.stringify(payloadData);

  // --- ШАГ 2: Регистрация ---
  // Замените URL на ваш реальный эндпоинт регистрации
  const registerRes = http.post('http://localhost:8080/api/v1/user-manager/users', payload, params);

  check(registerRes, {
    'Registration success (201)': (r) => r.status === 201 || r.status === 200,
  });

  // --- ШАГ 3: Аутентификация (Логин) ---
  // Обычно для логина отправляют те же данные
  const loginRes = http.post('http://localhost:8080/api/v1/auth/login', payload, params);

  // Проверяем, что логин прошел успешно
  const isLoginSuccessful = check(loginRes, {
    'Login success (200)': (r) => r.status === 200,
    'Has token': (r) => r.json('token') !== undefined, // Проверяем наличие поля token
  });

  // Если логин не прошел, нет смысла продолжать (прерываем итерацию)
  if (!isLoginSuccessful) {
    console.error(`Login failed for ${payloadData.login}`);
    return; 
  }

 let authToken = loginRes.headers['Authorization'];

  // Отладка: если токен не найден, выводим заголовки
  if (!authToken) {
    console.log(`No Auth token found. Headers: ${JSON.stringify(loginRes.headers)}`);
  }

  // --- 5. Авторизованный запрос ---
  
  // Если ваш сервер отдает в заголовке чистый токен (без "Bearer "), 
  // то при отправке обычно нужно добавить префикс.
  // Если сервер сразу отдает "Bearer <token>", то префикс добавлять не надо.
  // Предположим, сервер отдает чистый токен, тогда добавляем "Bearer ":
  /*
  const authHeaders = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': authToken, // Или `Bearer ${authToken}`, зависит от того, что именно кладет Golang
    },
  };
 
  // Делаем запрос к защищенному ресурсу
  const profileRes = http.get('https://your-api.com/api/user/profile', authParams);

  check(profileRes, {
    'Authenticated request success': (r) => r.status === 200,
    'Profile email matches': (r) => r.json('email') === payloadData.email,
  }); */

  sleep(1);
}