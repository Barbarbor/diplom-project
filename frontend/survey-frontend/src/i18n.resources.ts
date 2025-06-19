export const resources = {
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
      contacts: {
        title: "Contacts",
        section1: {
          intro:
            "If you have any questions about the service, suggestions, or complaints, you can always contact us in a convenient way:",
          email: {
            label: "Email:",
            address: "support@surveyplatform.com",
          },
          phone: {
            label: "Phone:",
            number: "+1 (555) 123-45-67",
          },
          address: {
            label: "Office Address:",
            details: "123 Example St, Suite 100, New York, NY 10001",
          },
          hours: {
            label: "Working Hours:",
            schedule: "Mon–Fri: 09:00–18:00 (EST)\nSat–Sun: Closed",
          },
        },
        section2: {
          title: "Feedback via Form",
          intro:
            "You can also send us a message directly from the website by filling out the feedback form:",
          form: {
            name: {
              label: "Your Name",
              placeholder: "John Doe",
            },
            email: {
              label: "Email",
              placeholder: "john@example.com",
            },
            message: {
              label: "Message",
              placeholder: "Your message",
            },
            submit: "Send",
          },
        },
      },
      landing: {
        layout: {
          feedback: "Feedback",
          contact_email: "Contact us: ",
          privacy: "Privacy Policy",
          terms: "Terms of Use",
          contacts: "Contacts",
          copyright: "© 2025 SurveyPlatform. All rights reserved.",
        },
        hero: {
          title: "Welcome to SurveyPlatform!",
          description:
            "Create, distribute, and analyze surveys with ease. Our service helps you gather feedback, conduct research, and improve your projects.",
          cta_logged_in: "Go to Surveys List",
          cta_logged_out: "Get Started Now",
        },
        advantages: {
          title: "Why Choose Us?",
          analytics: {
            title: "Detailed Analytics",
            description: "Get deep insights into survey results in real-time.",
          },
          ease: {
            title: "Ease of Use",
            description:
              "Intuitive interface for creating surveys without technical skills.",
          },
          security: {
            title: "Security",
            description: "Reliable protection of your respondents' data.",
          },
        },
        testimonials: {
          title: "What Our Users Say",
          testimonial1: {
            quote: '"Great tool for surveys! Very convenient and fast."',
            author: "— Anna S., Marketer",
          },
          testimonial2: {
            quote:
              '"Analytics helped us improve our product. Highly recommend!"',
            author: "— Ivan P., Developer",
          },
        },
      },
      survey: {
        surveyTitle: "Survey Title",
        creator: "Creator",
        createdAt: "Created At",
        state: "State",
        unknown: "Unknown",
        goToStats: "Go to Statistics",
        states: {
          DRAFT: "Draft",
          ACTIVE: "Active",
        },
        stats: {
          backToSurvey: "Back to Survey",
          generalStatsTitle: "General Statistics",
          startedInterviews: "Started interviews: {{count}}",
          completedInterviews: "Completed interviews: {{count}}",
          completionPercentage: "Completion percentage: {{percentage}}%",
          averageCompletionTime: "Average completion time: {{time}}s",
          error: "Error: {{message}}",
          noData: "Statistics data unavailable",
          type: "Type: {{type}}",
          answers: "Answers: {{answers}}",
        },
        question: {
          extraParams: {
            title: "Additional Parameters",
            required: "Required",
            minAnswersCount: "Minimum Answers Count",
            maxAnswersCount: "Maximum Answers Count",
            maxLength: "Maximum Length",
            minNumber: "Minimum Number",
            maxNumber: "Maximum Number",
            starsCount: "Number of Stars",
            minDate: "Minimum Date",
            maxDate: "Maximum Date",
          },
        },
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
      contacts: {
        title: "Контакты",
        section1: {
          intro:
            "Если у вас возникли вопросы по работе сервиса, предложения или жалобы, вы всегда можете связаться с нами удобным способом:",
          email: {
            label: "Электронная почта:",
            address: "support@surveyplatform.ru",
          },
          phone: {
            label: "Телефон:",
            number: "+7 (495) 123-45-67",
          },
          address: {
            label: "Адрес офиса:",
            details: "127006, г. Москва, ул. Примерная, д. 10, офис 5",
          },
          hours: {
            label: "График работы:",
            schedule: "Пн–Пт: 09:00–18:00 (МСК)\nСб–Вс: выходной",
          },
        },
        section2: {
          title: "Обратная связь через форму",
          intro:
            "Вы также можете отправить нам сообщение прямо с сайта, заполнив форму обратной связи:",
          form: {
            name: {
              label: "Ваше имя",
              placeholder: "Иван Иванов",
            },
            email: {
              label: "Электронная почта",
              placeholder: "ivan@example.com",
            },
            message: {
              label: "Сообщение",
              placeholder: "Ваше сообщение",
            },
            submit: "Отправить",
          },
        },
      },
      landing: {
        layout: {
          feedback: "Обратная связь",
          contact_email: "Свяжитесь с нами: ",
          privacy: "Политика конфиденциальности",
          terms: "Условия пользования",
          contacts: "Контакты",
          copyright: "© 2025 SurveyPlatform. Все права защищены.",
        },
        hero: {
          title: "Добро пожаловать в SurveyPlatform!",
          description:
            "Создавайте, распространяйте и анализируйте опросы с лёгкостью. Наш сервис помогает собирать обратную связь, проводить исследования и улучшать ваши проекты.",
          cta_logged_in: "Перейти к списку опросов",
          cta_logged_out: "Начать сейчас",
        },
        advantages: {
          title: "Почему выбирают нас?",
          analytics: {
            title: "Детальная аналитика",
            description:
              "Получайте глубокий анализ результатов опросов в реальном времени.",
          },
          ease: {
            title: "Простота использования",
            description:
              "Интуитивный интерфейс для создания опросов без технических навыков.",
          },
          security: {
            title: "Безопасность",
            description: "Надёжная защита данных ваших респондентов.",
          },
        },
        testimonials: {
          title: "Что говорят наши пользователи",
          testimonial1: {
            quote: '"Отличный инструмент для опросов! Очень удобно и быстро."',
            author: "— Анна С., Маркетолог",
          },
          testimonial2: {
            quote: '"Аналитика помогла нам улучшить продукт. Рекомендую!"',
            author: "— Иван П., Разработчик",
          },
        },
      },
      survey: {
        surveyTitle: "Название опроса",
        creator: "Автор",
        createdAt: "Дата создания",
        state: "Состояние",
        unknown: "Неизвестно",
        goToStats: "Перейти к статистике",
        states: {
          DRAFT: "Черновик",
          ACTIVE: "Активный",
        },
        stats: {
          backToSurvey: "Вернуться к опросу",
          generalStatsTitle: "Общая статистика",
          startedInterviews: "Начато интервью: {{count}}",
          completedInterviews: "Завершено интервью: {{count}}",
          completionPercentage: "Процент завершения: {{percentage}}%",
          averageCompletionTime: "Среднее время прохождения анкеты: {{time}}с",
          error: "Ошибка: {{message}}",
          noData: "Данные статистики недоступны",
          type: "Тип: {{type}}",
          answers: "Ответы: {{answers}}",
        },
        question: {
          extraParams: {
            title: "Дополнительные параметры",
            required: "Обязательный",
            minAnswersCount: "Минимальное количество ответов",
            maxAnswersCount: "Максимальное количество ответов",
            maxLength: "Максимальная длина",
            minNumber: "Минимальное число",
            maxNumber: "Максимальное число",
            starsCount: "Количество звёзд",
            minDate: "Минимальная дата",
            maxDate: "Максимальная дата",
          },
        },
      },
    },
  },
};
