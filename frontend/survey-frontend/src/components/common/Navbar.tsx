"use client";
import React from "react";
import { useTranslation } from "next-i18next";
import Profile from "../Profile";
import { useRouter } from "next/navigation";
import { createSurvey } from "@/api-client/survey";

const Navbar = ({ withProfile }: { withProfile?: boolean }) => {
  const { t } = useTranslation();
  const router = useRouter();

  const handleCreateSurvey = async () => {
    // Вызываем API для создания опроса
    const response = await createSurvey();
    if (response.status >= 400) {
      console.error("Failed to create survey:", response.error);
      // Можно показать уведомление об ошибке (toast)
    } else if (response.data?.hash) {
      // Перенаправляем пользователя на страницу опроса по хэшу
      router.push(`/survey/${response.data.hash}`);
    }
  };

  return (
    <nav className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <div className="flex space-x-4">
        {/* Заменяем ссылку на кнопку, которая вызывает создание опроса */}
        <button onClick={handleCreateSurvey} className="hover:underline">
          {t("auth.create_survey")}
        </button>
        <a href="/surveyslist" className="hover:underline">
          {t("auth.surveys_list")}
        </a>
      </div>
      <Profile withProfile={withProfile} />
    </nav>
  );
};

export default Navbar;