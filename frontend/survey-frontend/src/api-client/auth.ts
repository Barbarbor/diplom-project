import serverRequest from "@/lib/serverApi";
import request from "@/lib/api";

import { ApiUserCredentials, ApiUser } from "@/types/user";

// Регистрация пользователя
export const registerUser = async (data: ApiUserCredentials) => {
  const response = await request({
    method: "POST",
    prefix: "/api",
    url: "/auth/register",
    data,
    disableAuthCookie: true,
    cache: { disabled: true },
  });

  return response;
};

// Логин пользователя
export const loginUser = async (data: ApiUserCredentials) => {
  const response = await request({
    method: "POST",
    prefix: "/api",
    url: "/auth/login",
    data,
    disableAuthCookie: true,
    cache: { disabled: true },
  });

  return response;
};

export const loginUserClient = async (data: ApiUserCredentials) => {
  const response = await fetch("/api/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
  return await response.json();
};

// Получение данных пользователя
export const getUser = async () => {
  const response = await serverRequest<{ user: ApiUser }>({
    method: "GET",
    prefix: "/api",
    url: "/auth/user",
    cache: { disabled: true },
  });

  return response;
};

export const checkIsUserLogged = async () => {
  const user = await getUser();

  return !!user?.data?.user;
};
