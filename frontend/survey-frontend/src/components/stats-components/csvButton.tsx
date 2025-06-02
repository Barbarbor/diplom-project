import React from 'react';

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
    }>;
  };
}

export default function ExportToCsvButton({ data }: ExportToCsvButtonProps) {
  const exportToCsv = () => {
    const csvData: string[][] = [];

    // Общая статистика
    csvData.push(['Общая статистика']);
    csvData.push(['Начато интервью', data.started_interviews.toString()]);
    csvData.push(['Завершено интервью', data.completed_interviews.toString()]);
    csvData.push(['Среднее время прохождение анкеты', `${data.average_completion_time.toFixed(2)}s`]);
    csvData.push([]);

    // Статистика по вопросам
    data.questions.forEach((question) => {
      csvData.push([`Вопрос: ${question.label}`, `Тип: ${question.type}`]);

      if (question.type === 'single_choice') {
        const totalAnswers = question.answers.length;
        const answerCounts = question.answers.reduce((acc, ans) => {
          const numAns = Number(ans);
          acc[numAns] = (acc[numAns] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        question?.options?.forEach((option) => {
          const count = answerCounts[option.id] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([option.label, `${count} раз`, `${percentage.toFixed(1)}%`]);
        });
      } else if (question.type === 'multi_choice') {
        const totalInterviews = question.answers.length;
        const answerCounts = question.answers.reduce((acc, ans) => {
          const parsedAns = JSON.parse(ans) as number[];
          parsedAns.forEach((id) => {
            acc[id] = (acc[id] || 0) + 1;
          });
          return acc;
        }, {} as Record<number, number>);

        question?.options?.forEach((option) => {
          const count = answerCounts[option.id] || 0;
          const percentage = totalInterviews > 0 ? (count / totalInterviews) * 100 : 0;
          csvData.push([option.label, `${count} раз`, `${percentage.toFixed(1)}%`]);
        });
      } else if (question.type === 'consent') {
        const totalAnswers = question.answers.length;
        const trueCount = question.answers.filter((ans) => ans === 'true').length;
        const falseCount = totalAnswers - trueCount;
        const truePercentage = totalAnswers > 0 ? (trueCount / totalAnswers) * 100 : 0;
        const falsePercentage = totalAnswers > 0 ? (falseCount / totalAnswers) * 100 : 0;
        csvData.push(['Согласны', `${trueCount} раз`, `${truePercentage.toFixed(1)}%`]);
        csvData.push(['Не согласны', `${falseCount} раз`, `${falsePercentage.toFixed(1)}%`]);
      } else if (question.type === 'rating') {
        const totalAnswers = question.answers.length;
        const sum = question.answers.reduce((acc, ans) => acc + Number(ans), 0);
        const average = totalAnswers > 0 ? (sum / totalAnswers).toFixed(1) : '0.0';
        csvData.push(['Средний рейтинг', average]);

        const frequency = question.answers.reduce((acc, ans) => {
          const numAns = Number(ans);
          acc[numAns] = (acc[numAns] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        for (let star = 1; star <= 5; star++) {
          const count = frequency[star] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([`${star} звезда`, `${count} раз`, `${percentage.toFixed(1)}%`]);
        }
      } else if (question.type === 'short_text' || question.type === 'long_text') {
        const answerCounts = question.answers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        const sortedAnswers = Object.entries(answerCounts).sort((a, b) => b[1] - a[1]);
        const top5 = sortedAnswers.slice(0, 5);
        top5.forEach(([answer, count]) => {
          csvData.push([answer, `${count} раз`]);
        });
      } else if (question.type === 'date') {
        const dateCounts = question.answers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        Object.entries(dateCounts).forEach(([date, count]) => {
          csvData.push([date, `${count} раз`]);
        });
      } else if (question.type === 'number') {
        const numbers = question.answers.map(Number).filter((num) => !isNaN(num));
        const frequencyMap = numbers.reduce((acc, num) => {
          acc[num] = (acc[num] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        Object.entries(frequencyMap).forEach(([num, count]) => {
          csvData.push([num, `${count} раз`]);
        });
      } else if (question.type === 'email') {
        const uniqueEmails = [...new Set(question.answers)].sort();
        uniqueEmails.forEach((email) => {
          csvData.push([email]);
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