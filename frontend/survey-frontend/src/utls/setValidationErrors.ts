import { FieldValues, UseFormSetError } from 'react-hook-form';

interface ValidationError {
  formField: string;
  message: string;
}

export const setValidationErrors = <T extends FieldValues>(
  errors: ValidationError[],
  setError: UseFormSetError<T>
) => {
  errors.forEach((err) => {
    if (err.formField) {
      setError(err.formField as keyof T, {
        type: 'manual',
        message: err.message,
      });
    } else {
      // Если поле formField отсутствует, устанавливаем общую ошибку для формы
      setError('general' as keyof T, {
        type: 'manual',
        message: err.message,
      });
    }
  });
};