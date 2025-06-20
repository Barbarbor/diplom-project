"use client";
import React from "react";
import { useTranslation } from "next-i18next";
import Image from 'next/image';
import Profile from "../Profile";
import { useRouter } from "next/navigation";
import { createSurvey } from "@/api-client/survey";
import LanguageSelect from "./Language";

const Navbar = ({ withProfile = false }: { withProfile?: boolean }) => {
  const { t } = useTranslation();
  const router = useRouter();

  const handleCreateSurvey = async () => {
    const response = await createSurvey();
    if (response.status >= 400) {
      console.error("Failed to create survey:", response.error);
    } else if (response.data?.hash) {
      router.push(`/survey/${response.data.hash}`);
    }
  };

  return (
    <nav className="bg-gray-100 text-black p-4 flex justify-between items-center">
      <div className="flex items-center space-x-4">
        <div className="relative w-[150px] h-[55px]">
          <a href="/landing">
            <Image
              src="/logo.jpg"
              alt="Logo"
              layout="fill"
              objectFit="cover"
            />
          </a>
        </div>
        <button onClick={handleCreateSurvey} className="hover:underline">
          {t("auth.create_survey")}
        </button>
        <a href="/surveyslist" className="hover:underline">
          {t("auth.surveys_list")}
        </a>
      </div>
      <div className="flex items-center space-x-4">
        <LanguageSelect isUserLogged={withProfile} />
        <Profile withProfile={withProfile} />
      </div>
    </nav>
  );
};

export default Navbar;