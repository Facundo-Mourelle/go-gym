import React, { useState, useEffect } from 'react';
import { X, Dumbbell, Check } from 'lucide-react';
import { exercisesApi } from '../api/exercises';
import type { PatternInfo } from '../api/exercises';

interface Props {
    onClose: () => void;
    onCreated: () => void;
}

const EQUIPMENT_OPTIONS = [
    { value: 'barbell', label: 'Barbell' },
    { value: 'dumbbell', label: 'Dumbbell' },
    { value: 'machine', label: 'Machine' },
    { value: 'cable', label: 'Cable' },
    { value: 'bodyweight', label: 'Bodyweight' },
];

export const CreateExerciseModal: React.FC<Props> = ({ onClose, onCreated }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [patterns, setPatterns] = useState<PatternInfo[]>([]);
    const [selectedPattern, setSelectedPattern] = useState('');
    const [selectedEquipment, setSelectedEquipment] = useState<string[]>([]);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');
    const [patternsError, setPatternsError] = useState('');

    useEffect(() => {
        exercisesApi.listPatterns()
            .then(patterns => {
                const sorted = [...patterns].sort((a, b) => a.name.localeCompare(b.name));
                setPatterns(sorted);
            })
            .catch(err => {
                console.error('Failed to load patterns:', err);
                setPatternsError('Failed to load movement patterns. Is the backend running?');
            });
    }, []);

    const toggleEquipment = (eq: string) => {
        setSelectedEquipment(prev =>
            prev.includes(eq) ? prev.filter(e => e !== eq) : [...prev, eq]
        );
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!name.trim()) { setError('Name required'); return; }
        if (!selectedPattern) { setError('Select primary movement pattern'); return; }

        setSaving(true);
        try {
            await exercisesApi.create({
                name: name.trim(),
                description: description.trim(),
                primary_patterns: [{
                    pattern: selectedPattern,
                    contribution: 1.0,
                    range_of_motion: 'full',
                    notes: '',
                }],
                equipment: selectedEquipment,
            });
            onCreated();
            onClose();
        } catch (err: unknown) {
            const data = (err as { response?: { data?: string } })?.response?.data ?? 'Failed to create exercise';
            setError(data);
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/60" onClick={onClose}>
            <div className="bg-gray-800 w-full sm:max-w-lg rounded-t-2xl sm:rounded-2xl max-h-[90vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
                {/* Header */}
                <div className="flex items-center justify-between p-5 border-b border-gray-700 sticky top-0 bg-gray-800 z-10">
                    <div className="flex items-center gap-2">
                        <Dumbbell className="text-blue-500" size={22} />
                        <h2 className="text-xl font-bold text-white">New Exercise</h2>
                    </div>
                    <button onClick={onClose} className="p-1 text-gray-400 hover:text-white transition-colors">
                        <X size={24} />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="p-5 space-y-5">
                    {error && (
                        <div className="bg-red-900/50 border border-red-700 text-red-300 text-sm p-3 rounded-xl">
                            {error}
                        </div>
                    )}

                    {/* Name */}
                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Name</label>
                        <input
                            type="text"
                            value={name}
                            onChange={e => setName(e.target.value)}
                            className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                            placeholder="e.g. Incline Dumbbell Press"
                        />
                    </div>

                    {/* Description */}
                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Description</label>
                        <textarea
                            value={description}
                            onChange={e => setDescription(e.target.value)}
                            rows={3}
                            className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 resize-none"
                            placeholder="Optional notes..."
                        />
                    </div>

                    {/* Pattern */}
                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Primary Movement Pattern</label>
                        {patternsError ? (
                            <div className="text-red-400 text-sm bg-red-400/10 rounded-xl p-3 border border-red-400/30">
                                {patternsError}
                            </div>
                        ) : (
                            <>
                                <select
                                    value={selectedPattern}
                                    onChange={e => setSelectedPattern(e.target.value)}
                                    className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                                >
                                    <option value="">Select pattern...</option>
                                    {patterns.map(p => (
                                        <option key={p.pattern} value={p.pattern}>{p.name}</option>
                                    ))}
                                </select>
                                {selectedPattern && (
                                    <p className="text-xs text-gray-400 mt-1">
                                        {patterns.find(p => p.pattern === selectedPattern)?.description}
                                    </p>
                                )}
                            </>
                        )}
                    </div>

                    {/* Equipment */}
                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Equipment</label>
                        <div className="flex flex-wrap gap-2">
                            {EQUIPMENT_OPTIONS.map(eq => (
                                <button
                                    key={eq.value}
                                    type="button"
                                    onClick={() => toggleEquipment(eq.value)}
                                    className={`px-4 py-2 rounded-xl text-sm font-medium border transition-colors flex items-center gap-2 ${
                                        selectedEquipment.includes(eq.value)
                                            ? 'bg-blue-600 border-blue-400 text-white ring-2 ring-blue-400/50'
                                            : 'bg-gray-700 border-gray-600 text-gray-300 hover:border-gray-500 hover:bg-gray-600'
                                    }`}
                                >
                                    {selectedEquipment.includes(eq.value) && <Check size={14} className="text-white" />}
                                    {eq.label}
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Submit */}
                    <button
                        type="submit"
                        disabled={saving}
                        className="w-full py-4 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 text-white font-bold rounded-xl transition-colors text-lg"
                    >
                        {saving ? 'Creating...' : 'Create Exercise'}
                    </button>
                </form>
            </div>
        </div>
    );
};
