export type SessionID = string;
export type PerformedSetID = string;
export type WorkoutPlanID = string;
export type EquipmentID = string;

export interface Session {
    session_id: SessionID;
    workout_plan_id?: WorkoutPlanID;
    started_at: string;
    completed_at?: string;
    duration?: number; // in seconds
    notes: string;
    exercise_groups: ExerciseGroup[];
    total_sets: number;
    total_volume: number;
}

export interface ExerciseGroup {
    exercise_id: string;
    exercise_name: string;
    sets: PerformedSet[];
    total_volume: number;
    set_count: number;
}

export interface PerformedSet {
    set_id: PerformedSetID;
    set_number: number;
    reps: number;
    raw_load: number;
    effective_load: number;
    volume: number;
    equipment_id: EquipmentID;
    reps_in_reserve?: number;
    notes: string;
    performed_at: string;
}

export interface RecordSetRequest {
    exercise_id: string;
    set_number: number;
    reps: number;
    reps_in_reserve: number;
    raw_load: number;
    equipment_id: EquipmentID;
    notes?: string;
}

export interface RecordSetResponse {
    performed_set_id: PerformedSetID;
    effective_load: number;
    volume: number;
    performed_at: string;
}
