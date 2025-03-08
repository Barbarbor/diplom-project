import type { Metadata } from "next";
import Navbar from "@/components/Navbar";
import "./globals.css";
import "../i18n";
import { checkIsUserLogged } from "@/api-client/auth";
import I18nWrapper from "@/components/I18nWrapper";

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const isLogged = await checkIsUserLogged();

  return (
    <html lang="ru">
      <body>
        <I18nWrapper>
          <Navbar withProfile={isLogged} />
          {children}
        </I18nWrapper>
      </body>
    </html>
  );
}
