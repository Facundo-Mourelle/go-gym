import React, { useState, useEffect } from 'react';
import type { RecordSetRequest } from '../../types/session';
import type { EquipmentType } from '../../types/exercise';
import { equipmentApi, type EquipmentItem } from '../../api/equipment';

interface SetRecorderProps {
    exerciseId: string;
    exerciseName: string;
    nextSetNumber: number;
    onRecordSet: (data: RecordSetRequest) => Promise<void>;
    isLoading: boolean;
    initialReps?: number;
    initialWeight?: number;
    lastSessionReps?: number;
    lastSessionWeight?: number;
    equipmentTypes?: EquipmentType[];
}

export const SetRecorder: React.FC<SetRecorderProps> = ({
    exerciseId,
    exerciseName,
    nextSetNumber,
    onRecordSet,
    isLoading,
    initialReps,
    initialWeight,
    lastSessionReps,
    lastSessionWeight,
    equipmentTypes,
}) => {
    const [allEquipment, setAllEquipment] = useState<EquipmentItem[]>([]);
    const [equipmentId, setEquipmentId] = useState<string>('');
    const [reps, setReps] = useState<number | ''>(initialReps ?? '');
    const [weight, setWeight] = useState<number>(initialWeight ?? lastSessionWeight ?? 0);
    const [rir, setRir] = useState<number>(2);

    useEffect(() => {
        equipmentApi.list().then((items) => {
            setAllEquipment(items);
            const compatible = items.filter(
                (eq) => !equipmentTypes || equipmentTypes.length === 0 || equipmentTypes.includes(eq.type as EquipmentType)
            );
            if (compatible.length > 0 && !equipmentId) {
                setEquipmentId(compatible[0].id);
            }
        });
    }, []);

    const availableEquipment = allEquipment.filter(
        (eq) => !equipmentTypes || equipmentTypes.length === 0 || equipmentTypes.includes(eq.type as EquipmentType)
    );

    const canSubmit = reps !== '' && reps > 0 && equipmentId !== '' && !isLoading;

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (reps === '' || reps <= 0 || equipmentId === '') return;

        await onRecordSet({
            exercise_id: exerciseId,
            set_number: nextSetNumber,
            reps,
            reps_in_reserve: rir,
            raw_load: weight,
            equipment_id: equipmentId,
            notes: `RIR: ${rir}`
        });

        setReps(Math.max(1, reps - 1));
    };

    return (
        <div className="bg-gray-800 rounded-2xl p-6 border border-gray-700/50 shadow-lg">
            <div className="flex justify-between items-start mb-4">
                <div>
                    <h3 className="text-xl font-bold text-white">{exerciseName}</h3>
                    <div className="text-gray-400 text-sm font-medium mt-1 uppercase tracking-wide">Set {nextSetNumber}</div>
                </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-5">
                <div>
                    <label className="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">
                        Equipment
                    </label>
                    <div className="flex flex-wrap gap-2">
                        {availableEquipment.map((eq) => (
                            <button
                                key={eq.id}
                                type="button"
                                onClick={() => setEquipmentId(eq.id)}
                                disabled={isLoading}
                                className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                                    equipmentId === eq.id
                                        ? 'bg-blue-600 text-white'
                                        : 'bg-gray-700 text-gray-300 hover:bg-gray-600 hover:text-white'
                                }`}
                            >
                                {eq.name}
                            </button>
                        ))}
                    </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">
                            Weight (kg)
                        </label>
                        <div className="flex items-center bg-gray-700 rounded-xl border border-gray-600 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 overflow-hidden">
                            <button type="button" onClick={() => setWeight(w => Math.max(0, w - 2.5))} className="px-4 py-3 text-gray-400 hover:text-white hover:bg-gray-600">-</button>
                            <input
                                type="number"
                                step="0.5"
                                value={weight}
                                placeholder={lastSessionWeight ? `Last: ${lastSessionWeight}` : ''}
                                onChange={(e) => setWeight(parseFloat(e.target.value))}
                                className="w-full bg-transparent text-center text-white text-xl font-bold focus:outline-none appearance-none placeholder:text-gray-500"
                                disabled={isLoading}
                            />
                            <button type="button" onClick={() => setWeight(w => w + 2.5)} className="px-4 py-3 text-gray-400 hover:text-white hover:bg-gray-600">+</button>
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-bold text-gray-400 uppercase tracking-wider mb-2">
                            Reps
                        </label>
                        <div className="flex items-center bg-gray-700 rounded-xl border border-gray-600 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 overflow-hidden">
                            <button type="button" onClick={() => setReps(r => Math.max(1, (r || 0) - 1))} className="px-4 py-3 text-gray-400 hover:text-white hover:bg-gray-600">-</button>
                            <input
                                type="number"
                                min="1"
                                value={reps}
                                placeholder={lastSessionReps ? `Last: ${lastSessionReps}` : ''}
                                onChange={(e) => {
                                    const v = e.target.value;
                                    setReps(v === '' ? '' : parseInt(v));
                                }}
                                className="w-full bg-transparent text-center text-white text-xl font-bold focus:outline-none placeholder:text-gray-500"
                                disabled={isLoading}
                            />
                            <button type="button" onClick={() => setReps(r => (r || 0) + 1)} className="px-4 py-3 text-gray-400 hover:text-white hover:bg-gray-600">+</button>
                        </div>
                    </div>
                </div>

                <div>
                    <div className="flex justify-between items-end mb-2">
                        <label className="block text-xs font-bold text-gray-400 uppercase tracking-wider">
                            RIR (Reps In Reserve)
                        </label>
                        <span className="text-white font-bold">{rir}</span>
                    </div>
                    <input
                        type="range"
                        min="0"
                        max="5"
                        step="1"
                        value={rir}
                        onChange={(e) => setRir(parseInt(e.target.value))}
                        className="w-full h-2 bg-gray-600 rounded-lg appearance-none cursor-pointer accent-blue-500"
                        disabled={isLoading}
                    />
                    <div className="flex justify-between text-[10px] text-gray-500 mt-1 px-1 font-medium">
                        <span>0 (Failure)</span>
                        <span>2-3</span>
                        <span>5 (Easy)</span>
                    </div>
                </div>

                <button
                    type="submit"
                    disabled={!canSubmit}
                    className="w-full py-4 mt-2 bg-blue-600 hover:bg-blue-500 active:bg-blue-700 disabled:bg-gray-700 disabled:text-gray-400 text-white font-bold rounded-xl transition-colors text-lg shadow-lg shadow-blue-900/20"
                >
                    {isLoading ? 'Recording...' : 'Record Set'}
                </button>
            </form>
        </div>
    );
};
