/* eslint-disable react/no-unescaped-entities */
"use client";

import { useState, useEffect } from "react";

const CookieBanner = () => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Проверяем, принимал ли пользователь условия ранее
    const hasAcceptedCookies = localStorage.getItem("cookieAccepted");
    if (!hasAcceptedCookies) {
      // Показываем баннер через 2 секунды, если это первое посещение
      const timer = setTimeout(() => {
        setIsVisible(true);
      }, 2000);

      // Очищаем таймер при размонтировании компонента
      return () => clearTimeout(timer);
    }
  }, []);

  const handleAccept = () => {
    localStorage.setItem("cookieAccepted", "true");
    setIsVisible(false);
  };

  if (!isVisible) return null;

  return (
    <div className="fixed bottom-0 left-0 w-full bg-gray-800 text-white p-4 text-center z-50">
      <p>
        Мы используем куки для улучшения работы сайта. Нажимая "Принять", вы соглашаетесь с нашей{" "}
        <a href="/privacy" className="underline text-blue-300">Политикой конфиденциальности</a>.
      </p>
      <button
        onClick={handleAccept}
        className="mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        Принять
      </button>
    </div>
  );
};

export default CookieBanner;