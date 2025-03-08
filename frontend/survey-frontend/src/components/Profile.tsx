'use client';
import { useState } from 'react';
import { useTranslation } from 'next-i18next';

export default function Profile({ withProfile }: { withProfile?: boolean }) {
  const { t } = useTranslation();
  const [isPopupVisible, setIsPopupVisible] = useState(false);
  const togglePopup = () => setIsPopupVisible(!isPopupVisible);

  return (
    <div>
      {withProfile ? (
        <div className="relative">
          <button className="focus:outline-none" onClick={togglePopup}>
            {t('auth.profile')}
          </button>
          <div
            className="absolute right-0 mt-2 w-48 bg-white text-black shadow-lg rounded-lg"
            hidden={!isPopupVisible}
          >
            <a href="/profile" className="block px-4 py-2 hover:bg-gray-200">
              {t('auth.personal_account')}
            </a>
            <a href="/settings" className="block px-4 py-2 hover:bg-gray-200">
              {t('auth.settings')}
            </a>
            <a href="/logout" className="block px-4 py-2 hover:bg-gray-200">
              {t('auth.logout')}
            </a>
          </div>
        </div>
      ) : (
        <a href="/login" className="hover:underline">
          {t('profile.login')}
        </a>
      )}
    </div>
  );
}
