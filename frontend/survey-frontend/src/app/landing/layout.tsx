// app/layout.tsx
import React from "react";
import Image from "next/image";
import Link from "next/link";
import { getTranslations } from "@/i18n.server";

export const metadata = {
  title: "SurveyPlatform",
  description: "Платформа для создания и анализа опросов",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { t } = await getTranslations("ru", "translation");

  return (
    <div className="flex flex-col min-h-screen">
      {/* Header */}
      <header className="bg-gray-100 text-white p-4">
        <div className="max-w-7xl mx-auto flex items-center">
          <div className="relative w-[150px] h-[55px]">
            <Link href="/landing">
              <Image
                src="/logo.jpg"
                alt={t("landing.layout.logo_alt") || "Logo"} // Добавим перевод для alt, если нужно
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
          <h4 className="text-lg font-semibold mb-2">{t("landing.layout.feedback")}</h4>
          <p className="mb-2">
            {t("landing.layout.contact_email")}
            <a
              href="mailto:support@surveyplatform.ru"
              className="underline"
            >
              support@surveyplatform.ru
            </a>
          </p>
          <div className="flex justify-center gap-4 text-sm underline">
            <a href="/privacy">{t("landing.layout.privacy")}</a>
            <a href="/terms">{t("landing.layout.terms")}</a>
            <a href="/contacts">{t("landing.layout.contacts")}</a>
          </div>
          <p>{t("landing.layout.copyright")}</p>
        </div>
      </footer>
    </div>
  );
}