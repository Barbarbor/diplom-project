import request, { ApiResponse } from "@/lib/api";

/**
 * Response type for access list API
 */
interface AccessResponse {
  emails: string[];
}

/**
 * Retrieves the list of users with 'edit' access for a survey (GET /api/surveys/:hash/access).
 */
export const getAccessList = async (hash: string): Promise<ApiResponse<AccessResponse>> => {
  return await request<AccessResponse>({
    method: "GET",
    prefix: "/api",
    url: `/surveys/${hash}/access`,
  });
};

/**
 * Adds 'edit' access to a user for a survey (POST /api/surveys/:hash/access?email=<userEmailToAdd>).
 */
export const addEditAccess = async (hash: string, email: string): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "POST",
    prefix: "/api",
    url: `/surveys/${hash}/access?email=${encodeURIComponent(email)}`,
  });
};

/**
 * Removes 'edit' access from a user for a survey (DELETE /api/surveys/:hash/access?email=<userEmailToDelete>).
 */
export const removeEditAccess = async (hash: string, email: string): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "DELETE",
    prefix: "/api",
    url: `/surveys/${hash}/access?email=${encodeURIComponent(email)}`,
  });
};