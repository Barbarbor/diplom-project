import type { Metadata } from "next";
import "./globals.css";
import "../i18n";
import I18nWrapper from "@/components/I18nWrapper";
import QueryProvider from './providers';
import { checkIsUserLogged } from "@/api-client/auth";
import Navbar from "@/components/common/Navbar";
import Footer from "@/components/common/Footer";
import WithPathname from "@/components/common/WithPathname";
import CookieBanner from "@/components/banners/CookieBanner";
import { ErrorBoundary } from "@/components/common/ErrorBoundary";
import { getLanguage } from "@/lib/lang";

export const metadata: Metadata = {
  title: "Survey Platform",
  description: "Survey Platform",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const isLogged = await checkIsUserLogged();

const language = await getLanguage();

  return (
    <html lang="ru">
    <body className="flex flex-col min-h-screen">
       <ErrorBoundary>
        <QueryProvider>
          <I18nWrapper language={language} >
          
            <WithPathname restrictedPaths={['/landing', '/poll', '/privacy','/terms', '/contacts']}>
            <Navbar withProfile={isLogged} />
            </WithPathname>
            <main className="flex-grow">{children}</main>
             <WithPathname restrictedPaths={['/landing', '/poll', '/privacy','/terms', '/contacts']}>
            <Footer />
            </WithPathname>
             <WithPathname restrictedPaths={['/poll']}>
            <CookieBanner />
            </WithPathname>
          
          </I18nWrapper>
        </QueryProvider>
        </ErrorBoundary>
      </body>

    </html>
  );
}