'use client';

import { CustomError } from '@/lib/error';
import { ErrorBoundary as ReactErrorBoundary, FallbackProps } from 'react-error-boundary';


const FallbackComponent = ({ error, resetErrorBoundary }: FallbackProps) => {
  let title = 'Неизвестная ошибка';
  let message = 'Произошла непредвиденная ошибка.';

  // Check if error is a CustomError and handle based on statusCode
  if (error instanceof CustomError) {
    switch (error.statusCode) {
      case 404:
        title = 'Страница не найдена';
        message = `Ошибка 404: ${error.message || 'Запрошенный опрос не существует.'}`;
        break;
      case 403:
        title = 'Доступ запрещён';
        message = `Ошибка 403: ${error.message || 'У вас нет прав для просмотра этого опроса.'}`;
        break;
      case 401:
        title = 'Неавторизован';
        message = `Ошибка 401: ${error.message || 'Пожалуйста, войдите в систему.'}`;
        break;
      case 400:
        title = 'Неверный запрос';
        message = `Ошибка 400: ${error.message || 'Проверьте данные и попробуйте снова.'}`;
        break;
      case 500:
        title = 'Внутренняя ошибка сервера';
        message = `Ошибка 500: ${error.message || 'Проблема на стороне сервера. Попробуйте позже.'}`;
        break;
      default:
        title = 'Неизвестная ошибка';
        message = `Ошибка ${error.statusCode}: ${error.message || 'Неизвестная проблема.'}`;
    }
  } else {
    // Handle non-CustomError cases (e.g., unexpected JavaScript errors)
    message += ` ${error.message}`;
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="text-center p-6 bg-white rounded-lg shadow-lg">
        <h1 className="text-4xl font-bold text-red-600 mb-4">{title}</h1>
        <p className="text-lg text-gray-700 mb-4">{message}</p>
        <button
          onClick={resetErrorBoundary}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Попробовать снова
        </button>
      </div>
    </div>
  );
};

export const ErrorBoundary = ({ children }: { children: React.ReactNode }) => {
  return (
    <ReactErrorBoundary
      FallbackComponent={FallbackComponent}
      onReset={() => {
        // Optional: Reset any query or state on retry
        // e.g., queryClient.resetQueries(SURVEY_QUERY_KEY);
      }}
      onError={(error, info) => {
        console.error('Error caught by boundary:', error, info);
      }}
    >
      {children}
    </ReactErrorBoundary>
  );
};