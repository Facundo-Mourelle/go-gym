import { apiClient } from './client';
import type { Routine } from '../types/routine';

export const routinesApi = {
    list: async (): Promise<Routine[]> => {
        const response = await apiClient.get<{ routines: Routine[]; count: number }>('/api/v1/routines');
        return response.data.routines;
    },

    get: async (id: string): Promise<Routine> => {
        const response = await apiClient.get<Routine>(`/api/v1/routines/${id}`);
        return response.data;
    },

    create: async (data: { name: string; description?: string; movement_patterns: string[] }): Promise<Routine> => {
        const response = await apiClient.post<Routine>('/api/v1/routines', data);
        return response.data;
    },

    delete: async (id: string): Promise<void> => {
        await apiClient.delete(`/api/v1/routines/${id}`);
    },
};