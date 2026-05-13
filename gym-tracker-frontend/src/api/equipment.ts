import { apiClient } from './client';

export interface EquipmentData {
    id: string;
    name: string;
    type: 'freeweight' | 'machine' | 'cable';
    manufacturer?: string;
    user_id?: string;
    actual_weight?: number;
    pulley_type?: string;
    stack_weights?: number[];
    resistance_profile_id?: string;
    resistance_profile_name?: string;
}

export interface CreateEquipmentData {
    name: string;
    type: 'freeweight' | 'machine' | 'cable';
    manufacturer?: string;
    actual_weight?: number;
    pulley_type?: string;
    stack_weights?: number[];
    resistance_profile_id?: string;
    resistance_profile_name?: string;
}

export const equipmentApi = {
    list: async (): Promise<EquipmentData[]> => {
        const response = await apiClient.get<{ equipment: EquipmentData[]; count: number }>(
            '/api/v1/equipment'
        );
        return response.data.equipment;
    },

    get: async (id: string): Promise<EquipmentData> => {
        const response = await apiClient.get<EquipmentData>(`/api/v1/equipment/${id}`);
        return response.data;
    },

    create: async (data: CreateEquipmentData): Promise<EquipmentData> => {
        const response = await apiClient.post<EquipmentData>('/api/v1/equipment', data);
        return response.data;
    },

    update: async (id: string, data: Partial<CreateEquipmentData>): Promise<EquipmentData> => {
        const response = await apiClient.put<EquipmentData>(`/api/v1/equipment/${id}`, data);
        return response.data;
    },

    delete: async (id: string): Promise<void> => {
        await apiClient.delete(`/api/v1/equipment/${id}`);
    },
};
