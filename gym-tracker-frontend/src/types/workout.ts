export type WorkoutPlanID = string;
export type WorkoutExerciseID = string;
export type ExerciseID = string;
export type EquipmentType = string;

export interface WorkoutPlan {
    ID: WorkoutPlanID;
    Name: string;
    Description: string;
    CreatedAt: string;
    UpdatedAt: string;
    Exercises: WorkoutExercise[];
}

export interface WorkoutExercise {
    ID: WorkoutExerciseID;
    Order: number;
    ExerciseID: ExerciseID;
    ExerciseName: string;
    EquipmentType: EquipmentType;
    Sets: number;
    Reps: number;
    RepsInReserve: number;
    Notes: string;
}

export interface WorkoutSummary {
    ID: WorkoutPlanID;
    Name: string;
    Description: string;
    ExerciseCount: number;
    TotalSets: number;
    CreatedAt: string;
    UpdatedAt: string;
}

export interface CreateWorkoutData {
    Name: string;
    Description?: string;
    Exercises: CreateWorkoutExercise[];
}

export interface CreateWorkoutExercise {
    Order: number;
    ExerciseID: ExerciseID;
    Sets: number;
    Reps: number;
    RepsInReserve: number;
    Notes?: string;
}
