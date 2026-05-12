import { apiClient } from './client';
import type { Exercise } from '../types/exercise';

export interface PatternInfo {
    pattern: string;
    name: string;
    plane: string;
    primary_muscles: string[];
    joint_actions: string[];
    description: string;
}

export interface CreateExerciseData {
    name: string;
    description?: string;
    primary_patterns: { pattern: string; contribution: number; range_of_motion: string; notes: string }[];
    equipment: string[];
}

export const exercisesApi = {
    list: async (filters?: { search?: string; pattern?: string; muscle?: string; equipment?: string }): Promise<Exercise[]> => {
        const params = new URLSearchParams();
        if (filters?.search) params.append('search', filters.search);
        if (filters?.pattern) params.append('pattern', filters.pattern);
        if (filters?.muscle) params.append('muscle', filters.muscle);
        if (filters?.equipment) params.append('equipment', filters.equipment);

        const response = await apiClient.get<{ exercises: Exercise[]; count: number }>(
            `/api/v1/exercises?${params.toString()}`
        );
        return response.data.exercises;
    },

    get: async (id: string): Promise<Exercise> => {
        const response = await apiClient.get<Exercise>(`/api/v1/exercises/${id}`);
        return response.data;
    },

    create: async (data: CreateExerciseData): Promise<Exercise> => {
        const response = await apiClient.post<Exercise>('/api/v1/exercises/custom', data);
        return response.data;
    },

    update: async (id: string, data: Partial<CreateExerciseData>): Promise<Exercise> => {
        const response = await apiClient.put<Exercise>(`/api/v1/exercises/${id}`, data);
        return response.data;
    },

    listPatterns: async (): Promise<PatternInfo[]> => {
        const response = await apiClient.get<{ patterns: PatternInfo[]; count: number }>('/api/v1/patterns');
        return response.data.patterns;
    },
};
