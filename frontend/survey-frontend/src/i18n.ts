"use client";
import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      auth: {
        register: {
          title: "Register",
          submit: "Register",
        },
        login: {
          title: "Login",
          submit: "Login",
        },
        email: {
          label: "Email",
        },
        password: {
          label: "Password",
        },
        switch_to_login: "Already have an account? Login",
        switch_to_register: "No account? Register",
        create_poll: "Create Poll",
        polls_list: "Polls List",
        profile: "Profile",
        personal_account: "Personal Account",
        settings: "Settings",
        logout: "Logout",
      },
      profile: {
        title: "Profile Setup",
        firstName: "First Name",
        lastName: "Last Name",
        birthDate: "Birth Date",
        phoneNumber: "Phone Number",
        language: "Preferred Language",
        save: "Save Changes",
        login: "Sign in",
      },
    },
  },
  ru: {
    translation: {
      auth: {
        register: {
          title: "Регистрация",
          submit: "Зарегистрироваться",
        },
        login: {
          title: "Вход",
          submit: "Войти",
        },
        email: {
          label: "Эл. почта",
        },
        password: {
          label: "Пароль",
        },
        switch_to_login: "Уже есть аккаунт? Войти",
        switch_to_register: "Нет аккаунта? Зарегистрироваться",
        create_survey: "Создать опрос",
        surveys_list: "Список опросов",
        profile: "Профиль",
        personal_account: "Личный кабинет",
        settings: "Настройки",
        logout: "Выйти",
      },
      profile: {
        title: "Настройка Профиля",
        firstName: "Имя",
        lastName: "Фамилия",
        birthDate: "Дата Рождения",
        phoneNumber: "Номер Телефона",
        language: "Предпочитаемый Язык",
        save: "Сохранить Изменения",
        login: "Войти",
      },
    },
  },
};

i18n.use(initReactI18next).init({
  resources,
  lng: "ru",
  fallbackLng: "ru",
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
