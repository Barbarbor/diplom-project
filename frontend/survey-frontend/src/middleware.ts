import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { checkIsUserLogged } from "./api-client/auth";

export async function middleware(request: NextRequest) {
  const isUserLogged = await checkIsUserLogged();

  // Проверяем, посещает ли пользователь страницы /login или /register
  const url = request.nextUrl.clone();
  if (
    (url.pathname === "/login" || url.pathname === "/register") &&
    isUserLogged
  ) {
    url.pathname = "/"; // Перенаправляем на главную
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

// Указываем, на какие маршруты применять middleware
export const config = {
  matcher: ["/login", "/register"],
};
