import { useEffect, useState } from 'react';
import { useStartInterview } from '@/hooks/react-query/interview';

export const useGetInterviewId = (hash: string, isDemo?: string) => {
  const [interviewId, setInterviewId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { mutate: startInterview, isLoading: startLoading, error: startError } = useStartInterview();

  useEffect(() => {
    const generateRandomId = () => {
      const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
      let result = '';
      for (let i = 0; i < 16; i++) {
        result += characters.charAt(Math.floor(Math.random() * characters.length));
      }
      return result;
    };

    const getOrCreateInterviewId = async (isDemo?: boolean) => {
      const storedId = localStorage.getItem(`survey-${hash}`);
      if (storedId && !isDemo) {
        setInterviewId(storedId);
        setLoading(false);
        return;
      }

      const newId = generateRandomId();
      try {
        await startInterview(
          { hash, interviewId: newId, isDemo },
          {
            onSuccess: () => {
            
              if(!isDemo) localStorage.setItem(`survey-${hash}`, newId);
              setInterviewId(newId);
              setLoading(false);
            },
            onError: () => {
              setError('Ошибка при старте интервью');
              setLoading(false);
            },
          }
        );
      } catch (err) {
        setError('Не удалось начать интервью');
        setLoading(false);
      }
    };

    getOrCreateInterviewId(isDemo);
  }, [hash, startInterview]);

  return { interviewId, loading: loading || startLoading, error: error || startError?.message };
};