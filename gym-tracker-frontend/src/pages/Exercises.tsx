import React, { useState, useEffect } from 'react';
import { Dumbbell, Search, Plus, ListFilter, Pencil } from 'lucide-react';
import { exercisesApi } from '../api/exercises';
import { CreateExerciseModal } from '../components/CreateExerciseModal';
import { EditExerciseModal } from '../components/EditExerciseModal';
import type { Exercise } from '../types/exercise';
import type { PatternInfo } from '../api/exercises';

export const Exercises: React.FC = () => {
    const [exercises, setExercises] = useState<Exercise[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedPattern, setSelectedPattern] = useState('');
    const [patterns, setPatterns] = useState<PatternInfo[]>([]);
    const [showCreate, setShowCreate] = useState(false);
    const [editingExercise, setEditingExercise] = useState<Exercise | null>(null);

    const fetchExercises = async (search?: string, pattern?: string) => {
        setLoading(true);
        try {
            const filters: { search?: string; pattern?: string } = {};
            if (search) filters.search = search;
            if (pattern) filters.pattern = pattern;
            const data = await exercisesApi.list(Object.keys(filters).length ? filters : undefined);
            setExercises(data);
        } catch (err) {
            console.error("Failed to fetch exercises", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        exercisesApi.listPatterns()
            .then(patterns => {
                const sorted = [...patterns].sort((a, b) => a.name.localeCompare(b.name));
                setPatterns(sorted);
            })
            .catch(err => console.error('Failed to load patterns:', err));
    }, []);

    useEffect(() => {
        const timer = setTimeout(() => {
            fetchExercises(searchTerm || undefined, selectedPattern || undefined);
        }, 200);
        return () => clearTimeout(timer);
    }, [searchTerm, selectedPattern]);

    const filteredExercises = exercises;

    return (
        <div className="bg-gray-900 h-screen flex flex-col text-white">
            <div className="p-4 bg-gray-800 border-b border-gray-700 shadow-md">
                <div className="flex items-center justify-between mb-4">
                    <h1 className="text-2xl font-bold flex items-center gap-2">
                        <Dumbbell className="text-blue-500" />
                        Exercise Dictionary
                    </h1>
                    <button
                        onClick={() => setShowCreate(true)}
                        className="flex items-center gap-1.5 bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-xl font-semibold transition-colors text-sm"
                    >
                        <Plus size={18} />
                        <span>Create</span>
                    </button>
                </div>
                
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Search className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                        type="text"
                        className="bg-gray-900 border border-gray-700 text-white text-sm rounded-xl focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-3 outline-none transition-colors"
                        placeholder="Search exercises..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>

                <div className="mt-3 flex items-center gap-2">
                    <ListFilter size={16} className="text-gray-500 shrink-0" />
                    <select
                        value={selectedPattern}
                        onChange={e => setSelectedPattern(e.target.value)}
                        className="bg-gray-900 border border-gray-700 text-white text-sm rounded-xl px-3 py-2 outline-none focus:ring-blue-500 focus:border-blue-500 transition-colors"
                    >
                        <option value="">All patterns</option>
                        {patterns.map(p => (
                            <option key={p.pattern} value={p.pattern}>{p.name}</option>
                        ))}
                    </select>
                    {selectedPattern && (
                        <button
                            onClick={() => setSelectedPattern('')}
                            className="text-xs text-blue-400 hover:text-blue-300 transition-colors"
                        >
                            Clear
                        </button>
                    )}
                </div>
            </div>

            <div className="flex-1 overflow-y-auto p-4">
                {loading ? (
                    <div className="flex justify-center p-10">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                    </div>
                ) : filteredExercises.length === 0 ? (
                    <div className="text-center text-gray-400 mt-10 bg-gray-800 p-8 rounded-2xl border border-gray-700/50">
                        <Dumbbell className="mx-auto h-12 w-12 text-gray-600 mb-3" />
                        <p className="font-medium">No exercises found.</p>
                        <p className="text-sm mt-1">Try a different search term.</p>
                    </div>
                ) : (
                    <div className="grid gap-3">
                        {filteredExercises.map((ex) => (
                            <div key={ex.id} className="bg-gray-800 rounded-2xl p-4 border border-gray-700/50 hover:border-gray-600 transition-colors group">
                                <div className="flex items-start justify-between gap-2">
                                    <h3 className="font-bold text-lg">{ex.name}</h3>
                                    {(ex.source === 'user' || ex.is_custom) && (
                                        <button
                                            onClick={() => setEditingExercise(ex)}
                                            className="p-1.5 text-gray-500 hover:text-blue-400 hover:bg-blue-400/10 rounded-lg transition-colors opacity-0 group-hover:opacity-100 shrink-0"
                                            title="Edit exercise"
                                        >
                                            <Pencil size={16} />
                                        </button>
                                    )}
                                </div>
                                {ex.description && (
                                    <p className="text-gray-400 text-sm mt-1 line-clamp-2">{ex.description}</p>
                                )}
                                <div className="mt-3 flex flex-wrap gap-2">
                                    {ex.primary_patterns?.map((p, i) => (
                                        <span key={i} className="bg-blue-900/50 text-blue-300 text-xs px-2 py-1 rounded-md border border-blue-800/50">
                                            {p.pattern_name || p.pattern}
                                        </span>
                                    ))}
                                    {ex.equipment && ex.equipment.length > 0 && (
                                        <span className="bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded-md border border-gray-600">
                                            {ex.equipment.join(', ')}
                                        </span>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {showCreate && (
                <CreateExerciseModal
                    onClose={() => setShowCreate(false)}
                    onCreated={fetchExercises}
                />
            )}

            {editingExercise && (
                <EditExerciseModal
                    exercise={editingExercise}
                    onClose={() => setEditingExercise(null)}
                    onUpdated={fetchExercises}
                />
            )}
        </div>
    );
};
