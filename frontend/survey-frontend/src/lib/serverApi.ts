"use server";
import { cookies } from "next/headers";

type ApiRequestParams = {
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  url: string;
  disableAuthCookie?: boolean;
  prefix?: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data?: Record<string, any>;
  cache?: {
    disabled?: boolean;
    revalidateTime?: number; // Время в секундах
    tags?: string[]; // Реалидационные теги
  };
};

export type ApiResponse<T = unknown> = {
  data?: T;
  error?: string;
  headers?: Headers;
  status: number;
};
const serverRequest = async <T = unknown>({
  method,
  url,
  disableAuthCookie = false,
  prefix = "",
  data,
  cache = {},
}: ApiRequestParams): Promise<ApiResponse<T>> => {
  const fullUrl = `${process.env.NEXT_PUBLIC_API_URL}${prefix}${url}`;
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (!disableAuthCookie) {
    const cookieStore = await cookies();
    const authToken = cookieStore.get("auth_token")?.value;
    if (!authToken) {
      return {
        error: "Unauthorized",
        status: 401,
      };
    }

    headers["Authorization"] = `Bearer ${authToken}`;
    headers["Cookie"] = `auth_token=${authToken}`;
  }

  const fetchOptions: RequestInit = {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
    credentials: "include",
    cache: cache.disabled ? "no-store" : "force-cache",
    next: {
      revalidate: cache.revalidateTime ?? undefined,
      tags: cache.tags ?? undefined,
    },
  };

  const response = await fetch(fullUrl, fetchOptions);
  const responseData = await response.json();

  if (!response.ok) {
    const responseError = responseData.error;
    return { error: responseError, status: response.status };
  }

  return {
    data: responseData,
    status: response.status,
    headers: response.headers,
  };
};

export default serverRequest;
