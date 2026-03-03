import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getWorkflows, putWorkflow, deleteWorkflow, getTenant, getWorkflow } from '../api/workflows';
import { WorkflowUpdatePayload } from '../types';

/**
 * Key for caching workflow and tenant data
 */
export const workflowKeys = {
  all: ['workflows'] as const,
  list: () => [...workflowKeys.all, 'list'] as const,
  details: (id: string) => [...workflowKeys.all, 'detail', id] as const,
};

export const tenantKeys = {
  all: ['tenant'] as const,
  current: () => [...tenantKeys.all, 'current'] as const,
};

/**
 * Hook for fetching current tenant branding.
 */
export function useTenant() {
  return useQuery({
    queryKey: tenantKeys.current(),
    queryFn: getTenant,
    staleTime: 1000 * 60 * 30, // 30 minutes, branding rarely changes
  });
}

/**
 * Hook for fetching all workflows.
 */
export function useWorkflows() {
  return useQuery({
    queryKey: workflowKeys.list(),
    queryFn: getWorkflows,
  });
}

/**
 * Hook for fetching a single workflow by ID.
 */
export function useWorkflow(id: string | null) {
  return useQuery({
    queryKey: workflowKeys.details(id || ''),
    queryFn: () => getWorkflow(id!),
    enabled: !!id,
  });
}

/**
 * Hook for creating or updating a workflow.
 */
export function usePutWorkflow() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (workflow: WorkflowUpdatePayload) => putWorkflow(workflow),
    onSuccess: () => {
      // Invalidate the list to trigger a refetch
      queryClient.invalidateQueries({ queryKey: workflowKeys.list() });
    },
  });
}

/**
 * Hook for deleting a workflow.
 */
export function useDeleteWorkflow() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteWorkflow(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workflowKeys.list() });
    },
  });
}
