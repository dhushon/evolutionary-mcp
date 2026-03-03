import apiClient from './client';
import { Memory, MemoryFeedback } from '../types';

export const getMemories = async (): Promise<Memory[]> => {
  const response = await apiClient.get<Memory[]>('/memories');
  return response.data || [];
};

export const searchMemories = async (query: string): Promise<Memory[]> => {
  const response = await apiClient.post<Memory[]>('/memories/search', { query });
  return response.data || [];
};

export const giveMemoryFeedback = async (id: string, feedback: MemoryFeedback): Promise<void> => {
  await apiClient.post(`/memories/${id}/feedback`, feedback);
};
