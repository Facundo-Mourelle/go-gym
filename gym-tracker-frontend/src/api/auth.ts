import { apiClient } from './client';
import type { AuthResponse, LoginRequest, RegisterRequest, User } from '../types/auth';

export const authApi = {
    register: async (data: RegisterRequest): Promise<AuthResponse> => {
        const response = await apiClient.post<AuthResponse>('/api/v1/auth/register', data);
        return response.data;
    },

    login: async (data: LoginRequest): Promise<AuthResponse> => {
        const response = await apiClient.post<AuthResponse>('/api/v1/auth/login', data);
        return response.data;
    },

    getCurrentUser: async (): Promise<User> => {
        const response = await apiClient.get<User>('/api/v1/auth/me');
        return response.data;
    },
};
