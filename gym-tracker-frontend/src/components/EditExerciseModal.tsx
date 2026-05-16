import React, { useState, useEffect } from 'react';
import { X, Dumbbell, Check } from 'lucide-react';
import { exercisesApi } from '../api/exercises';
import type { Exercise } from '../types/exercise';
import type { PatternInfo } from '../api/exercises';

interface Props {
    exercise: Exercise;
    onClose: () => void;
    onUpdated: () => void;
}

const EQUIPMENT_OPTIONS = [
    { value: 'barbell', label: 'Barbell' },
    { value: 'dumbbell', label: 'Dumbbell' },
    { value: 'machine', label: 'Machine' },
    { value: 'cable', label: 'Cable' },
    { value: 'bodyweight', label: 'Bodyweight' },
];

export const EditExerciseModal: React.FC<Props> = ({ exercise, onClose, onUpdated }) => {
    const [name, setName] = useState(exercise.name);
    const [patterns, setPatterns] = useState<PatternInfo[]>([]);
    const [selectedPattern, setSelectedPattern] = useState(
        exercise.primary_patterns?.[0]?.pattern ?? ''
    );
    const [selectedEquipment, setSelectedEquipment] = useState<string[]>([...exercise.equipment]);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        exercisesApi.listPatterns()
            .then(patterns => {
                const sorted = [...patterns].sort((a, b) => a.name.localeCompare(b.name));
                setPatterns(sorted);
            })
            .catch(err => {
                console.error('Failed to load patterns:', err);
                setError('Failed to load movement patterns.');
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

        if (!name.trim()) { setError('Name is required'); return; }
        if (!selectedPattern) { setError('Select a movement pattern'); return; }

        setSaving(true);
        try {
            await exercisesApi.update(exercise.id, {
                name: name.trim(),
                primary_patterns: [{
                    pattern: selectedPattern,
                    contribution: 1.0,
                    range_of_motion: 'full',
                    notes: '',
                }],
                equipment: selectedEquipment,
            });
            onUpdated();
            onClose();
        } catch (err: unknown) {
            const data = (err as { response?: { data?: string } })?.response?.data ?? 'Failed to update exercise';
            setError(data);
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-start sm:items-center justify-center bg-black/60 pt-12 sm:pt-0" onClick={onClose}>
            <div className="bg-night-surface w-full sm:max-w-lg rounded-t-2xl sm:rounded-2xl max-h-[90vh] flex flex-col" onClick={e => e.stopPropagation()}>
                {/* Header */}
                <div className="flex items-center justify-between p-5 border-b border-night-border shrink-0">
                    <div className="flex items-center gap-2">
                        <Dumbbell className="text-night-blue" size={22} />
                        <h2 className="text-xl font-bold text-white">Edit Exercise</h2>
                    </div>
                    <button onClick={onClose} className="p-1 text-night-muted hover:text-white transition-colors">
                        <X size={24} />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="p-5 space-y-5 overflow-y-auto">
                    {error && (
                        <div className="bg-night-red/10 border border-night-red/30 text-night-red text-sm p-3 rounded-xl">
                            {error}
                        </div>
                    )}

                    {/* Name */}
                    <div>
                        <label className="block text-sm font-bold text-night-muted uppercase tracking-wider mb-2">Name</label>
                        <input
                            type="text"
                            value={name}
                            onChange={e => setName(e.target.value)}
                            className="w-full bg-night-surfaceAlt border border-night-border text-white rounded-xl p-3 outline-none focus:border-night-blue focus:ring-1 focus:ring-night-blue"
                            placeholder="Exercise name"
                        />
                    </div>

                    {/* Movement Pattern - scrollable list */}
                    <div>
                        <label className="block text-sm font-bold text-night-muted uppercase tracking-wider mb-2">Movement Pattern</label>
                        <div className="max-h-48 overflow-y-auto space-y-1 rounded-xl border border-night-border bg-night-surfaceAlt p-1">
                            {patterns.map(p => (
                                <button
                                    key={p.pattern}
                                    type="button"
                                    onClick={() => setSelectedPattern(p.pattern)}
                                    className={`w-full flex items-center justify-between px-3 py-2.5 rounded-lg text-sm text-left transition-colors ${
                                        selectedPattern === p.pattern
                                            ? 'bg-night-blue/20 text-night-blue'
                                            : 'text-night-text hover:bg-night-surfaceAlt/50'
                                    }`}
                                >
                                    <div>
                                        <span className="font-medium">{p.name}</span>
                                        <p className="text-xs text-night-muted mt-0.5">{p.description}</p>
                                    </div>
                                    {selectedPattern === p.pattern && (
                                        <Check size={18} className="shrink-0" />
                                    )}
                                </button>
                            ))}
                            {patterns.length === 0 && (
                                <p className="text-night-muted text-sm text-center py-4">Loading patterns...</p>
                            )}
                        </div>
                    </div>

                    {/* Equipment - scrollable list */}
                    <div>
                        <label className="block text-sm font-bold text-night-muted uppercase tracking-wider mb-2">Equipment</label>
                        <div className="max-h-40 overflow-y-auto space-y-1 rounded-xl border border-night-border bg-night-surfaceAlt p-1">
                            {EQUIPMENT_OPTIONS.map(eq => {
                                const isSelected = selectedEquipment.includes(eq.value);
                                return (
                                    <button
                                        key={eq.value}
                                        type="button"
                                        onClick={() => toggleEquipment(eq.value)}
                                        className={`w-full flex items-center justify-between px-3 py-2.5 rounded-lg text-sm text-left transition-colors ${
                                            isSelected
                                                ? 'bg-night-teal/20 text-night-teal'
                                                : 'text-night-text hover:bg-night-surfaceAlt/50'
                                        }`}
                                    >
                                        <span className="font-medium capitalize">{eq.label}</span>
                                        {isSelected && (
                                            <Check size={18} className="shrink-0" />
                                        )}
                                    </button>
                                );
                            })}
                        </div>
                    </div>

                    {/* Actions */}
                    <div className="flex gap-3 pt-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="flex-1 py-3 bg-night-surfaceAlt hover:bg-night-surfaceAlt/80 text-night-text rounded-xl font-medium transition-colors"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={saving}
                            className="flex-1 py-3 bg-night-blue hover:bg-night-blue/80 disabled:bg-night-surfaceAlt text-white rounded-xl font-medium transition-colors"
                        >
                            {saving ? 'Saving...' : 'Save Changes'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};
