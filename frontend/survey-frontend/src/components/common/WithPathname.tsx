'use client';

import { usePathname } from 'next/navigation';
import React from 'react';

export default function WithPathname({
  children,
  restrictedPaths,
}: {
  children: React.ReactNode;
  restrictedPaths: string[];
}) {
  const pathname = usePathname();
  // Проверяем, начинается ли pathname с любого из restrictedPaths
  const shouldRender = !restrictedPaths?.some((restrictedPath) =>
    pathname.startsWith(restrictedPath)
  );

  // Если путь не ограничен, рендерим children, иначе ничего не рендерим
  return shouldRender ? <>{children}</> : null;
}