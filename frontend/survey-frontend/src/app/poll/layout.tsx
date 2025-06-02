import React from "react";
import Link from "next/link";

export const metadata = {
  title: "SurveyPlatform",
  description: "Прохождение опросов",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
     <div className="flex flex-col min-h-screen">
      <main className="flex-grow">{children}</main>
      <footer className="bg-gray-800 text-white p-4">
        <div className="max-w-7xl mx-auto flex flex-col sm:flex-row justify-between items-center">
          <p className="mb-2 sm:mb-0">
            SurveyPlatform не имеет отношения к опросам, созданным третьими лицами, и не несёт ответственности за их содержание. (С) 2025 SurveyPlatform
          </p>
          <div className="flex space-x-4">
            <Link href="/privacy" className="hover:underline">
              Политика конфиденциальности
            </Link>
            <Link href="/terms" className="hover:underline">
              Условия использования
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}