import React, { useState } from 'react';
import { X, Wrench } from 'lucide-react';
import { equipmentApi } from '../api/equipment';

interface Props {
    onClose: () => void;
    onCreated: () => void;
}

const TYPE_OPTIONS = [
    { value: 'freeweight', label: 'Free Weight' },
    { value: 'cable', label: 'Cable' },
    { value: 'machine', label: 'Machine' },
] as const;

const PULLEY_OPTIONS = [
    { value: '2:1', label: '2:1' },
    { value: '1:1', label: '1:1' },
    { value: '4:1', label: '4:1' },
];

const MOVEMENT_PATTERN_OPTIONS = [
    { value: 'horizontal_push', label: 'Horizontal Push (Chest Press)' },
    { value: 'horizontal_pull', label: 'Horizontal Pull (Rows)' },
    { value: 'vertical_push', label: 'Vertical Push (Overhead Press)' },
    { value: 'vertical_pull', label: 'Vertical Pull (Lat Pulldown)' },
    { value: 'shoulder_flexion', label: 'Shoulder Flexion' },
    { value: 'shoulder_extension', label: 'Shoulder Extension' },
    { value: 'shoulder_abduction', label: 'Shoulder Abduction (Lateral Raise)' },
    { value: 'shoulder_adduction', label: 'Shoulder Adduction' },
    { value: 'hip_hinge', label: 'Hip Hinge (Deadlift)' },
    { value: 'hip_adduction', label: 'Hip Adduction' },
    { value: 'squat_pattern', label: 'Squat Pattern' },
    { value: 'knee_extension', label: 'Knee Extension' },
    { value: 'elbow_flexion', label: 'Elbow Flexion (Curls)' },
    { value: 'elbow_extension', label: 'Elbow Extension (Triceps)' },
    { value: 'spinal_flexion', label: 'Spinal Flexion' },
    { value: 'spinal_extension', label: 'Spinal Extension' },
];

export const CreateEquipmentModal: React.FC<Props> = ({ onClose, onCreated }) => {
    const [name, setName] = useState('');
    const [type, setType] = useState<'freeweight' | 'cable' | 'machine'>('freeweight');
    const [manufacturer, setManufacturer] = useState('');
    const [actualWeight, setActualWeight] = useState('');
    const [pulleyType, setPulleyType] = useState('');
    const [weightIncrement, setWeightIncrement] = useState('');
    const [resistanceProfileName, setResistanceProfileName] = useState('');
    const [movementPattern, setMovementPattern] = useState('');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!name.trim()) {
            setError('Name is required');
            return;
        }

        setSaving(true);
        try {
            const data: any = {
                name: name.trim(),
                type,
                manufacturer: manufacturer.trim() || undefined,
            };

            if (type === 'freeweight' && actualWeight) {
                data.actual_weight = parseFloat(actualWeight);
            }

            if (type === 'cable') {
                if (pulleyType) data.pulley_type = pulleyType;
                if (weightIncrement) data.weight_increment = parseFloat(weightIncrement);
            }

            if (type === 'machine') {
                if (resistanceProfileName.trim()) data.resistance_profile_name = resistanceProfileName.trim();
                if (movementPattern) data.movement_pattern = movementPattern;
            }

            await equipmentApi.create(data);
            onCreated();
            onClose();
        } catch (err: any) {
            setError(err.response?.data?.message || err.response?.data || 'Failed to create equipment');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/60" onClick={onClose}>
            <div className="bg-gray-800 w-full sm:max-w-lg rounded-t-2xl sm:rounded-2xl max-h-[90vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
                <div className="flex items-center justify-between p-5 border-b border-gray-700 sticky top-0 bg-gray-800 z-10">
                    <div className="flex items-center gap-2">
                        <Wrench className="text-blue-500" size={22} />
                        <h2 className="text-xl font-bold text-white">New Equipment</h2>
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

                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Type</label>
                        <div className="flex gap-2">
                            {TYPE_OPTIONS.map(opt => (
                                <button
                                    key={opt.value}
                                    type="button"
                                    onClick={() => setType(opt.value)}
                                    className={`flex-1 py-2.5 rounded-xl text-sm font-medium border transition-colors ${
                                        type === opt.value
                                            ? 'bg-blue-600 border-blue-400 text-white ring-2 ring-blue-400/50'
                                            : 'bg-gray-700 border-gray-600 text-gray-300 hover:border-gray-500'
                                    }`}
                                >
                                    {opt.label}
                                </button>
                            ))}
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Name</label>
                        <input
                            type="text"
                            value={name}
                            onChange={e => setName(e.target.value)}
                            className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                            placeholder="e.g. Barbell 20kg"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Manufacturer</label>
                        <input
                            type="text"
                            value={manufacturer}
                            onChange={e => setManufacturer(e.target.value)}
                            className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                            placeholder="e.g. Rogue, Eleiko"
                        />
                    </div>

                    
                    {type === 'freeweight' && (
                        <div>
                            <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Actual Weight (kg)</label>
                            <input
                                type="number"
                                step="0.5"
                                value={actualWeight}
                                onChange={e => setActualWeight(e.target.value)}
                                className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                                placeholder="e.g. 20"
                            />
                        </div>
                    )}

                    {type === 'cable' && (
                        <>
                            <div>
                                <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Pulley Type</label>
                                <select
                                    value={pulleyType}
                                    onChange={e => setPulleyType(e.target.value)}
                                    className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                                >
                                    <option value="">Select pulley type...</option>
                                    {PULLEY_OPTIONS.map(opt => (
                                        <option key={opt.value} value={opt.value}>{opt.label}</option>
                                    ))}
                                </select>
                            </div>
                            <div>
                                <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Weight Increase</label>
                                <select
                                    value={weightIncrement}
                                    onChange={e => setWeightIncrement(e.target.value)}
                                    className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                                >
                                    <option value="">Select weight increase...</option>
                                    <option value="1.25">1.25 kg</option>
                                    <option value="2.5">2.5 kg</option>
                                    <option value="4.5">4.5 kg (10 lb)</option>
                                </select>
                                <p className="text-xs text-gray-500 mt-1">The smallest increment between plates in the stack</p>
                            </div>
                        </>
                    )}

                    {type === 'machine' && (
                        <>
                            <div>
                                <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Resistance Profile Name</label>
                                <input
                                    type="text"
                                    value={resistanceProfileName}
                                    onChange={e => setResistanceProfileName(e.target.value)}
                                    className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                                    placeholder="e.g. Cybex VR2"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Movement Pattern</label>
                                <select
                                    value={movementPattern}
                                    onChange={e => setMovementPattern(e.target.value)}
                                    className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                                >
                                    <option value="">Select movement pattern...</option>
                                    {MOVEMENT_PATTERN_OPTIONS.map(opt => (
                                        <option key={opt.value} value={opt.value}>{opt.label}</option>
                                    ))}
                                </select>
                                <p className="text-xs text-gray-500 mt-1">Which movement pattern this machine is designed for</p>
                            </div>
                        </>
                    )}

                    <button
                        type="submit"
                        disabled={saving}
                        className="w-full py-4 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 text-white font-bold rounded-xl transition-colors text-lg"
                    >
                        {saving ? 'Creating...' : 'Create Equipment'}
                    </button>
                </form>
            </div>
        </div>
    );
};