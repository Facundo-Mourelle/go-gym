import React from 'react';
import type { PerformedSet } from '../../types/session';
import { CheckCircle } from 'lucide-react';

interface CompletedSetsProps {
    sets: PerformedSet[];
}

export const CompletedSets: React.FC<CompletedSetsProps> = ({ sets }) => {
    if (sets.length === 0) {
        return (
            <div className="text-gray-500 text-center py-8">
                No sets recorded yet
            </div>
        );
    }

    return (
        <div className="space-y-2">
            <h4 className="text-sm font-medium text-gray-400 mb-3">Completed Sets</h4>
            {sets.map((set) => (
                <div
                    key={set.set_id}
                    className="bg-gray-800 rounded-lg p-4 flex items-center justify-between"
                >
                    <div className="flex items-center gap-3">
                        <CheckCircle className="text-green-500" size={20} />
                        <div>
                            <div className="text-white font-medium">
                                Set {set.set_number}
                            </div>
                            <div className="text-gray-400 text-sm">
                                {set.effective_load}kg × {set.reps} reps
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </div>
    );
};
