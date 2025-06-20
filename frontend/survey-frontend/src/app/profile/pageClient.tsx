'use client';

import React, { useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'next-i18next';
import Input from '@/components/common/Input';
import Select from '@/components/common/Select';
import Button from '@/components/common/Button';
import { saveUserProfile } from '@/api-client/profile';

type FormData = {
  first_name: string;
  last_name: string;
  birth_date: string;
  phone_number: string;
  lang: 'en' | 'ru';
};

type ProfilePageClientProps = {
  initialData: FormData;
};

const formatDateToYMD = (isoDate: string): string => {
  if (!isoDate) return '';
  const date = new Date(isoDate);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

export default function ProfilePageClient({ initialData }: ProfilePageClientProps) {
  const { t } = useTranslation('translation', { keyPrefix: 'profile' });
  const [isSuccessModalOpen, setIsSuccessModalOpen] = useState(false);

  const formattedInitialData = {
    ...initialData,
    birth_date: formatDateToYMD(initialData.birth_date),
  };

  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    defaultValues: formattedInitialData,
  });

  const onSubmit = async (data: FormData) => {
    try {
      const formattedData = {
        ...data,
        birth_date: formatDateToYMD(data.birth_date),
      };
      if (!formattedData.birth_date || formattedData.birth_date === '1970-01-01') {
        formattedData.birth_date = '';
      }
      await saveUserProfile(formattedData);
      setIsSuccessModalOpen(true); // показать модалку
    } catch (error: any) {
      console.error('Failed to save profile:', error);
      alert(t('saveError', { message: error.message }));
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-200 relative">
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="bg-white p-6 rounded-lg shadow-md w-full max-w-lg"
      >
        <h1 className="text-xl font-bold mb-4">{t('title')}</h1>

        <Input
          label={t('firstName')}
          type="text"
          name="first_name"
          register={register}
          errors={errors}
        />

        <Input
          label={t('lastName')}
          type="text"
          name="last_name"
          register={register}
          errors={errors}
        />

        <Input
          label={t('birthDate')}
          type="date"
          name="birth_date"
          register={register}
          errors={errors}
        />

        <Input
          label={t('phoneNumber')}
          type="tel"
          name="phone_number"
          register={register}
          errors={errors}
        />

        <Controller
          name="lang"
          control={control}
          render={({ field }) => (
            <Select
              label={t('language')}
              name={field.name}
              value={field.value}
              onChange={field.onChange}
              options={[
                { value: 'en', label: t('language_en') },
                { value: 'ru', label: t('language_ru') },
              ]}
            />
          )}
        />

        <Button>{t('save')}</Button>
      </form>

      {/* ✅ Успешная модалка */}
      {isSuccessModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg shadow-md max-w-sm w-full text-center">
            <h2 className="text-lg font-semibold mb-4">
              {t('saveSuccess') || 'Изменения успешно сохранены!'}
            </h2>
            <button
              className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
              onClick={() => setIsSuccessModalOpen(false)}
            >
              OK
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
