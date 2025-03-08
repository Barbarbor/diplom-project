import serverRequest from '@/lib/api';

export const fetchUserProfile = async () => {
  return serverRequest<{ firstName: string; lastName: string; birthDate: string; phoneNumber: string; language: string }>({
    method: 'GET',
    prefix: '/api',
    url: '/profile',
  });
};

export const saveUserProfile = async (data: {
  firstName: string;
  lastName: string;
  birthDate: string;
  phoneNumber: string;
  language: string;
}) => {
  return serverRequest({
    method: 'POST',
    prefix: '/api',
    url: '/profile',
    data,
  });
};