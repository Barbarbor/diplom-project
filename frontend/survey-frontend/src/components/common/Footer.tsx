"use client";
import React from "react";
import Image from "next/image";
import Link from "next/link";

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-100 text-black py-8">
      <div className="max-w-7xl mx-auto px-4 flex flex-col md:flex-row justify-between items-center space-y-6 md:space-y-0">
        {/* Логотип и копирайт */}
        <div className="flex items-center space-x-2">
          <div className="relative w-[150px] h-[55px]">
            <Link href="/landing">

                <Image
                  src="/logo.jpg"
                  alt="SurveyPlatform Logo"
                  layout="fill"
                  objectFit="cover"
                  priority
                />

            </Link>
          </div>
          <p className="text-sm">© 2025 SurveyPlatform. Все права защищены.</p>
        </div>

        {/* Ссылки */}
        <nav className="flex space-x-6">
          <Link href="/terms">
            <p>Условия пользования</p>
          </Link>
          <Link href="/privacy">
            <p >Политика конфиденциальности</p>
          </Link>
          <Link href="/contacts">
            <p >Контакты</p>
          </Link>
        </nav>
      </div>
    </footer>
  );
};

export default Footer;
