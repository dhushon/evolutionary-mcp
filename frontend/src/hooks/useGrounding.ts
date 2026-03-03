import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  getGroundingRules, 
  createGroundingRule, 
  updateGroundingRule, 
  deleteGroundingRule 
} from '../api/grounding';
import { GroundingRule } from '../types';

export const groundingKeys = {
  all: ['grounding'] as const,
  list: () => [...groundingKeys.all, 'list'] as const,
  detail: (id: string) => [...groundingKeys.all, 'detail', id] as const,
};

export function useGroundingRules() {
  return useQuery({
    queryKey: groundingKeys.list(),
    queryFn: getGroundingRules,
  });
}

export function useCreateGroundingRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (rule: Partial<GroundingRule>) => createGroundingRule(rule),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: groundingKeys.list() });
    },
  });
}

export function useUpdateGroundingRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, rule }: { id: string; rule: Partial<GroundingRule> }) => 
      updateGroundingRule(id, rule),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: groundingKeys.list() });
      queryClient.invalidateQueries({ queryKey: groundingKeys.detail(variables.id) });
    },
  });
}

export function useDeleteGroundingRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteGroundingRule(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: groundingKeys.list() });
    },
  });
}
