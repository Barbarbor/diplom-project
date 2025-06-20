// i18n.client.ts
"use client";

import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import { resources } from "./i18n.resources";

  i18n
    .use(initReactI18next)
    .use(LanguageDetector) // Опционально, для автоматического определения языка
    .init({
      resources,
      fallbackLng: "ru",
      interpolation: {
        escapeValue: false,
      },
      detection: {
      order: ["cookie"],
        caches: ["cookie"],
      },
       react: { useSuspense: false },
    });


export default i18n;
