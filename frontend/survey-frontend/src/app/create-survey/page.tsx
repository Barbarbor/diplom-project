'use client'
import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';  // Путь к файлу, где настроен axios или fetch

// Интерфейс для Survey (или модель)
interface Survey {
  id?: number;
  title: string;
  description: string;
  content: string;  // Добавляем поле для контента
}

export default function CreateSurvey() {
  const queryClient = useQueryClient();

  // Мутация для создания нового опроса
  const mutation = useMutation({
    mutationFn: async (survey: Survey) => {
      const res = await api.post('/surveys', survey);
      return res.data;
    },
    onSuccess: () => {
      // После успешного создания опроса, обновляем кэш с опросами
      queryClient.invalidateQueries({ queryKey: ['surveys'] });
    },
  });

  // Обработчик отправки формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    const newSurvey: Survey = {
      title: formData.get('title') as string,
      description: formData.get('description') as string,
      content: formData.get('content') as string,  // Получаем значение поля content
    };

    mutation.mutate(newSurvey);
  };

  return (
    <div>
      <h1>Create Survey</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="title">Title</label>
          <input type="text" id="title" name="title" required />
        </div>
        <div>
          <label htmlFor="description">Description</label>
          <input type="text" id="description" name="description" required />
        </div>
        <div>
          <label htmlFor="content">Content</label>
          <textarea id="content" name="content" required />  {/* Добавляем поле для контента */}
        </div>
        <button type="submit" disabled={mutation.isLoading}>
          {mutation.isLoading ? 'Creating...' : 'Create Survey'}
        </button>
      </form>

      {mutation.isError && (
        <p>Error creating survey: {mutation.error instanceof Error ? mutation.error.message : JSON.stringify(mutation.error)}</p>
      )}
      {mutation.isSuccess && <p>Survey created successfully!</p>}
    </div>
  );
}
