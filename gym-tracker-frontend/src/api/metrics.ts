import { apiClient } from './client';

export interface ProgressDataPoint {
    date: string;
    session_id: string;
    set_id: string;
    set_number: number;
    reps: number;
    score: number;
    effective_load: number;
    weight: number;
    rir: number;
    equipment_type: string;
}

export interface SetNumberSummary {
    set_number: number;
    data_point_count: number;
    best_score: number;
    worst_score: number;
    average_score: number;
    latest_score: number;
    first_score: number;
    trend_direction: string;
    trend_strength: number;
    slope: number;
    total_improvement: number;
    percent_improvement: number;
}

export interface ProgressSummaryResponse {
    exercise_id: string;
    exercise_name: string;
    date_range: { start_date: string; end_date: string };
    set_summaries: Record<number, SetNumberSummary>;
    overall_best_score: number;
    overall_average_score: number;
    overall_latest_score: number;
    total_data_points: number;
}

export interface ProgressResponse {
    exercise_id: string;
    exercise_name: string;
    all_data_points: ProgressDataPoint[];
    data_by_set_number: Record<number, ProgressDataPoint[]>;
    summary: ProgressSummaryResponse;
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
