export interface Routine {
    id: string;
    user_id: string;
    name: string;
    description: string;
    movement_patterns: string[];
    is_preset: boolean;
    created_at: string;
    updated_at: string;
}