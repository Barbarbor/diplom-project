/* eslint-disable react/no-unescaped-entities */
import React from 'react';
import Image from 'next/image'; // Импортируем Image из next/image
import { FaChartPie, FaComments, FaUserCheck } from 'react-icons/fa'; // Иконки для преимуществ
import { checkIsUserLogged } from '@/api-client/auth';

export default async function LandingPage() {
  // Вызов серверной функции для проверки авторизации
  const isUserLogged = await checkIsUserLogged();

  return (
    <div className="min-h-screen bg-gray-100">
    

      {/* Hero Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto text-center">
          <h2 className="text-4xl font-bold mb-4">Добро пожаловать в SurveyPlatform!</h2>
          <p className="text-lg text-gray-600 mb-6">
            Создавайте, распространяйте и анализируйте опросы с лёгкостью. Наш сервис помогает собирать обратную связь, проводить исследования и улучшать ваши проекты.
          </p>
          <div className="relative mx-auto mb-4 w-[600px] h-[300px]">
            <Image
              src="/graph.png"
              alt="Survey Dashboard"
              layout="fill"
              objectFit="contain"
           
              priority // Улучшаем LCP
            />
          </div>
          <a
            href={isUserLogged ? '/surveyslist' : '/login'}
            className="px-6 py-3 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            {isUserLogged ? 'Перейти к списку опросов' : 'Начать сейчас'}
          </a>
        </div>
      </section>

      {/* Advantages Section */}
      <section className="py-16 bg-gray-200">
        <div className="max-w-7xl mx-auto text-center">
          <h3 className="text-3xl font-semibold mb-8">Почему выбирают нас?</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="p-6 bg-white rounded-lg shadow-md">
              <FaChartPie className="text-blue-600 mx-auto mb-4 w-12 h-12" />
              <h4 className="text-xl font-semibold mb-2">Детальная аналитика</h4>
              <p className="text-gray-600">Получайте глубокий анализ результатов опросов в реальном времени.</p>
            </div>
            <div className="p-6 bg-white rounded-lg shadow-md">
              <FaComments className="text-blue-600 mx-auto mb-4 w-12 h-12" />
              <h4 className="text-xl font-semibold mb-2">Простота использования</h4>
              <p className="text-gray-600">Интуитивный интерфейс для создания опросов без технических навыков.</p>
            </div>
            <div className="p-6 bg-white rounded-lg shadow-md">
              <FaUserCheck className="text-blue-600 mx-auto mb-4 w-12 h-12" />
              <h4 className="text-xl font-semibold mb-2">Безопасность</h4>
              <p className="text-gray-600">Надёжная защита данных ваших респондентов.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto text-center">
          <h3 className="text-3xl font-semibold mb-8">Что говорят наши пользователи</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="p-6 bg-gray-100 rounded-lg shadow-md">
              <div className="relative mx-auto mb-2 w-[100px] h-[100px]">
                <Image
    src="/user1.jpg"
    alt="User 1"
    width={100}
    height={100}
    objectFit="cover"
    className="rounded-full"
  />
              </div>
              <p className="text-gray-600 italic mb-2">"Отличный инструмент для опросов! Очень удобно и быстро."</p>
              <p className="font-semibold">— Анна С., Маркетолог</p>
            </div>
            <div className="p-6 bg-gray-100 rounded-lg shadow-md">
             <div className="mx-auto mb-2 w-[100px] h-[100px] rounded-full overflow-hidden">
    <Image
      src="/user2.jpg"
      alt="User 2"
      width={100}
      height={100}
      objectFit="cover"
      className="rounded-full"
    />
  </div>
              <p className="text-gray-600 italic mb-2">"Аналитика помогла нам улучшить продукт. Рекомендую!"</p>
              <p className="font-semibold">— Иван П., Разработчик</p>
            </div>
          </div>
        </div>
      </section>

    
    </div>
  );
}