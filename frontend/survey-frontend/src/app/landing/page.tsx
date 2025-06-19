/* eslint-disable react/no-unescaped-entities */
import React from "react";
import Image from "next/image";
import { FaChartPie, FaComments, FaUserCheck } from "react-icons/fa";
import { checkIsUserLogged } from "@/api-client/auth";
import { getTranslations } from "@/i18n.server";

export default async function LandingPage() {
  const { t } = await getTranslations("ru", "translation");
  const isUserLogged = await checkIsUserLogged();

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Hero Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto text-center">
          <h2 className="text-4xl font-bold mb-4">{t("landing.hero.title")}</h2>
          <p className="text-lg text-gray-600 mb-6">{t("landing.hero.description")}</p>
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
            href={isUserLogged ? "/surveyslist" : "/login"}
            className="px-6 py-3 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            {isUserLogged ? t("landing.hero.cta_logged_in") : t("landing.hero.cta_logged_out")}
          </a>
        </div>
      </section>

      {/* Advantages Section */}
      <section className="py-16 bg-gray-200">
        <div className="max-w-7xl mx-auto text-center">
          <h3 className="text-3xl font-semibold mb-8">{t("landing.advantages.title")}</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="p-6 bg-white rounded-lg shadow-md">
              <FaChartPie className="text-blue-600 mx-auto mb-4 w-12 h-12" />
              <h4 className="text-xl font-semibold mb-2">{t("landing.advantages.analytics.title")}</h4>
              <p className="text-gray-600">{t("landing.advantages.analytics.description")}</p>
            </div>
            <div className="p-6 bg-white rounded-lg shadow-md">
              <FaComments className="text-blue-600 mx-auto mb-4 w-12 h-12" />
              <h4 className="text-xl font-semibold mb-2">{t("landing.advantages.ease.title")}</h4>
              <p className="text-gray-600">{t("landing.advantages.ease.description")}</p>
            </div>
            <div className="p-6 bg-white rounded-lg shadow-md">
              <FaUserCheck className="text-blue-600 mx-auto mb-4 w-12 h-12" />
              <h4 className="text-xl font-semibold mb-2">{t("landing.advantages.security.title")}</h4>
              <p className="text-gray-600">{t("landing.advantages.security.description")}</p>
            </div>
          </div>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto text-center">
          <h3 className="text-3xl font-semibold mb-8">{t("landing.testimonials.title")}</h3>
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
              <p className="text-gray-600 italic mb-2">{t("landing.testimonials.testimonial1.quote")}</p>
              <p className="font-semibold">{t("landing.testimonials.testimonial1.author")}</p>
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
              <p className="text-gray-600 italic mb-2">{t("landing.testimonials.testimonial2.quote")}</p>
              <p className="font-semibold">{t("landing.testimonials.testimonial2.author")}</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}