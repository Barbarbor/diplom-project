import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  i18n: {
    locales: ["en", "ru"],
    defaultLocale: "ru",
  },
  images: {
    domains: ['images.unsplash.com'] // Allow images from via.placeholder.com
  },
};

export default nextConfig;
