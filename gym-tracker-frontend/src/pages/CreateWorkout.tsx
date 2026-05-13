import React, { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { ArrowLeft, Plus, Trash2, Pencil, Save, Search, X } from 'lucide-react';
import { workoutsApi } from '../api/workouts';
import { exercisesApi } from '../api/exercises';
import type { WorkoutPlan } from '../types/workout';
import type { Exercise } from '../types/exercise';

// Helper to map equipment type to display label and color
const getEquipmentBadge = (equipmentType: string): { label: string; colorClass: string } => {
    switch (equipmentType) {
        case 'freeweight':
            return { label: 'Free Weight', colorClass: 'bg-blue-600' };
        case 'machine':
            return { label: 'Machine', colorClass: 'bg-green-600' };
        case 'cable':
            return { label: 'Cable', colorClass: 'bg-purple-600' };
        default:
            return { label: equipmentType, colorClass: 'bg-gray-600' };
    }
};

// Interface for local exercise state (used in form)
interface WorkoutExerciseFormData {
    id: string; // temp ID for new exercises
    exerciseId: string;
    exerciseName: string;
    equipmentType: string;
    sets: number;
    repsInReserve: number;
    isEditing: boolean;
}

export const CreateWorkout: React.FC = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const workoutId = searchParams.get('id');
    const isEditMode = searchParams.get('edit') === 'true';
    const isViewMode = !!workoutId && !isEditMode;

    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [workout, setWorkout] = useState<WorkoutPlan | null>(null);
    const [error, setError] = useState<string | null>(null);

    // Form state
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [exercises, setExercises] = useState<WorkoutExerciseFormData[]>([]);

    // Modal state
    const [showExerciseModal, setShowExerciseModal] = useState(false);
    const [availableExercises, setAvailableExercises] = useState<Exercise[]>([]);
    const [exerciseSearch, setExerciseSearch] = useState('');
    const [loadingExercises, setLoadingExercises] = useState(false);
    const [deleting, setDeleting] = useState(false);

    // Fetch workout if in view or edit mode
    useEffect(() => {
        if (workoutId) {
            fetchWorkout();
        }
    }, [workoutId]);

    const fetchWorkout = async () => {
        if (!workoutId) return;
        setLoading(true);
        setError(null);
        try {
            const data = await workoutsApi.get(workoutId);
            setWorkout(data);
            // Populate form for edit mode
            if (isEditMode) {
                setName(data.Name);
                setDescription(data.Description);
                setExercises(data.Exercises.map((ex, idx) => ({
                    id: ex.ID || `temp-${idx}`,
                    exerciseId: ex.ExerciseID,
                    exerciseName: ex.ExerciseName,
                    equipmentType: ex.EquipmentType,
                    sets: ex.Sets,
                    repsInReserve: ex.RepsInReserve,
                    isEditing: false,
                })));
            }
        } catch (err) {
            console.error('Failed to fetch workout:', err);
            setError('Failed to load workout');
        } finally {
            setLoading(false);
        }
    };

    // Open exercise modal
    const openExerciseModal = async () => {
        setShowExerciseModal(true);
        setLoadingExercises(true);
        setExerciseSearch('');
        try {
            const data = await exercisesApi.list();
            setAvailableExercises(data);
        } catch (err) {
            console.error('Failed to fetch exercises:', err);
        } finally {
            setLoadingExercises(false);
        }
    };

    // Add exercise from modal
    const addExercise = (exercise: Exercise) => {
        const equipmentType = exercise.equipment?.[0] || 'freeweight';
        const newExercise: WorkoutExerciseFormData = {
            id: `temp-${Date.now()}`,
            exerciseId: exercise.id,
            exerciseName: exercise.name,
            equipmentType,
            sets: 3,
            repsInReserve: 0,
            isEditing: false,
        };
        setExercises([...exercises, newExercise]);
        setShowExerciseModal(false);
    };

    // Remove exercise
    const removeExercise = (id: string) => {
        setExercises(exercises.filter(ex => ex.id !== id));
    };

    // Toggle edit mode for an exercise
    const toggleExerciseEdit = (id: string) => {
        setExercises(exercises.map(ex =>
            ex.id === id ? { ...ex, isEditing: !ex.isEditing } : ex
        ));
    };

    // Update exercise field
    const updateExerciseField = (id: string, field: 'sets' | 'repsInReserve', value: number) => {
        setExercises(exercises.map(ex =>
            ex.id === id ? { ...ex, [field]: value } : ex
        ));
    };

    // Save workout (create or update)
    const saveWorkout = async () => {
        if (!name.trim()) {
            setError('Workout name is required');
            return;
        }

        setSaving(true);
        setError(null);

        try {
            const exercisesPayload = exercises.map((ex, idx) => ({
                order: idx,
                exercise_id: ex.exerciseId,
                sets: ex.sets,
                reps: 0,
                reps_in_reserve: ex.repsInReserve,
            }));

            const payload: Record<string, unknown> = {
                name: name.trim(),
                description: description.trim() || undefined,
                exercises: exercisesPayload,
            };

            let savedWorkout: WorkoutPlan;
            if (workoutId && isEditMode) {
                savedWorkout = await workoutsApi.update(workoutId, payload);
            } else {
                savedWorkout = await workoutsApi.create(payload);
            }

            // Navigate to view mode
            navigate(`/create-workout?id=${savedWorkout.ID}`);
        } catch (err) {
            console.error('Failed to save workout:', err);
            setError('Failed to save workout');
        } finally {
            setSaving(false);
        }
    };

    const handleDelete = async () => {
        if (!workoutId) return;
        if (!window.confirm('Are you sure you want to delete this workout plan?')) return;

        setDeleting(true);
        try {
            await workoutsApi.delete(workoutId);
            navigate('/history');
        } catch (err) {
            console.error('Failed to delete workout:', err);
            setError('Failed to delete workout');
            setDeleting(false);
        }
    };

    // Filter exercises in modal
    const filteredExercises = availableExercises.filter(ex =>
        ex.name.toLowerCase().includes(exerciseSearch.toLowerCase())
    );

    // Loading state
    if (loading) {
        return (
            <div className="p-6 flex items-center justify-center min-h-screen">
                <div className="text-gray-400">Loading...</div>
            </div>
        );
    }

    // View Mode
    if (isViewMode && workout) {
        return (
            <div className="p-6 min-h-screen bg-gray-900 text-white">
                <button
                    onClick={() => navigate('/history')}
                    className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-6"
                >
                    <ArrowLeft size={20} />
                    <span className="text-sm">Back to Workouts</span>
                </button>

                <div className="max-w-3xl mx-auto">
                    <div className="flex items-start justify-between mb-6">
                        <div>
                            <h1 className="text-3xl font-bold mb-2">{workout.Name}</h1>
                            {workout.Description && (
                                <p className="text-gray-400">{workout.Description}</p>
                            )}
                            <p className="text-sm text-gray-500 mt-2">
                                Created: {new Date(workout.CreatedAt).toLocaleDateString()}
                            </p>
                        </div>
                        <div className="flex items-center gap-2">
                            <button
                                onClick={() => navigate(`/create-workout?id=${workout.ID}&edit=true`)}
                                className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
                            >
                                <Pencil size={18} />
                                Edit
                            </button>
                            <button
                                onClick={handleDelete}
                                disabled={deleting}
                                className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-800 rounded-lg transition-colors"
                            >
                                <Trash2 size={18} />
                                {deleting ? 'Deleting...' : 'Delete'}
                            </button>
                        </div>
                    </div>

                    {workout.Exercises.length > 0 && (
                        <div className="mt-8">
                            <h2 className="text-xl font-semibold mb-4">Exercises</h2>
                            <div className="space-y-3">
                                {workout.Exercises.map((exercise) => {
                                    const badge = getEquipmentBadge(exercise.EquipmentType);
                                    return (
                                        <div
                                            key={exercise.ID}
                                            className="bg-gray-800 rounded-lg p-4 border border-gray-700"
                                        >
                                            <div className="flex items-center justify-between">
                                                <div>
                                                    <h3 className="font-medium text-lg">{exercise.ExerciseName}</h3>
                                                    <div className="flex items-center gap-3 mt-2">
                                                        <span className="text-gray-400">
                                                            {exercise.Sets} sets
                                                        </span>
                                                        {exercise.RepsInReserve > 0 && (
                                                            <span className="text-gray-400">
                                                                {exercise.RepsInReserve} RIR
                                                            </span>
                                                        )}
                                                        <span className={`px-2 py-1 rounded text-xs font-medium ${badge.colorClass}`}>
                                                            {badge.label}
                                                        </span>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    )}

                    {workout.Exercises.length === 0 && (
                        <div className="text-center py-12 text-gray-500">
                            No exercises in this workout
                        </div>
                    )}
                </div>
            </div>
        );
    }

    // Edit Mode or Create Mode
    return (
        <div className="p-6 min-h-screen bg-gray-900 text-white">
            <button
                onClick={() => navigate('/history')}
                className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors mb-6"
            >
                <ArrowLeft size={20} />
                <span className="text-sm">Back to Workouts</span>
            </button>

            <div className="max-w-3xl mx-auto">
                <div className="flex items-center justify-between mb-6">
                    <h1 className="text-2xl font-bold">
                        {workoutId && isEditMode ? 'Edit Workout' : 'Create Workout'}
                    </h1>
                    {workoutId && isEditMode && (
                        <button
                            onClick={handleDelete}
                            disabled={deleting}
                            className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-800 rounded-lg transition-colors"
                        >
                            <Trash2 size={18} />
                            {deleting ? 'Deleting...' : 'Delete'}
                        </button>
                    )}
                </div>

                {error && (
                    <div className="mb-4 p-3 bg-red-900/50 border border-red-700 rounded-lg text-red-200">
                        {error}
                    </div>
                )}

                {/* Form Fields */}
                <div className="space-y-4 mb-8">
                    <div>
                        <label className="block text-sm font-medium text-gray-300 mb-2">
                            Name *
                        </label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="e.g., Push Day"
                            className="w-full bg-gray-800 border border-gray-700 rounded-lg p-3 text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-300 mb-2">
                            Description
                        </label>
                        <textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="Optional description..."
                            rows={3}
                            className="w-full bg-gray-800 border border-gray-700 rounded-lg p-3 text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 resize-none"
                        />
                    </div>
                </div>

                {/* Exercises Section */}
                <div className="mb-8">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-xl font-semibold">Exercises</h2>
                        <button
                            onClick={openExerciseModal}
                            className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
                        >
                            <Plus size={18} />
                            Import Exercise
                        </button>
                    </div>

                    {exercises.length === 0 ? (
                        <div className="text-center py-8 text-gray-500 bg-gray-800/50 rounded-lg border border-gray-700 border-dashed">
                            No exercises added yet. Click "Import Exercise" to add one.
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {exercises.map((exercise) => {
                                const badge = getEquipmentBadge(exercise.equipmentType);
                                return (
                                    <div
                                        key={exercise.id}
                                        className="bg-gray-800 rounded-lg p-4 border border-gray-700"
                                    >
                                        <div className="flex items-center justify-between">
                                            <div className="flex-1">
                                                <div className="flex items-center gap-3 mb-3">
                                                    <h3 className="font-medium">{exercise.exerciseName}</h3>
                                                    <span className={`px-2 py-1 rounded text-xs font-medium ${badge.colorClass}`}>
                                                        {badge.label}
                                                    </span>
                                                </div>

                                                {exercise.isEditing ? (
                                                    <div className="flex items-center gap-4">
                                                        <div className="flex items-center gap-2">
                                                            <label className="text-sm text-gray-400">Sets:</label>
                                                            <input
                                                                type="number"
                                                                min={1}
                                                                value={exercise.sets}
                                                                onChange={(e) => updateExerciseField(exercise.id, 'sets', parseInt(e.target.value) || 1)}
                                                                className="w-16 bg-gray-700 border border-gray-600 rounded px-2 py-1 text-white text-center"
                                                            />
                                                        </div>
                                                        <div className="flex items-center gap-2">
                                                            <label className="text-sm text-gray-400">RIR:</label>
                                                            <input
                                                                type="number"
                                                                min={0}
                                                                value={exercise.repsInReserve}
                                                                onChange={(e) => updateExerciseField(exercise.id, 'repsInReserve', parseInt(e.target.value) || 0)}
                                                                className="w-16 bg-gray-700 border border-gray-600 rounded px-2 py-1 text-white text-center"
                                                            />
                                                        </div>
                                                        <button
                                                            onClick={() => toggleExerciseEdit(exercise.id)}
                                                            className="text-green-400 hover:text-green-300"
                                                        >
                                                            <Save size={18} />
                                                        </button>
                                                    </div>
                                                ) : (
                                                    <div className="flex items-center gap-4 text-sm text-gray-400">
                                                        <span>{exercise.sets} sets</span>
                                                        <span>{exercise.repsInReserve} RIR</span>
                                                    </div>
                                                )}
                                            </div>

                                            <div className="flex items-center gap-2 ml-4">
                                                <button
                                                    onClick={() => toggleExerciseEdit(exercise.id)}
                                                    className="p-2 text-gray-400 hover:text-white transition-colors"
                                                    title="Edit"
                                                >
                                                    <Pencil size={18} />
                                                </button>
                                                <button
                                                    onClick={() => removeExercise(exercise.id)}
                                                    className="p-2 text-gray-400 hover:text-red-400 transition-colors"
                                                    title="Remove"
                                                >
                                                    <Trash2 size={18} />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </div>

                {/* Save Button */}
                <button
                    onClick={saveWorkout}
                    disabled={saving}
                    className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-800 rounded-lg transition-colors font-medium"
                >
                    {saving ? (
                        'Saving...'
                    ) : (
                        <>
                            <Save size={20} />
                            {workoutId && isEditMode ? 'Update Workout' : 'Save Workout'}
                        </>
                    )}
                </button>
            </div>

            {/* Exercise Import Modal */}
            {showExerciseModal && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-gray-800 rounded-lg w-full max-w-md max-h-[80vh] flex flex-col border border-gray-700">
                        <div className="flex items-center justify-between p-4 border-b border-gray-700">
                            <h2 className="text-lg font-semibold">Import Exercise</h2>
                            <button
                                onClick={() => setShowExerciseModal(false)}
                                className="text-gray-400 hover:text-white"
                            >
                                <X size={20} />
                            </button>
                        </div>

                        <div className="p-4 border-b border-gray-700">
                            <div className="relative">
                                <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                                <input
                                    type="text"
                                    value={exerciseSearch}
                                    onChange={(e) => setExerciseSearch(e.target.value)}
                                    placeholder="Search exercises..."
                                    className="w-full bg-gray-700 border border-gray-600 rounded-lg pl-10 pr-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:border-blue-500"
                                />
                            </div>
                        </div>

                        <div className="flex-1 overflow-y-auto p-2">
                            {loadingExercises ? (
                                <div className="text-center py-8 text-gray-400">Loading...</div>
                            ) : filteredExercises.length === 0 ? (
                                <div className="text-center py-8 text-gray-400">No exercises found</div>
                            ) : (
                                <div className="space-y-1">
                                    {filteredExercises.map((exercise) => {
                                        const equipment = exercise.equipment?.[0] || 'freeweight';
                                        const badge = getEquipmentBadge(equipment);
                                        return (
                                            <button
                                                key={exercise.id}
                                                onClick={() => addExercise(exercise)}
                                                className="w-full text-left p-3 rounded-lg hover:bg-gray-700 transition-colors"
                                            >
                                                <div className="flex items-center justify-between">
                                                    <span className="font-medium">{exercise.name}</span>
                                                    <span className={`px-2 py-1 rounded text-xs font-medium ${badge.colorClass}`}>
                                                        {badge.label}
                                                    </span>
                                                </div>
                                            </button>
                                        );
                                    })}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};