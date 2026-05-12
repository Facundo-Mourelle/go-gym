import { apiClient } from './client';
import type { Session, RecordSetRequest, RecordSetResponse } from '../types/session';

export interface StartSessionRequest {
    workout_plan_id?: string;
    notes?: string;
}

export interface StartSessionResponse {
    session_id: string;
    started_at: string;
    workout_plan_id?: string;
}

export const sessionsApi = {
    start: async (data: StartSessionRequest): Promise<StartSessionResponse> => {
        console.log('Starting session with data:', data);
        const response = await apiClient.post<StartSessionResponse>('/api/v1/sessions', data);
        console.log('Session start response:', response.data);
        return response.data;
    },

    recordSet: async (sessionId: string, data: RecordSetRequest): Promise<RecordSetResponse> => {
        const response = await apiClient.post<RecordSetResponse>(
            `/api/v1/sessions/${sessionId}/sets`,
            data
        );
        return response.data;
    },

    complete: async (sessionId: string, notes?: string): Promise<void> => {
        await apiClient.put(`/api/v1/sessions/${sessionId}/complete`, { notes });
    },

    get: async (sessionId: string): Promise<Session> => {
        const response = await apiClient.get<Session>(`/api/v1/sessions/${sessionId}`);
        return response.data;
    },

    list: async (limit = 50): Promise<Session[]> => {
        const response = await apiClient.get<Session[]>(`/api/v1/sessions?limit=${limit}`);
        return response.data;
    },
};
