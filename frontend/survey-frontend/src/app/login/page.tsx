"use client";
import React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { loginUser } from "@/api-client/auth";
import Input from "@/components/common/Input";
import Button from "@/components/common/Button";
import { useRouter } from "next/navigation";
import { useTranslation } from "next-i18next";

// Определяем схему валидации с помощью zod для логина

interface LoginFormData {
  email: string;
  password: string;
}

const LoginPage = () => {
  const { t } = useTranslation();

  const loginSchema = z.object({
    email: z.string().email({ message: t("auth.login.invalid_email") }),
    password: z.string().min(1, { message: t("auth.login.invalid_password") }),
  });

  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });



const onSubmit = async (data: LoginFormData) => {
  const response = await loginUser(data);
  if (response.status >= 400) {
    setError("root", {
      type: "server",
      message: response.error,
    });
  } else {
    router.push("/");
    router.refresh();
  }
};

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-200">
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="bg-white p-6 rounded-lg shadow-md w-full max-w-md"
      >
        <h1 className="text-xl font-bold mb-4">{t("auth.login.title")}</h1>

        <Input
          label={t("auth.email.label")}
          type="email"
          name="email"
          register={register}
          errors={errors}
        />

        <Input
          label={t("auth.password.label")}
          type="password"
          name="password"
          register={register}
          errors={errors}
        />

        {errors.root && (
          <p className="text-red-500 text-sm text-center mb-3">
            {errors.root.message}
          </p>
        )}

        <Button>{t("auth.login.submit")}</Button>

        <button
          type="button"
          onClick={() => router.push("/register")}
          className="mt-4 ml-4  text-blue-500 hover:underline"
        >
          {t("auth.switch_to_register")}
        </button>
      </form>
    </div>
  );
};

export default LoginPage;
