import apiClient from './client';
import { GroundingRule } from '../types';

export const getGroundingRules = async (): Promise<GroundingRule[]> => {
  const response = await apiClient.get<GroundingRule[]>('/grounding');
  return response.data || [];
};

export const createGroundingRule = async (rule: Partial<GroundingRule>): Promise<GroundingRule> => {
  const response = await apiClient.post<GroundingRule>('/grounding', rule);
  return response.data;
};

export const updateGroundingRule = async (id: string, rule: Partial<GroundingRule>): Promise<GroundingRule> => {
  const response = await apiClient.put<GroundingRule>(`/grounding/${id}`, rule);
  return response.data;
};

export const deleteGroundingRule = async (id: string): Promise<void> => {
  await apiClient.delete(`/grounding/${id}`);
};
