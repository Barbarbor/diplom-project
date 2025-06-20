
import { fetchUserProfile } from '@/api-client/profile';
import ProfilePageClient from './pageClient';

export default async function ProfilePageServer() {
  let profileData =  {
    first_name: '',
    last_name: '',
    birth_date: '',
    phone_number: '',
    lang: 'ru'
  };

  try {
    const response = await fetchUserProfile();
    // @ts-expect-error Property 'profile' does not exist on type '{ first_name?: string | undefined; last_name?: string | undefined; birth_date?: string | undefined; phone_number?: string | undefined; lang?: string | undefined; }'.
    const profile = await response.data.profile;

    profileData = {
      first_name: profile?.first_name || '',
      last_name: profile?.last_name || '',
      birth_date: profile?.birth_date?.Time || '',
      phone_number: profile?.phone_number || '',
      lang: profile?.lang as 'ru' | 'en' || 'ru',
    };
  } catch (error) {
    console.error('Failed to load profile on server:', error);
  }

  return <ProfilePageClient initialData={profileData} />;
}