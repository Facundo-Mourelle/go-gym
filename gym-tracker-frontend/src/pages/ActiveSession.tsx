import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useSessionStore } from '../store/sessionStore';
import { sessionsApi } from '../api/sessions';
import { routinesApi } from '../api/routines';
import { workoutsApi } from '../api/workouts';
import { exercisesApi } from '../api/exercises';
import { SetRecorder } from '../components/session/SetRecorder';
import { CompletedSets } from '../components/session/CompletedSets';
import { ExerciseSelector } from '../components/session/ExerciseSelector';
import type { Exercise } from '../types/exercise';
import type { RecordSetRequest } from '../types/session';
import { StopCircle, Plus } from 'lucide-react';

export const ActiveSession: React.FC = () => {
    const navigate = useNavigate();
    const {
        activeSession,
        currentExerciseId,
        setActiveSession,
        setCurrentExercise,
        addPerformedSet,
        prepopulateExercises,
        clearActiveSession,
    } = useSessionStore();

    const [showExerciseSelector, setShowExerciseSelector] = useState(false);
    const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(null);
    const [isRecording, setIsRecording] = useState(false);
    const [isCompleting, setIsCompleting] = useState(false);
    const [searchParams] = useSearchParams();
    const [routinePatterns, setRoutinePatterns] = useState<string[]>([]);
    const [lastSessionMap, setLastSessionMap] = useState<Map<string, { reps: number; weight: number }>>(new Map());
    const routineSetupDone = useRef(false);
    const workoutSetupDone = useRef(false);

    const startNewSession = React.useCallback(async () => {
        try {
            const response = await sessionsApi.start({});
            setActiveSession({
                session_id: response.session_id,
                started_at: response.started_at,
                workout_plan_id: response.workout_plan_id,
                exercise_groups: [],
                total_sets: 0,
                total_volume: 0,
                notes: '',
            });
        } catch (error) {
            console.error('Failed to start session:', error);
            alert('Failed to start session');
        }
    }, [setActiveSession]);

    // Start session on mount if none exists
    useEffect(() => {
        if (!activeSession) {
            startNewSession();
        }
    }, [activeSession, startNewSession]);

    useEffect(() => {
        const routineId = searchParams.get('routineId');
        if (routineId && activeSession && !routineSetupDone.current) {
            routineSetupDone.current = true;
            const setupFromRoutine = async () => {
                try {
                    const routine = await routinesApi.get(routineId);
                    const patterns = routine.movement_patterns || [];
                    setRoutinePatterns(patterns);

                    const allExercises = await exercisesApi.list();
                    const matched = allExercises.filter(ex =>
                        ex.primary_patterns?.some(p => patterns.includes(p.pattern))
                    );

                    if (matched.length > 0) {
                        prepopulateExercises(matched.map(ex => ({ id: ex.id, name: ex.name })));
                        setSelectedExercise(matched[0]);
                    }
                } catch (error) {
                    console.error('Failed to setup from routine:', error);
                }
            };
            setupFromRoutine();
        }
    }, [searchParams, activeSession, prepopulateExercises]);

    useEffect(() => {
        const workoutId = searchParams.get('workoutId');
        if (workoutId && activeSession && !workoutSetupDone.current) {
            workoutSetupDone.current = true;
            const setupFromWorkout = async () => {
                try {
                    const workout = await workoutsApi.get(workoutId);
                    if (workout.Exercises && workout.Exercises.length > 0) {
                        const allExercises = await exercisesApi.list();
                        const exerciseIds = new Set(workout.Exercises.map(ex => ex.ExerciseID));
                        const matched = allExercises.filter(ex => exerciseIds.has(ex.id));

                        if (matched.length > 0) {
                            prepopulateExercises(matched.map(ex => ({ id: ex.id, name: ex.name })));
                            setSelectedExercise(matched[0]);
                        }
                    }
                } catch (error) {
                    console.error('Failed to setup from workout:', error);
                }
            };
            setupFromWorkout();
        }
    }, [searchParams, activeSession, prepopulateExercises]);

    useEffect(() => {
        const loadLastSessionData = async () => {
            try {
                const sessions = await sessionsApi.list(1);
                if (sessions.length > 0) {
                    const full = await sessionsApi.get(sessions[0].session_id);
                    const map = new Map<string, { reps: number; weight: number }>();
                    for (const group of full.exercise_groups) {
                        if (group.sets.length > 0) {
                            const last = group.sets[group.sets.length - 1];
                            map.set(group.exercise_id, { reps: last.reps, weight: last.raw_load });
                        }
                    }
                    setLastSessionMap(map);
                }
            } catch {
                console.warn('No previous session data');
            }
        };
        loadLastSessionData();
    }, []);

    const handleExerciseSelect = (exercise: Exercise) => {
        setSelectedExercise(exercise);
        setCurrentExercise(exercise.id);
        setShowExerciseSelector(false);
    };

    const handleRecordSet = async (data: RecordSetRequest) => {
        if (!activeSession || !selectedExercise) return;

        setIsRecording(true);
        try {
            const response = await sessionsApi.recordSet(activeSession.session_id, data);

            // Add to store
            addPerformedSet(selectedExercise.id, {
                set_id: response.performed_set_id,
                set_number: data.set_number,
                reps: data.reps,
                reps_in_reserve: data.reps_in_reserve,
                raw_load: data.raw_load,
                effective_load: response.effective_load,
                volume: response.volume,
                equipment_id: data.equipment_id,
                notes: data.notes || '',
                performed_at: response.performed_at,
            });

        } catch (error) {
            console.error('Failed to record set:', error);
            alert('Failed to record set. Please try again.');
        } finally {
            setIsRecording(false);
        }
    };

    const handleCompleteSession = async () => {
        if (!activeSession) return;

        const confirmed = window.confirm('Are you sure you want to complete this session?');
        if (!confirmed) return;

        setIsCompleting(true);
        try {
            await sessionsApi.complete(activeSession.session_id);
            clearActiveSession();
            navigate('/dashboard');
        } catch (error) {
            console.error('Failed to complete session:', error);
            alert('Failed to complete session. Please try again.');
        } finally {
            setIsCompleting(false);
        }
    };

    if (!activeSession) {
        return (
            <div className="min-h-screen bg-night-bg flex items-center justify-center">
                <div className="text-white">Starting session...</div>
            </div>
        );
    }

    const currentExerciseGroup = activeSession.exercise_groups.find(
        (g) => g.exercise_id === currentExerciseId
    );

    const nextSetNumber = currentExerciseGroup
        ? currentExerciseGroup.set_count + 1
        : 1;

    return (
        <div className="min-h-screen bg-night-bg pb-20">
            {/* Header */}
            <div className="bg-night-surface border-b border-night-border p-4">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-xl font-semibold text-white">Active Session</h1>
                        <div className="text-night-muted text-sm mt-1">
                            {activeSession.total_sets} sets • {activeSession.total_volume.toFixed(0)}kg volume
                        </div>
                    </div>
                    <button
                        onClick={handleCompleteSession}
                        disabled={isCompleting || activeSession.total_sets === 0}
                        className="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-night-surfaceAlt text-white rounded-lg font-medium transition-colors flex items-center gap-2"
                    >
                        <StopCircle size={18} />
                        Complete
                    </button>
                </div>
            </div>

            <div className="p-4 space-y-6">
                {activeSession.exercise_groups.length > 1 && (
                    <div className="flex gap-2 overflow-x-auto pb-2 -mx-4 px-4">
                        {activeSession.exercise_groups.map((group) => (
                            <button
                                key={group.exercise_id}
                                onClick={() => {
                                    setCurrentExercise(group.exercise_id);
                                    setSelectedExercise({
                                        id: group.exercise_id,
                                        name: group.exercise_name,
                                        description: '',
                                        primary_patterns: [],
                                        secondary_patterns: [],
                                        derived_muscles: { primary: [], secondary: [] },
                                        equipment: [],
                                        source: 'user',
                                        is_custom: true,
                                    });
                                }}
                                className={`px-3 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
                                    currentExerciseId === group.exercise_id
                                        ? 'bg-night-blue text-night-text'
                                        : 'bg-night-surface text-night-muted hover:text-night-text hover:bg-night-surfaceAlt'
                                } ${group.sets.length > 0 ? 'ring-1 ring-green-500/50' : ''}`}
                            >
                                {group.exercise_name || 'Exercise'}
                                {group.set_count > 0 && (
                                    <span className="ml-1.5 text-xs opacity-70">({group.set_count})</span>
                                )}
                            </button>
                        ))}
                    </div>
                )}

                {/* Exercise Selection */}
                {!selectedExercise ? (
                    <button
                        onClick={() => setShowExerciseSelector(true)}
                        className="w-full py-12 bg-night-surface hover:bg-night-surfaceAlt border-2 border-dashed border-night-border rounded-lg text-night-muted hover:text-night-text transition-colors flex flex-col items-center gap-3"
                    >
                        <Plus size={32} />
                        <span className="text-lg font-medium">Select Exercise</span>
                    </button>
                ) : (
                    <>
                        {/* Set Recorder */}
                        <SetRecorder
                            exerciseId={selectedExercise.id}
                            exerciseName={selectedExercise.name}
                            nextSetNumber={nextSetNumber}
                            onRecordSet={handleRecordSet}
                            isLoading={isRecording}
                            equipmentTypes={selectedExercise.equipment}
                            exercisePrimaryPatterns={selectedExercise.primary_patterns?.map(p => p.pattern)}
                            lastSessionReps={lastSessionMap.get(selectedExercise.id)?.reps}
                            lastSessionWeight={lastSessionMap.get(selectedExercise.id)?.weight}
                        />

                        {/* Completed Sets */}
                        {currentExerciseGroup && currentExerciseGroup.sets.length > 0 && (
                            <CompletedSets sets={currentExerciseGroup.sets} />
                        )}

                        {/* Change Exercise Button */}
                        <button
                            onClick={() => setShowExerciseSelector(true)}
                            className="w-full py-3 bg-night-surface hover:bg-night-surfaceAlt text-night-text rounded-lg font-medium transition-colors"
                        >
                            Change Exercise
                        </button>
                    </>
                )}

                {/* All Exercise Groups */}
                {activeSession.exercise_groups.length > 0 && (
                    <div className="mt-8">
                        <h3 className="text-lg font-semibold text-white mb-4">Session Summary</h3>
                        <div className="space-y-4">
                            {activeSession.exercise_groups.map((group) => (
                                <div key={group.exercise_id} className="bg-night-surface rounded-lg p-4">
                                    <div className="flex items-center justify-between">
                                        <div className="text-night-text font-medium">{group.exercise_name || 'Exercise'}</div>
                                        <div className="text-night-muted text-sm">{group.set_count} sets</div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>

            {/* Exercise Selector Modal */}
            {showExerciseSelector && (
                <ExerciseSelector
                    onSelect={handleExerciseSelect}
                    onClose={() => setShowExerciseSelector(false)}
                    allowedPatterns={routinePatterns.length > 0 ? routinePatterns : undefined}
                />
            )}
        </div>
    );
};
