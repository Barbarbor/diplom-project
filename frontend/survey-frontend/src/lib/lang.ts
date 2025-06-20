import { cookies } from "next/headers";
import { fetchUserProfile } from "@/api-client/profile";

export const getLanguage = async (): Promise<string> => {
  try {
    // Fetch user profile to get the language
    const profileResponse = await fetchUserProfile();
    // @ts-expect-error Property 'profile' does not exist on type '{ firstName: string; lastName: string; birthDate: string; phoneNumber: string; language: string; }'.
    if (profileResponse.status >= 200 && profileResponse.data?.profile?.lang) {
      // @ts-expect-error Property 'profile' does not exist on type '{ firstName: string; lastName: string; birthDate: string; phoneNumber: string; language: string; }'.
      return profileResponse.data?.profile?.lang;
    }
  } catch (error) {
    console.error("Failed to fetch user profile:", error);
  }

  // Fallback to cookie if profile fetch fails or no language is set
  const cookieStore = await cookies();
  const i18nextLng = cookieStore.get("i18next")?.value;

  if (i18nextLng) {
    return i18nextLng;
  }

  // Default to "ru" if no language is found
  return "ru";
};
