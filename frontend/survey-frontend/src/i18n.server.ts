// i18n.server.ts
import i18n from "i18next";
import { resources } from "./i18n.resources";


export async function initI18nServer(lng: string = "ru") {
  await i18n.init({
    lng,
    fallbackLng: "ru",
    resources,
    interpolation: {
      escapeValue: false,
    },
  });
  return i18n;
}

export async function getTranslations(lng: string, ns: string = "translation") {
  const i18nInstance = await initI18nServer(lng);
  return {
    t: (key: string) => i18nInstance.t(`${ns}:${key}`),
    i18n: i18nInstance,
  };
}