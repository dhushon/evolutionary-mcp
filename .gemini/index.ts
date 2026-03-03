import axios from 'axios';
import { Workflow } from '../types';

const apiClient = axios.create({
    baseURL: '/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

export const getWorkflows = async (): Promise<Workflow[]> => {
    const response = await apiClient.get('/workflows');
    // The backend now correctly returns [] for empty lists.
    return response.data;
};