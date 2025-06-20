import request from '@/lib/api';
import serverRequest from '@/lib/serverApi';

export const fetchUserProfile = async () => {
  return serverRequest<{ first_name?: string; last_name?: string; birth_date?: string; phone_number?: string; lang?: string }>({
    method: 'GET',
    prefix: '/api',
    url: '/profile',
    cache:{disabled: true}
  });
};

export const saveUserProfile = async (data: {
  first_name?: string;
  last_name?: string;
  birth_date?: string;
  phone_number?: string;
  lang?: string;
}) => {
  return request({
    method: 'PUT',
    prefix: '/api',
    url: '/profile',
    data,
  });
};