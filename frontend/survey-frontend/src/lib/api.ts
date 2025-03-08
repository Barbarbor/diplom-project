// lib/api.ts

type ApiRequestParams = {
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  url: string;
  disableAuthCookie?: boolean;
  prefix?: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data?: Record<string, any>;
  cache?: {
    disabled?: boolean;
    revalidateTime?: number;
    tags?: string[];
  };
};

type ApiResponse<T = unknown> = {
  data?: T;
  error?: string;
  headers?: Headers;
  status: number;
};

const request = async <T = unknown>({
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

  // Если не нужно отключать использование куки, пробуем их использовать.
  if (!disableAuthCookie) {
    const authToken = document.cookie
      .split("; ")
      .find((row) => row.startsWith("auth_token="))
      ?.split("=")[1];

    if (authToken) {
      headers["Authorization"] = `Bearer ${authToken}`;
    } else {
      return {
        error: "Unauthorized",
        status: 401,
      };
    }
  }

  const fetchOptions: RequestInit = {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
    credentials: "include", // Включаем автоматическую передачу куки
    cache: cache.disabled ? "no-store" : "force-cache",
  };

  const response = await fetch(fullUrl, fetchOptions);

  if (!response.ok) {
    const responseBody = await response.json();
    return { error: responseBody.error, status: response.status };
  }

  const responseData = await response.json();
  return {
    data: responseData as T,
    status: response.status,
    headers: response.headers,
  };
};

export default request;
