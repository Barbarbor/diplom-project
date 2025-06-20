// app/contacts/page.tsx
import React from "react";
import { getTranslations } from "@/i18n.server";

export default async function ContactsPage() {
  const { t } = await getTranslations("translation");

  return (
    <div className="max-w-3xl mx-auto py-12 px-4 space-y-6">
      <h1 className="text-3xl font-bold mb-4">{t("contacts.title")}</h1>

      <section className="space-y-4">
        <p>{t("contacts.section1.intro")}</p>
        <ul className="list-disc list-inside ml-6 space-y-2">
          <li>
            <strong>{t("contacts.section1.email.label")}</strong>{" "}
            <a
              href={`mailto:${t("contacts.section1.email.address")}`}
              className="text-blue-600 hover:underline"
            >
              {t("contacts.section1.email.address")}
            </a>
          </li>
          <li>
            <strong>{t("contacts.section1.phone.label")}</strong>{" "}
            {t("contacts.section1.phone.number")}
          </li>
          <li>
            <strong>{t("contacts.section1.address.label")}</strong>
            <br />
            {t("contacts.section1.address.details")}
          </li>
          <li>
            <strong>{t("contacts.section1.hours.label")}</strong>
            <br />
            {t("contacts.section1.hours.schedule")}
          </li>
        </ul>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">{t("contacts.section2.title")}</h2>
        <p>{t("contacts.section2.intro")}</p>
        <form className="space-y-4 max-w-md">
          <div>
            <label htmlFor="name" className="block text-sm font-medium">
              {t("contacts.section2.form.name.label")}
            </label>
            <input
              id="name"
              name="name"
              type="text"
              placeholder={t("contacts.section2.form.name.placeholder")}
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label htmlFor="email" className="block text-sm font-medium">
              {t("contacts.section2.form.email.label")}
            </label>
            <input
              id="email"
              name="email"
              type="email"
              placeholder={t("contacts.section2.form.email.placeholder")}
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label htmlFor="message" className="block text-sm font-medium">
              {t("contacts.section2.form.message.label")}
            </label>
            <textarea
              id="message"
              name="message"
              rows={4}
              className="mt-1 block w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder={t("contacts.section2.form.message.placeholder")}
              required
            />
          </div>

          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            {t("contacts.section2.form.submit")}
          </button>
        </form>
      </section>
    </div>
  );
}