// lib/api.ts
export type ApiRequestParams = {
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

export type ApiResponse<T = unknown> = {
  data?: T;
  error?: string;
  headers?: Headers;
  status: number;
  success?: boolean;
};

const request = async <T = unknown>({
  method,
  url,
  disableAuthCookie = false,
  prefix = "",
  data
}: ApiRequestParams): Promise<ApiResponse<T>> => {
  const fullUrl = `${process.env.NEXT_PUBLIC_API_URL}${prefix}${url}`;
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  

  const fetchOptions: RequestInit = {
    method,
    headers,
    body: data ? JSON.stringify(data) : undefined,
    credentials: !disableAuthCookie ? "include": 'omit', // Включаем автоматическую передачу куки
  
  };
try {
    const response = await fetch(fullUrl, fetchOptions);
    
    if (!response.ok) {
      const responseBody = await response.json();
      return { error: responseBody.error, status: response.status, success: false };
    }

    // Проверяем статус 204 "No Content"
    if (response.status === 204) {

      return {
        status: response.status,
        headers: response.headers,
        success: true,
      };
    }
    // Парсим тело ответа, если оно есть
    const responseData = await response.json();

    return {
      data: responseData as T,
      status: response.status,
      headers: response.headers,
      success: true,
    };
  } catch (error) {
  
    return {
      error: error instanceof Error ? error.message : 'Unknown error',
      status: 0,
      headers: new Headers(),
      success: false,
    };
  }
};

export default request;
