'use client';
import { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import Cookies from 'js-cookie';
import { saveUserProfile } from '@/api-client/profile';

interface LanguageSelectProps {
  isUserLogged: boolean;
}

const languageOptions = [
  { value: 'en', label: 'English' },
  { value: 'ru', label: 'Русский' },
];

export default function LanguageSelect({ isUserLogged }: LanguageSelectProps) {
  const { t, i18n } = useTranslation();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Only render translation after hydration
  useEffect(() => {
    setMounted(true);
  }, []);

  const handleLanguageChange = async (newLanguage: string) => {
    try {
      await i18n.changeLanguage(newLanguage);
      Cookies.set('i18nextLng', newLanguage, { expires: 365 });
      if (isUserLogged) {
        await saveUserProfile({ lang: newLanguage });
      }
      setIsDropdownOpen(false);
      window.location.reload();
    } catch (error) {
      console.error('Ошибка при смене языка:', error);
    }
  };

  const handleLabelClick = () => {
    setIsDropdownOpen(true);
  };

  // Close dropdown on outside click
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // During SSR, skip rendering label to avoid mismatch
  if (typeof window === 'undefined' || !mounted) {
    return null;
  }

  return (
    <div className="mb-4 relative" ref={dropdownRef}>
      <label
        className="block text-gray-700 font-bold mt-4 mr-24 cursor-pointer hover:underline"
        onClick={handleLabelClick}
      >
        {t('profile.language')}
      </label>
      {isDropdownOpen && (
        <ul className="absolute z-10 min-w-[200px] bg-white border border-gray-300 rounded-md shadow-lg top-full left-0 mt-1">
          {languageOptions.map((option) => (
            <li
              key={option.value}
              className="px-4 py-2 hover:bg-gray-100 cursor-pointer"
              onClick={() => handleLanguageChange(option.value)}
            >
              {option.label}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
