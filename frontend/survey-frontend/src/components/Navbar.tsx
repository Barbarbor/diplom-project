'use client';
import { useTranslation } from 'next-i18next';
import Profile from './Profile';

const Navbar = ({ withProfile }: { withProfile?: boolean }) => {
  const { t } = useTranslation();

  return (
    <nav className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <div className="flex space-x-4">
        <a href="/create-poll" className="hover:underline">
          {t('auth.create_poll')}
        </a>
        <a href="/polls" className="hover:underline">
          {t('auth.polls_list')}
        </a>
      </div>
      <Profile withProfile={withProfile} />
    </nav>
  );
};

export default Navbar;
