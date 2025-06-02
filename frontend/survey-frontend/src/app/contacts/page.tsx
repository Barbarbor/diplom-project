// app/contacts/page.tsx
import React from "react";

export default function ContactsPage() {
  return (
    <div className="max-w-3xl mx-auto py-12 px-4 space-y-6">
      <h1 className="text-3xl font-bold mb-4">Контакты</h1>

      <section className="space-y-4">
        <p>
          Если у вас возникли вопросы по работе сервиса, предложений или жалобы, вы всегда можете связаться 
          с нами удобным способом:
        </p>
        <ul className="list-disc list-inside ml-6 space-y-2">
          <li>
            <strong>Электронная почта:</strong>{" "}
            <a href="mailto:support@surveyplatform.ru" className="text-blue-600 hover:underline">
              support@surveyplatform.ru
            </a>
          </li>
          <li>
            <strong>Телефон:</strong> +7 (495) 123-45-67
          </li>
          <li>
            <strong>Адрес офиса:</strong><br />
            127006, г. Москва, ул. Примерная, д. 10, офис 5
          </li>
          <li>
            <strong>График работы:</strong><br />
            Пн–Пт: 09:00–18:00 (МСК)<br />
            Сб–Вс: выходной
          </li>
        </ul>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">Обратная связь через форму</h2>
        <p>
          Вы также можете отправить нам сообщение прямо с сайта, заполнив форму обратной связи:
        </p>
        <form className="space-y-4 max-w-md">
          <div>
            <label htmlFor="name" className="block text-sm font-medium">
              Ваше имя
            </label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder="Иван Иванов"
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label htmlFor="email" className="block text-sm font-medium">
              Электронная почта
            </label>
            <input
              id="email"
              name="email"
              type="email"
              placeholder="ivan@example.com"
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label htmlFor="message" className="block text-sm font-medium">
              Сообщение
            </label>
            <textarea
              id="message"
              name="message"
              rows={4}
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Ваше сообщение"
              required
            />
          </div>

          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Отправить
          </button>
        </form>
      </section>
    </div>
  );
}
