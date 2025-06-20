import React from 'react';
import { useTranslation } from 'react-i18next';

interface ExportToCsvButtonProps {
  data: {
    started_interviews: number;
    completed_interviews: number;
    average_completion_time: number;
    questions: Array<{
      id: string;
      label: string;
      type: string;
      answers: string[];
      options?: Array<{ id: number; label: string }>;
      extra_params: {
        starsCount?: number; // Для rating
      };
    }>;
  };
}

// Функция фильтрации ответов
function filterValidAnswers(question: ExportToCsvButtonProps['data']['questions'][0]): string[] {
  const { type, answers, options, extra_params } = question;
  const nonNullableAnswers = answers === null ? [] : answers;

  switch (type) {
    case 'single_choice':
      const validOptionIds = options?.map(opt => opt.id.toString()) || [];
      return nonNullableAnswers.filter(ans => validOptionIds.includes(ans));

    case 'multi_choice':
      const validOptionIdsSet = new Set(options?.map(opt => opt.id.toString()) || []);
      return nonNullableAnswers
        .map(ans => {
          try {
            const parsedAns = JSON.parse(ans) as string[];
            const filteredAns = parsedAns.filter(id => validOptionIdsSet.has(id));
            return filteredAns.length > 0 ? JSON.stringify(filteredAns) : null;
          } catch {
            return null;
          }
        })
        .filter(ans => ans !== null) as string[];

    case 'rating':
      const starsCount = extra_params?.starsCount || 5; // Динамический starsCount
      return nonNullableAnswers.filter(ans => {
        const num = Number(ans);
        return !isNaN(num) && num >= 1 && num <= starsCount;
      });

    case 'number':
      return nonNullableAnswers.filter(ans => !isNaN(Number(ans)));

    case 'date':
      return nonNullableAnswers.filter(ans => !isNaN(new Date(ans).getTime()));

    case 'email':
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      return nonNullableAnswers.filter(ans => emailRegex.test(ans));

    case 'consent':
      return nonNullableAnswers.filter(ans => ans === 'true' || ans === 'false');

    case 'short_text':
    case 'long_text':
      return nonNullableAnswers; // Все текстовые ответы валидны

    default:
      return nonNullableAnswers;
  }
}

export default function ExportToCsvButton({ data }: ExportToCsvButtonProps) {
  const { t } = useTranslation();

  const exportToCsv = () => {
    const csvData: string[][] = [];

    // Общая статистика
    csvData.push([t('survey.stats.csv.generalStatistics')]);
    csvData.push([t('survey.stats.csv.startedInterviews'), data.started_interviews.toString()]);
    csvData.push([t('survey.stats.csv.completedInterviews'), data.completed_interviews.toString()]);
    csvData.push([t('survey.stats.csv.averageCompletionTime'), `${data.average_completion_time.toFixed(2)}s`]);
    csvData.push([]);

    // Статистика по вопросам
    data.questions.forEach((question) => {
      const validAnswers = filterValidAnswers(question); // Фильтруем ответы
    csvData.push([
        t('survey.stats.csv.question') + `: ${question.label}`,
        t('survey.stats.csv.type') + `: ${t(`survey.question.type.${question.type}`)}`
      ]);

      if (question.type === 'single_choice') {
        const totalAnswers = validAnswers.length;
        const answerCounts = validAnswers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        question?.options?.forEach((option) => {
          const count = answerCounts[option.id.toString()] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([t('survey.stats.csv.singleChoiceOption', { label: option.label, count, percentage: percentage.toFixed(1) })]);
        });
      } else if (question.type === 'multi_choice') {
        const totalInterviews = validAnswers.length;
        const answerCounts = validAnswers.reduce((acc, ans) => {
          try {
            const parsedAns = JSON.parse(ans) as string[];
            parsedAns.forEach((id) => {
              acc[id] = (acc[id] || 0) + 1;
            });
          } catch {}
          return acc;
        }, {} as Record<string, number>);

        question?.options?.forEach((option) => {
          const count = answerCounts[option.id.toString()] || 0;
          const percentage = totalInterviews > 0 ? (count / totalInterviews) * 100 : 0;
          csvData.push([t('survey.stats.csv.multiChoiceOption', { label: option.label, count, percentage: percentage.toFixed(1) })]);
        });
      } else if (question.type === 'consent') {
        const totalAnswers = validAnswers.length;
        const trueCount = validAnswers.filter((ans) => ans === 'true').length;
        const falseCount = totalAnswers - trueCount;
        const truePercentage = totalAnswers > 0 ? (trueCount / totalAnswers) * 100 : 0;
        const falsePercentage = totalAnswers > 0 ? (falseCount / totalAnswers) * 100 : 0;
        csvData.push([t('survey.stats.csv.consentYes', { count: trueCount, percentage: truePercentage.toFixed(1) })]);
        csvData.push([t('survey.stats.csv.consentNo', { count: falseCount, percentage: falsePercentage.toFixed(1) })]);
      } else if (question.type === 'rating') {
        const totalAnswers = validAnswers.length;
        const sum = validAnswers.reduce((acc, ans) => acc + Number(ans), 0);
        const average = totalAnswers > 0 ? (sum / totalAnswers).toFixed(1) : '0.0';
        csvData.push([t('survey.stats.csv.ratingAverage'), average]);

        const starsCount = question.extra_params?.starsCount || 5; // Динамический starsCount
        const frequency = validAnswers.reduce((acc, ans) => {
          const numAns = Number(ans);
          acc[numAns] = (acc[numAns] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        for (let star = 1; star <= starsCount; star++) {
          const count = frequency[star] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([t('survey.stats.csv.ratingStar', { star, count, percentage: percentage.toFixed(1) })]);
        }
      } else if (question.type === 'short_text' || question.type === 'long_text') {
        const answerCounts = validAnswers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        const sortedAnswers = Object.entries(answerCounts).sort((a, b) => b[1] - a[1]);
        const top5 = sortedAnswers.slice(0, 5);
        top5.forEach(([answer, count]) => {
          csvData.push([t('survey.stats.csv.textAnswer', { answer, count })]);
        });
      } else if (question.type === 'date') {
        const dateCounts = validAnswers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        Object.entries(dateCounts).forEach(([date, count]) => {
          csvData.push([t('survey.stats.csv.dateAnswer', { date, count })]);
        });
      } else if (question.type === 'number') {
        const numbers = validAnswers.map(Number).filter((num) => !isNaN(num));
        const frequencyMap = numbers.reduce((acc, num) => {
          acc[num] = (acc[num] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        Object.entries(frequencyMap).forEach(([num, count]) => {
          csvData.push([t('survey.stats.csv.numberAnswer', { number: num, count })]);
        });
      } else if (question.type === 'email') {
        const uniqueEmails = [...new Set(validAnswers)].sort();
        uniqueEmails.forEach((email) => {
          csvData.push([t('survey.stats.csv.emailAnswer', { email })]);
        });
      }

      csvData.push([]); // Пустая строка между вопросами
    });

    // Генерация CSV строки
    const csv = csvData.map((row) => row.join(';')).join('\n');

    // Добавляем BOM для UTF-8
    const BOM = '\uFEFF'; // UTF-8 BOM
    const csvWithBom = BOM + csv;

    // Создание и скачивание файла
    const blob = new Blob([csvWithBom], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `survey_stats_${new Date().toISOString()}.csv`;
    link.click();
  };

  return (
    <button
      className="px-4 py-2 bg-blue-600 text-white rounded"
      onClick={exportToCsv}
    >
      Экспорт в CSV
    </button>
  );
}