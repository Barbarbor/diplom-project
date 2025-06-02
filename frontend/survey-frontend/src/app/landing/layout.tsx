// app/layout.tsx
import React from "react";
import Image from "next/image";

import Link from "next/link";

export const metadata = {
  title: "SurveyPlatform",
  description: "Платформа для создания и анализа опросов",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {


  return (
   
      <div className="flex flex-col min-h-screen">
        {/* Header */}
        <header className="bg-gray-100 text-white p-4">
          <div className="max-w-7xl mx-auto flex items-center">
            <div className="relative w-[150px] h-[55px]">
              <Link href="/landing">

                  <Image
                    src="/logo.jpg"
                    alt="Logo"
                    layout="fill"
                    objectFit="cover"
                  />

              </Link>
            </div>
           
          </div>
        </header>

        {/* Основной контент */}
        <main className="flex-grow">{children}</main>

        {/* Footer */}
        <footer className="bg-gray-600 text-white py-6">
          <div className="max-w-7xl mx-auto text-center">
            <h4 className="text-lg font-semibold mb-2">Обратная связь</h4>
            <p className="mb-2">
              Свяжитесь с нами:{" "}
              <a
                href="mailto:support@surveyplatform.ru"
                className="underline"
              >
                support@surveyplatform.ru
              </a>
            </p>
             <div className="flex justify-center gap-4 text-sm underline">
              <a href="/privacy">Политика конфиденциальности</a>
              <a href="/terms">Условия пользования</a>
              <a href="/contacts">Контакты</a>
            </div>
            <p>© 2025 SurveyPlatform. Все права защищены.</p>
          </div>
        </footer>
      </div>

  );
}
