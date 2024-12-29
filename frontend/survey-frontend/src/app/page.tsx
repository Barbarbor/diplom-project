'use client'
import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';

export default function Home() {
  const { data, isLoading } = useQuery({
      queryKey: ['surveys'],
      queryFn: async () => {
          const res = await api.get('/surveys');
          console.log(res.data);  // Проверьте структуру данных
          return res.data;
      }
  });

  if (isLoading) return <div>Loading...</div>;

  return (
      <div>
          {data.map((survey: any) => (
              <div key={survey.id}>
                  <h2>{survey.title}</h2>
                  <p>{survey.description}</p>
              </div>
          ))}
      </div>
  );
}