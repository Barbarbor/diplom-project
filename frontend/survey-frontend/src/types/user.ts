export enum Lang {
  Ru = "ru",
  Eng = "en",
}

export interface User {
  id: number;
  email: string;
  password: string;
  created_at: Date;
}

export interface UserProfile {
  id: number;
  user_id: number;
  first_name?: string;
  last_name?: string;
  birth_date: Date;
  phone_number: string;
  lang: Lang;
  created_at?: Date;
}

export interface ApiUserCredentials {
  email: string;
  password: string;
}

export interface ApiUser {
  email: string;
  user_id: number;
  first_name?: string;
  last_name?: string;
  birth_date: Date;
  phone_number: string;
  lang: Lang;
  created_at?: Date;
}
