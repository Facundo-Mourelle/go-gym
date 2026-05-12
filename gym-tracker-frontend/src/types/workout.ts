export type WorkoutPlanID = string;
export type WorkoutExerciseID = string;
export type ExerciseID = string;

export interface WorkoutPlan {
    ID: WorkoutPlanID;
    UserID: string;
    Name: string;
    Description: string;
    Exercises: WorkoutExercise[];
}

export interface WorkoutExercise {
    ID: WorkoutExerciseID;
    Order: number;
    ExerciseID: ExerciseID;
    Sets: number;
    Reps: number;
    RepsInReserve: number;
    Notes: string;
}
