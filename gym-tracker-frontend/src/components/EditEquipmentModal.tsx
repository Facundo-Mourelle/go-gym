import React, { useState } from 'react';
import { X, Wrench, Trash2 } from 'lucide-react';
import { equipmentApi } from '../api/equipment';
import type { EquipmentData } from '../api/equipment';

interface Props {
    equipment: EquipmentData;
    onClose: () => void;
    onUpdated: () => void;
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

export const EditEquipmentModal: React.FC<Props> = ({ equipment, onClose, onUpdated }) => {
    const [name, setName] = useState(equipment.name);
    const [type, setType] = useState<'freeweight' | 'cable' | 'machine'>(equipment.type);
    const [manufacturer, setManufacturer] = useState(equipment.manufacturer || '');
    const [actualWeight, setActualWeight] = useState(equipment.actual_weight?.toString() || '');
    const [pulleyType, setPulleyType] = useState(equipment.pulley_type || '');
    const [stackWeights, setStackWeights] = useState(equipment.stack_weights?.join(', ') || '');
    const [resistanceProfileName, setResistanceProfileName] = useState(equipment.resistance_profile_name || '');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [deleting, setDeleting] = useState(false);

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
                if (stackWeights.trim()) {
                    data.stack_weights = stackWeights
                        .split(',')
                        .map(s => parseFloat(s.trim()))
                        .filter(n => !isNaN(n));
                }
            }

            if (type === 'machine' && resistanceProfileName.trim()) {
                data.resistance_profile_name = resistanceProfileName.trim();
            }

            await equipmentApi.update(equipment.id, data);
            onUpdated();
            onClose();
        } catch (err: any) {
            setError(err.response?.data?.message || err.response?.data || 'Failed to update equipment');
        } finally {
            setSaving(false);
        }
    };

    const handleDelete = async () => {
        setDeleting(true);
        try {
            await equipmentApi.delete(equipment.id);
            onUpdated();
            onClose();
        } catch (err: any) {
            setError(err.response?.data?.message || err.response?.data || 'Failed to delete equipment');
        } finally {
            setDeleting(false);
            setShowDeleteConfirm(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-start sm:items-center justify-center bg-black/60 pt-12 sm:pt-0" onClick={onClose}>
            <div className="bg-gray-800 w-full sm:max-w-lg rounded-t-2xl sm:rounded-2xl max-h-[90vh] flex flex-col" onClick={e => e.stopPropagation()}>
                <div className="flex items-center justify-between p-5 border-b border-gray-700 shrink-0">
                    <div className="flex items-center gap-2">
                        <Wrench className="text-blue-500" size={22} />
                        <h2 className="text-xl font-bold text-white">Edit Equipment</h2>
                    </div>
                    <button onClick={onClose} className="p-1 text-gray-400 hover:text-white transition-colors">
                        <X size={24} />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="p-5 space-y-5 overflow-y-auto">
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
                                <label className="block text-sm font-bold text-gray-400 uppercase tracking-wider mb-2">Stack Weights (comma-separated)</label>
                                <textarea
                                    value={stackWeights}
                                    onChange={e => setStackWeights(e.target.value)}
                                    rows={2}
                                    className="w-full bg-gray-700 border border-gray-600 text-white rounded-xl p-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 resize-none"
                                    placeholder="e.g. 5, 10, 15, 20, 25"
                                />
                                <p className="text-xs text-gray-500 mt-1">Enter weights separated by commas</p>
                            </div>
                        </>
                    )}

                    {type === 'machine' && (
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
                    )}

                    <div className="flex gap-3 pt-2">
                        {showDeleteConfirm ? (
                            <>
                                <button
                                    type="button"
                                    onClick={() => setShowDeleteConfirm(false)}
                                    className="flex-1 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-xl font-medium transition-colors"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="button"
                                    onClick={handleDelete}
                                    disabled={deleting}
                                    className="flex-1 py-3 bg-red-600 hover:bg-red-500 disabled:bg-red-800 text-white rounded-xl font-medium transition-colors"
                                >
                                    {deleting ? 'Deleting...' : 'Confirm Delete'}
                                </button>
                            </>
                        ) : (
                            <>
                                <button
                                    type="button"
                                    onClick={() => setShowDeleteConfirm(true)}
                                    className="py-3 px-4 bg-red-900/50 hover:bg-red-900 text-red-400 rounded-xl font-medium transition-colors flex items-center gap-2"
                                >
                                    <Trash2 size={18} />
                                    Delete
                                </button>
                                <div className="flex-1" />
                                <button
                                    type="button"
                                    onClick={onClose}
                                    className="py-3 px-4 bg-gray-700 hover:bg-gray-600 text-white rounded-xl font-medium transition-colors"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={saving}
                                    className="py-3 px-6 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 text-white rounded-xl font-medium transition-colors"
                                >
                                    {saving ? 'Saving...' : 'Save Changes'}
                                </button>
                            </>
                        )}
                    </div>
                </form>
            </div>
        </div>
    );
};