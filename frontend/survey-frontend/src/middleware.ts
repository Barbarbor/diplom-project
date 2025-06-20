import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { checkIsUserLogged } from "./api-client/auth";

export async function middleware(request: NextRequest) {
  const isUserLogged = await checkIsUserLogged();

  const url = request.nextUrl.clone();
    const response = NextResponse.next();

  // Проверяем, посещает ли пользователь страницы /login или /register
  if ((url.pathname === "/login" || url.pathname === "/register" || url.pathname === '/') && isUserLogged) {
    url.pathname = "/surveyslist";
    return NextResponse.redirect(url);
  }

  // Проверяем, если пользователь не зарегистрирован и пытается зайти на защищённые страницы
  if (!isUserLogged) {
    const protectedPaths = [
      "/surveyslist",
      "/survey/:hash",
      "/survey/:hash/stats",
      "/",
      "/profile"
    ];

    const isProtectedPath = protectedPaths.some((path) => {
      if (!path.includes(":hash")) {
        return url.pathname === path;
      }
      const regexPath = path.replace(":hash", "[^/]+");
      const regex = new RegExp(`^${regexPath}$`);
      return regex.test(url.pathname);
    });

    if (isProtectedPath) {
      url.pathname = "/landing";
      return NextResponse.redirect(url);
    }
  }

  return response;
}

export const config = {
  matcher: ["/login", "/register", "/surveyslist", "/survey/:path*", "/landing", "/poll/:path*", "/", "/profile"],
};