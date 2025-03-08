
import { getSurvey } from "@/api-client/survey";

interface SurveyPageProps {
  params: { hash: string };
}

export default async function SurveyPage({ params }: SurveyPageProps) {
  const { hash } = await params;
  const response = await getSurvey(hash);

  if (response.status >= 400 || !response.data?.survey) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-red-500 text-lg">{response.error}</p>
      </div>
    );
  }

  const { survey } = response.data; // предполагаем, что getSurvey возвращает { survey, email }

  return (
    <div className="max-w-3xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-4">{survey.title}</h1>
      <div className="mb-4">
        <p className="text-gray-700">
          <span className="font-semibold">Created at:</span>{" "}
          {new Date(survey.created_at).toLocaleString()}
        </p>
        <p className="text-gray-700">
          <span className="font-semibold">Updated at:</span>{" "}
          {new Date(survey.updated_at).toLocaleString()}
        </p>
        <p className="text-gray-700">
          <span className="font-semibold">Creator:</span> {survey.creator}
        </p>
      </div>
      {/* Здесь можно добавить дополнительное отображение данных опроса, например вопросы */}
    </div>
  );
}
