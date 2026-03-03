import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getMemories, searchMemories, giveMemoryFeedback } from '../api/memories';
import { MemoryFeedback } from '../types';

export const memoryKeys = {
  all: ['memories'] as const,
  list: () => [...memoryKeys.all, 'list'] as const,
  search: (query: string) => [...memoryKeys.all, 'search', query] as const,
};

export function useMemories() {
  return useQuery({
    queryKey: memoryKeys.list(),
    queryFn: getMemories,
  });
}

export function useSearchMemories(query: string) {
  return useQuery({
    queryKey: memoryKeys.search(query),
    queryFn: () => searchMemories(query),
    enabled: query.length > 2,
  });
}

export function useGiveMemoryFeedback() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, feedback }: { id: string; feedback: MemoryFeedback }) => 
      giveMemoryFeedback(id, feedback),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: memoryKeys.all });
    },
  });
}
