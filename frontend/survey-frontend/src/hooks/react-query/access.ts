import {useQuery, useMutation, UseQueryOptions, UseMutationOptions, useQueryClient} from '@tanstack/react-query';
import {getAccessList, addEditAccess, removeEditAccess} from '../../api-client/access';
import {ApiResponse} from '@/lib/api';

// Assuming queryClient is available globally or imported
import {QueryClient} from '@tanstack/react-query';
import {useSurveyHash} from '../survey';

export const useAccessList = (
    hash: string,
    options?: UseQueryOptions<ApiResponse<string[]>, Error>
) => {
    return useQuery<ApiResponse<string[]>, Error>({
        queryKey: ['accessList', hash],
        queryFn: () => getAccessList(hash),
        enabled: !!hash, // Only run query if hash is provided
        select: data => data.data || [], // Extract emails from ApiResponse
        ...options
    });
};

export const useAddEditAccess = (
    options?: UseMutationOptions<ApiResponse<void>, Error, {hash: string; email: string}>
) => {
    const queryClient = useQueryClient();

    const hash = useSurveyHash();
    return useMutation<ApiResponse<void>, Error, {hash: string; email: string}>({
        mutationFn: ({hash, email}) => addEditAccess(hash, email),
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ['accessList', hash]});
        },
        ...options
    });
};

export const useRemoveEditAccess = (
    options?: UseMutationOptions<ApiResponse<void>, Error, {hash: string; email: string}>
) => {
    const queryClient = useQueryClient();

    const hash = useSurveyHash();
    return useMutation<ApiResponse<void>, Error, {hash: string; email: string}>({
        mutationFn: ({hash, email}) => removeEditAccess(hash, email),
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ['accessList', hash]});
        },
        ...options
    });
};
