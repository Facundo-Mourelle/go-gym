import { apiClient } from './client';
import type { WorkoutPlan, WorkoutSummary } from '../types/workout';

interface ListWorkoutsResponse {
    workouts: WorkoutSummary[];
    count: number;
}

export const workoutsApi = {
    list: async (): Promise<WorkoutSummary[]> => {
        const response = await apiClient.get<ListWorkoutsResponse>('/api/v1/workouts');
        return response.data.workouts;
    },

    get: async (id: string): Promise<WorkoutPlan> => {
        const response = await apiClient.get<WorkoutPlan>(`/api/v1/workouts/${id}`);
        return response.data;
    },

    create: async (data: Record<string, unknown>): Promise<WorkoutPlan> => {
        const response = await apiClient.post<WorkoutPlan>('/api/v1/workouts', data);
        return response.data;
    },

    update: async (id: string, data: Record<string, unknown>): Promise<WorkoutPlan> => {
        const response = await apiClient.put<WorkoutPlan>(`/api/v1/workouts/${id}`, data);
        return response.data;
    },

    delete: async (id: string): Promise<void> => {
        await apiClient.delete(`/api/v1/workouts/${id}`);
    },
};
