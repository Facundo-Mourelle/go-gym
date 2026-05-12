export type ExerciseID = string;
export type EquipmentID = string;
export type EquipmentType = 'dumbbell' | 'barbell' | 'machine';

export interface Exercise {
    id: ExerciseID;
    name: string;
    description: string;
    primary_patterns: PatternContribution[];
    secondary_patterns: PatternContribution[];
    derived_muscles: {
        primary: string[];
        secondary: string[];
    };
    equipment: EquipmentType[];
    source: 'system' | 'user';
    is_custom: boolean;
}

export interface PatternContribution {
    pattern: string;
    pattern_name: string;
    contribution: number;
    range_of_motion: string;
    notes: string;
}

export interface Equipment {
    id: EquipmentID;
    type: EquipmentType;
    name: string;
    manufacturer?: string;
    model_number?: string;
}
