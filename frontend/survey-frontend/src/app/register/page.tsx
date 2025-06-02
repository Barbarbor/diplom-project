"use client";
import React, { useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { registerUser } from "@/api-client/auth";
import Input from "@/components/common/Input";
import Button from "@/components/common/Button";
import { useRouter } from "next/navigation";
import { useTranslation } from "next-i18next";

const RegisterPage = () => {
  const { t } = useTranslation();

  const registerSchema = useMemo(
    () =>
      z.object({
        email: z.string().email({ message: t("auth.login.invalid_email") }),
        password: z
          .string()
          .min(9, { message: t("auth.login.invalid_password_length") })
          .regex(/[0-9]/, { message: t("auth.login.invalid_password_nums") })
          .regex(/[!@#$%^&*(),.?":{}|<>_+&=-]/, {
            message: t("auth.login.invalid_password_symbols"),
          }),
      }),
    [t]
  );

  type RegisterFormData = z.infer<typeof registerSchema>;
  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors },
    setError
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });

  const onSubmit = async (data: RegisterFormData) => {
    const response = await registerUser(data);
    if (response.status === 409) {
    
      setError("email", {
        type: "server",
        message: response.error,
      });
    } else {
      router.push("/login");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="bg-white p-6 rounded-lg shadow-md w-full max-w-md"
      >
        <h1 className="text-xl font-bold mb-4">Регистрация</h1>

        <Input
          label={t("auth.email.label")}
          type="text"
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

        <Button>{t("auth.register.submit")}</Button>

        <button
          type="button"
          onClick={() => router.push("/login")}
          className="mt-4 ml-4 text-blue-500 hover:underline"
        >
          {t("auth.switch_to_login")}
        </button>
      </form>
    </div>
  );
};

export default RegisterPage;
