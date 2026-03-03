import apiClient from './client';
import { Workflow, HealthStatus, Tenant, WorkflowUpdatePayload } from '../types';

/**
 * Retrieves the current tenant's branding configuration.
 */
export const getTenant = async (): Promise<Tenant> => {
  const response = await apiClient.get<Tenant>('/tenant');
  return response.data;
};

/**
 * Retrieves all workflows for the current tenant.
 */
export const getWorkflows = async (): Promise<Workflow[]> => {
  const response = await apiClient.get<Workflow[]>('/workflows');
  return response.data || [];
};

/**
 * Retrieves a specific workflow by ID.
 */
export const getWorkflow = async (id: string): Promise<Workflow> => {
  const response = await apiClient.get<Workflow>(`/workflows/${id}`);
  return response.data;
};

/**
 * Creates or updates a workflow (supports versioning via save_as_new_version flag).
 */
export const putWorkflow = async (workflow: WorkflowUpdatePayload): Promise<Workflow> => {
  const response = await apiClient.put<Workflow>('/workflows', workflow);
  return response.data;
};

/**
 * Partially updates an existing workflow.
 */
export const patchWorkflow = async (workflow: Partial<Workflow>): Promise<Workflow> => {
  const response = await apiClient.patch<Workflow>('/workflows', workflow);
  return response.data;
};

/**
 * Deletes a specific workflow version.
 */
export const deleteWorkflow = async (id: string): Promise<void> => {
  await apiClient.delete(`/workflows?id=${id}`);
};

/**
 * Public health check (bypass auth)
 */
export const getHealth = async (): Promise<HealthStatus> => {
  // Use absolute URL to bypass the /api/v1 prefix if necessary, 
  // or relative if standard.
  const response = await apiClient.get<HealthStatus>('/health');
  return response.data;
};
