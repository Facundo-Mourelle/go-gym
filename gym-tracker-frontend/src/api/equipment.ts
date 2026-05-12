import { apiClient } from './client';

export interface EquipmentItem {
    id: string;
    name: string;
    type: string;
}

export const equipmentApi = {
    list: async (): Promise<EquipmentItem[]> => {
        const response = await apiClient.get<{ equipment: EquipmentItem[]; count: number }>(
            '/api/v1/equipment'
        );
        return response.data.equipment;
    },
};
