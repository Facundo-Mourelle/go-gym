import { create } from 'zustand';
import type { Session, PerformedSet } from '../types/session';

interface SessionState {
    // Current active session
    activeSession: Session | null;
    currentExerciseId: string | null;

    // UI state
    isRecording: boolean;

    // Actions
    setActiveSession: (session: Session) => void;
    setCurrentExercise: (exerciseId: string) => void;
    addPerformedSet: (exerciseId: string, set: PerformedSet) => void;
    prepopulateExercises: (exercises: { id: string; name: string }[]) => void;
    clearActiveSession: () => void;
}

export const useSessionStore = create<SessionState>((set, get) => ({
    activeSession: null,
    currentExerciseId: null,
    isRecording: false,

    setActiveSession: (session) => set({ activeSession: session }),

    setCurrentExercise: (exerciseId) => set({ currentExerciseId: exerciseId }),

    addPerformedSet: (exerciseId, performedSet) => {
        const { activeSession } = get();
        if (!activeSession) return;

        // Find or create exercise group
        const exerciseGroupIndex = activeSession.exercise_groups.findIndex(
            (g) => g.exercise_id === exerciseId
        );

        const updatedGroups = [...activeSession.exercise_groups];

        const setVolume = performedSet.volume ?? 0;

        if (exerciseGroupIndex >= 0) {
            // Add to existing group
            updatedGroups[exerciseGroupIndex].sets.push(performedSet);
            updatedGroups[exerciseGroupIndex].set_count++;
            updatedGroups[exerciseGroupIndex].total_volume += setVolume;
        } else {
            // Create new group
            updatedGroups.push({
                exercise_id: exerciseId,
                exercise_name: '', // Will be populated from exercise data
                sets: [performedSet],
                total_volume: setVolume,
                set_count: 1,
            });
        }

        set({
            activeSession: {
                ...activeSession,
                exercise_groups: updatedGroups,
                total_sets: activeSession.total_sets + 1,
                total_volume: activeSession.total_volume + setVolume,
            },
        });
    },

    prepopulateExercises: (exercises) => {
        const { activeSession, currentExerciseId } = get();
        if (!activeSession) return;

        const existingIds = new Set(activeSession.exercise_groups.map(g => g.exercise_id));
        const newGroups = exercises
            .filter(ex => !existingIds.has(ex.id))
            .map(ex => ({
                exercise_id: ex.id,
                exercise_name: ex.name,
                sets: [],
                total_volume: 0,
                set_count: 0,
            }));

        if (newGroups.length === 0) return;

        set({
            activeSession: {
                ...activeSession,
                exercise_groups: [...activeSession.exercise_groups, ...newGroups],
            },
            currentExerciseId: currentExerciseId || newGroups[0].exercise_id,
        });
    },

    clearActiveSession: () =>
        set({
            activeSession: null,
            currentExerciseId: null,
            isRecording: false,
        })
}));
