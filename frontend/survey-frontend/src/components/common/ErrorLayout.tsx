'use client';

import Link from 'next/link';

export default function ErrorLayout({ statusCode, message, children }: { statusCode: number; message: string; children?: React.ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="text-center p-6 bg-white rounded-lg shadow-lg">
        <h1 className="text-6xl font-bold text-red-600 mb-4">{statusCode}</h1>
        <p className="text-xl text-gray-700 mb-6">{message}</p>
        {children}
        <Link href="/" className="mt-4 inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
          Вернуться на главную
        </Link>
      </div>
    </div>
  );
}