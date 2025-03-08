import { NextRequest, NextResponse } from "next/server";

import { loginUser } from "@/api-client/auth";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const response = await loginUser(body);

    if (response.status >= 400) {
      const errorData = await response.error;
      return NextResponse.json(
        { error: errorData },
        { status: response.status }
      );
    }

    const nextResponse = NextResponse.json({ message: "Login successful" });

    // Устанавливаем куки из заголовка Set-Cookie
    const setCookieHeader = response?.headers?.get("set-cookie");

    if (setCookieHeader) {
      nextResponse.headers.set("Set-Cookie", setCookieHeader);
    }

    return nextResponse;
  } catch (error) {
    return NextResponse.json(
      { message: "Login failed", error: (error as Error).message },
      { status: 401 }
    );
  }
}
