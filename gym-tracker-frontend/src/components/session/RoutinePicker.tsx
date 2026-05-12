import React, { useState, useEffect } from 'react';
import { routinesApi } from '../../api/routines';
import type { Routine } from '../../types/routine';
import { ListFilter, Dumbbell, X } from 'lucide-react';

interface RoutinePickerProps {
    onSelect: (routineId: string) => void;
    onFreestyle: () => void;
    onClose: () => void;
}

export const RoutinePicker: React.FC<RoutinePickerProps> = ({
    onSelect,
    onFreestyle,
    onClose,
}) => {
    const [routines, setRoutines] = useState<Routine[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadRoutines();
    }, []);

    const loadRoutines = async () => {
        try {
            const data = await routinesApi.list();
            setRoutines(data);
        } catch (error) {
            console.error('Failed to load routines:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-end md:items-center justify-center z-50">
            <div className="bg-gray-900 w-full md:max-w-2xl md:rounded-lg max-h-[80vh] flex flex-col">
                {/* Header */}
                <div className="p-4 border-b border-gray-700">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                            <div className="p-2 bg-blue-600 rounded-lg">
                                <ListFilter className="text-white" size={20} />
                            </div>
                            <h2 className="text-xl font-semibold text-white">Choose a Routine</h2>
                        </div>
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-white transition-colors"
                        >
                            <X size={24} />
                        </button>
                    </div>
                    <p className="text-gray-400 text-sm">
                        Select a workout template based on movement patterns, or start freestyle.
                    </p>
                </div>

                {/* Routine List */}
                <div className="flex-1 overflow-y-auto p-4">
                    {loading ? (
                        <div className="text-center text-gray-400 py-8">Loading routines...</div>
                    ) : routines.length === 0 ? (
                        <div className="text-center text-gray-400 py-8">
                            <Dumbbell className="mx-auto mb-2" size={40} />
                            <p>No routines found</p>
                            <p className="text-sm mt-1">Create a routine to get started</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {routines.map((routine) => (
                                <button
                                    key={routine.id}
                                    onClick={() => onSelect(routine.id)}
                                    className="w-full text-left p-4 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded-lg transition-colors"
                                >
                                    <div className="flex items-start justify-between">
                                        <div className="flex-1">
                                            <div className="flex items-center gap-2">
                                                <div className="text-white font-medium">{routine.name}</div>
                                                {routine.is_preset && (
                                                    <span className="px-2 py-0.5 bg-blue-900/50 text-blue-300 text-xs rounded-full">
                                                        Preset
                                                    </span>
                                                )}
                                            </div>
                                            {routine.description && (
                                                <div className="text-gray-400 text-sm mt-1">
                                                    {routine.description}
                                                </div>
                                            )}
                                            <div className="text-gray-500 text-xs mt-2">
                                                {(routine.movement_patterns ?? []).length} movement pattern{(routine.movement_patterns ?? []).length !== 1 ? 's' : ''}
                                            </div>
                                            <div className="flex flex-wrap gap-1 mt-2">
                                                {(routine.movement_patterns ?? []).slice(0, 4).map((pattern) => (
                                                    <span
                                                        key={pattern}
                                                        className="px-2 py-0.5 bg-blue-900/50 text-blue-300 text-xs rounded"
                                                    >
                                                        {pattern}
                                                    </span>
                                                ))}
                                                {(routine.movement_patterns ?? []).length > 4 && (
                                                    <span className="px-2 py-0.5 bg-gray-700 text-gray-400 text-xs rounded">
                                                        +{(routine.movement_patterns ?? []).length - 4}
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                        <Dumbbell className="text-gray-600" size={20} />
                                    </div>
                                </button>
                            ))}
                        </div>
                    )}
                </div>

                {/* Footer Actions */}
                <div className="p-4 border-t border-gray-700 flex gap-3">
                    <button
                        onClick={onFreestyle}
                        className="flex-1 py-3 px-4 bg-gray-800 hover:bg-gray-700 text-white font-medium rounded-lg transition-colors"
                    >
                        Skip (Freestyle)
                    </button>
                    <button
                        onClick={onClose}
                        className="flex-1 py-3 px-4 bg-gray-700 hover:bg-gray-600 text-white font-medium rounded-lg transition-colors"
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    );
};