import React, { useState, useEffect } from 'react';
import type { Exercise } from '../../types/exercise';
import { exercisesApi } from '../../api/exercises';
import { Search, Dumbbell, X } from 'lucide-react';

interface ExerciseSelectorProps {
    onSelect: (exercise: Exercise) => void;
    onClose: () => void;
    allowedPatterns?: string[];
}

export const ExerciseSelector: React.FC<ExerciseSelectorProps> = ({
    onSelect,
    onClose,
    allowedPatterns,
}) => {
    const [exercises, setExercises] = useState<Exercise[]>([]);
    const [search, setSearch] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadExercises();
    }, []);

    const loadExercises = async (search?: string) => {
        try {
            const data = await exercisesApi.list(search ? { search } : undefined);
            setExercises(data);
        } catch (error) {
            console.error('Failed to load exercises:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const timer = setTimeout(() => {
            loadExercises(search || undefined);
        }, 200);
        return () => clearTimeout(timer);
    }, [search]);

    const byPattern = allowedPatterns && allowedPatterns.length > 0
        ? exercises.filter(ex =>
            ex.primary_patterns?.some(p => allowedPatterns.includes(p.pattern))
          )
        : exercises;

    const filteredExercises = byPattern;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-end md:items-center justify-center z-50">
            <div className="bg-gray-900 w-full md:max-w-2xl md:rounded-lg max-h-[80vh] flex flex-col">
                {/* Header */}
                <div className="p-4 border-b border-gray-700">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-xl font-semibold text-white">Select Exercise</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-white transition-colors"
                        >
                            <X size={24} />
                        </button>
                    </div>

                    {/* Search */}
                    <div className="relative">
                        <Search
                            className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                            size={20}
                        />
                        <input
                            type="text"
                            placeholder="Search exercises..."
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            className="w-full pl-10 pr-4 py-3 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            autoFocus
                        />
                    </div>

                    {/* Pattern Filter Info Banner */}
                    {allowedPatterns && allowedPatterns.length > 0 && (
                        <div className="mt-3 px-3 py-2 bg-blue-900/30 border border-blue-800 rounded-lg">
                            <p className="text-blue-300 text-xs">
                                Showing exercises matching {allowedPatterns.length} movement pattern{allowedPatterns.length !== 1 ? 's' : ''}
                            </p>
                        </div>
                    )}
                </div>

                {/* Exercise List */}
                <div className="flex-1 overflow-y-auto p-4">
                    {loading ? (
                        <div className="text-center text-gray-400 py-8">Loading exercises...</div>
                    ) : filteredExercises.length === 0 ? (
                        <div className="text-center text-gray-400 py-8">No exercises found</div>
                    ) : (
                        <div className="space-y-2">
                            {filteredExercises.map((exercise) => (
                                <button
                                    key={exercise.id}
                                    onClick={() => onSelect(exercise)}
                                    className="w-full text-left p-4 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors"
                                >
                                    <div className="flex items-start gap-3">
                                        <Dumbbell className="text-blue-500 mt-1" size={20} />
                                        <div className="flex-1">
                                            <div className="text-white font-medium">{exercise.name}</div>
                                            <div className="text-gray-400 text-sm mt-1">
                                                {exercise.derived_muscles.primary.join(', ')}
                                            </div>
                                            <div className="flex gap-2 mt-2">
                                                {exercise.primary_patterns.slice(0, 2).map((pattern) => (
                                                    <span
                                                        key={pattern.pattern}
                                                        className="text-xs bg-blue-900 text-blue-200 px-2 py-1 rounded"
                                                    >
                                                        {pattern.pattern_name}
                                                    </span>
                                                ))}
                                            </div>
                                        </div>
                                    </div>
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};
