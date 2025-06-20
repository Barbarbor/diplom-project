import { useState, useEffect } from 'react';
import { useTranslation } from 'next-i18next';
import Cookies from 'js-cookie';

export default function Profile({ withProfile }: { withProfile?: boolean }) {
  const { t } = useTranslation();
  const [isPopupVisible, setIsPopupVisible] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const togglePopup = () => setIsPopupVisible((v) => !v);

  const handleLogout = () => {
    Cookies.remove('auth_token');
    setIsPopupVisible(false);
    location.reload();
  };

  // Avoid SSR/client mismatch by not rendering until mounted
  if (!mounted) return null;

  return (
    <div>
      {withProfile ? (
        <div className="relative">
          <button className="focus:outline-none" onClick={togglePopup}>
            {t('auth.profile')}
          </button>
          {isPopupVisible && (
            <div className="absolute right-0 mt-2 w-48 bg-white text-black shadow-lg rounded-lg">
              <a href="/profile" className="block px-4 py-2 hover:bg-gray-200">
                {t('auth.personal_account')}
              </a>
              <button
                onClick={handleLogout}
                className="w-full text-left px-4 py-2 hover:bg-gray-200 focus:outline-none"
              >
                {t('auth.logout')}
              </button>
            </div>
          )}
        </div>
      ) : (
        <a href="/login" className="hover:underline">
          {t('profile.login')}
        </a>
      )}
    </div>
  );
}
