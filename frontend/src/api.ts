import axios from 'axios';

export interface Workflow {
  id: string;
  workflow_id: string;
  tenant_id: string;
  name: string;
  description: string;
  status: "draft" | "active" | "archived";
  version: number;
  is_latest: boolean;
  input_schema: Record<string, any>;
  output_schema?: Record<string, any>;
  created_by: string;
  created_at: string;
  updated_at: string;
}

const apiClient = axios.create({
  baseURL: '/api/v1', // Use relative path to leverage Vite proxy
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getWorkflows = async (): Promise<Workflow[]> => {
  const response = await apiClient.get<Workflow[]>('/workflows');
  return response.data;
};

export const createWorkflow = async (workflow: Partial<Workflow>): Promise<Workflow> => {
  const response = await apiClient.put<Workflow>('/workflows', workflow);
  return response.data;
};

export const getHealth = async () => {
  const response = await apiClient.get('/health');
  return response.data;
};
