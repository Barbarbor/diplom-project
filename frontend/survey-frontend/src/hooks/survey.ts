import { useParams } from 'next/navigation';

export const useSurveyHash = () => {
  const params = useParams();
  const hash = params.hash as string;
  if (!hash) {
    throw new Error('useSurveyHash must be used within a route with a hash parameter');
  }
  return hash;
};