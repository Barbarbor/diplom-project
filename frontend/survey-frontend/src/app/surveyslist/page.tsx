// app/surveyslist/page.tsx
"use client";

import React, { useState, useMemo } from "react";
import { useGetSurveys } from "@/hooks/react-query/survey";
import { useRouter } from "next/navigation";
import { SurveyState } from "@/types/survey";
import Spinner from "@/components/common/Spinner";

const SURVEY_STATE = {'DRAFT': 'Черновик', 'ACTIVE': 'Активный'}
const SurveysListPage = () => {
  const { data, isLoading, error } = useGetSurveys();
  const router = useRouter();

  // Состояние для фильтров
  const [titleFilter, setTitleFilter] = useState("");
  const [stateFilter, setStateFilter] = useState<SurveyState | "all">("all");

  const surveys = data?.surveys; // Массив опросов находится в data.surveys

  // Фильтрация опросов
  const filteredSurveys = useMemo(() => {
    if (!surveys) return [];
    return surveys.filter((survey) => {
      const matchesTitle = survey.title
        .toLowerCase()
        .includes(titleFilter.toLowerCase());
      const matchesState =
        stateFilter === "all" || survey.state === stateFilter;
      return matchesTitle && matchesState;
    });
  }, [surveys, titleFilter, stateFilter]);

  // Определяем цвет строки в зависимости от состояния
  const getStateColor = (state: SurveyState) => {
    switch (state) {
      case "ACTIVE":
        return "bg-green-50";
      case "DRAFT":
        return "bg-yellow-50";
      default:
        return "bg-white";
    }
  };

  if (isLoading) return <Spinner />;
  if (error) return <div>Ошибка при загрузке опросов: {error.message}</div>;

  return (
    <div className="max-w-4xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Список опросов</h1>

      {/* Фильтры */}
      <div className="mb-4 flex flex-col sm:flex-row gap-4">
        <div className="flex-1">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Поиск по названию
          </label>
          <input
            type="text"
            value={titleFilter}
            onChange={(e) => setTitleFilter(e.target.value)}
            placeholder="Введите название опроса"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="flex-1">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Состояние
          </label>
          <select
            value={stateFilter}
            onChange={(e) => setStateFilter(e.target.value as SurveyState | "all")}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">Все</option>
            <option value="ACTIVE">Активный</option>
            <option value="DRAFT">Черновик</option>
          </select>
        </div>
      </div>

      {!surveys || surveys.length === 0 ? (
        <p>Нет доступных опросов</p>
      ) : filteredSurveys.length === 0 ? (
        <p>Нет опросов, соответствующих фильтру</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full bg-white border border-gray-300">
            <thead>
              <tr className="bg-gray-100 border-b">
                <th className="py-2 px-4 text-left">Название</th>
                <th className="py-2 px-4 text-left">Создано</th>
                <th className="py-2 px-4 text-left">Обновлено</th>
                <th className="py-2 px-4 text-left">Состояние</th>
                <th className="py-2 px-4 text-left">Завершённых интервью</th>
              </tr>
            </thead>
            <tbody>
              {filteredSurveys.map((survey) => (
                <tr
                  key={survey.hash}
                  className={`border-b cursor-pointer ${getStateColor(survey.state)} hover:bg-gray-100`}
                  onClick={() => router.push(`/survey/${survey.hash}`)}
                >
                  <td className="py-2 px-4">{survey.title}</td>
                  <td className="py-2 px-4">
                    {new Date(survey.created_at).toLocaleDateString("ru-RU", {
                      timeZone: "UTC",
                    })}
                  </td>
                  <td className="py-2 px-4">
                    {new Date(survey.updated_at).toLocaleDateString("ru-RU", {
                      timeZone: "UTC",
                    })}
                  </td>
                  <td className="py-2 px-4">{SURVEY_STATE[survey.state]}</td>
                  <td
                    className={`py-2 px-4 ${
                      survey.completed_interviews > 0
                        ? "text-green-600 font-semibold"
                        : "text-gray-600"
                    }`}
                  >
                    {survey.completed_interviews}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default SurveysListPage;