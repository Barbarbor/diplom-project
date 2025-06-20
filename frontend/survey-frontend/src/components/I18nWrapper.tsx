// components/I18nWrapper.tsx
'use client';

import { ReactNode, useEffect } from 'react';
import { I18nextProvider } from 'react-i18next';
import i18n from '../i18n.client';

export default function I18nWrapper({
  children,
  language = 'ru',
}: {
  children: ReactNode;
  language: string;
}) {
  useEffect(() => {
    // Вызываем только один раз при монтировании или когда language меняется
    if (i18n.language !== language) {
      i18n.changeLanguage(language);
    }
  }, [language]);

  return <I18nextProvider i18n={i18n}>{children}</I18nextProvider>;
}