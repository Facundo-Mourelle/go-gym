import { apiClient } from './client';

export interface ProgressDataPoint {
    Date: string;
    SessionID: string;
    Estimated1RM: number;
    Volume: number;
    MaxWeight: number;
    TotalSets: number;
    TotalReps: number;
}

export interface ProgressResponse {
    ExerciseID: string;
    ExerciseName: string;
    Data: ProgressDataPoint[];
    OverallTrend: string;
    ImprovementPercentage: number;
}

export const metricsApi = {
    getExerciseProgress: async (exerciseId: string, startDate?: string, endDate?: string): Promise<ProgressResponse> => {
        const params = new URLSearchParams();
        if (startDate) params.append('start_date', startDate);
        if (endDate) params.append('end_date', endDate);
        
        const response = await apiClient.get<ProgressResponse>(`/api/v1/metrics/progress/${exerciseId}?${params.toString()}`);
        return response.data;
    }
};
