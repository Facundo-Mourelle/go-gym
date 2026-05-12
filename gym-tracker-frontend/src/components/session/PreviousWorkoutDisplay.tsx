import React, { useEffect, useState } from 'react';
import { sessionsApi } from '../../api/sessions';
import type { Session } from '../../types/session';
import { Clock, Dumbbell, X, Play, AlertCircle } from 'lucide-react';

interface PreviousWorkoutDisplayProps {
    onStartSession: () => void;
    onClose: () => void;
}

export const PreviousWorkoutDisplay: React.FC<PreviousWorkoutDisplayProps> = ({
    onStartSession,
    onClose,
}) => {
    const [session, setSession] = useState<Session | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const load = async () => {
            try {
                const sessions = await sessionsApi.list(1);
                if (sessions.length === 0) {
                    setError('No previous workouts found');
                    setLoading(false);
                    return;
                }
                const full = await sessionsApi.get(sessions[0].session_id);
                setSession(full);
            } catch {
                setError('Failed to load previous workout');
            } finally {
                setLoading(false);
            }
        };
        load();
    }, []);

    const formatDate = (iso: string) => {
        const d = new Date(iso);
        return d.toLocaleDateString('en-US', {
            weekday: 'short', month: 'short', day: 'numeric',
            hour: '2-digit', minute: '2-digit',
        });
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-end md:items-center justify-center z-50">
            <div className="bg-gray-900 w-full md:max-w-2xl md:rounded-lg max-h-[80vh] flex flex-col">
                <div className="p-4 border-b border-gray-700">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                            <div className="p-2 bg-gray-700 rounded-lg">
                                <Clock className="text-blue-400" size={20} />
                            </div>
                            <h2 className="text-xl font-semibold text-white">Previous Workout</h2>
                        </div>
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-white transition-colors"
                        >
                            <X size={24} />
                        </button>
                    </div>
                    {session && (
                        <p className="text-gray-400 text-sm">
                            {formatDate(session.started_at)}
                            {session.total_sets > 0 && (
                                <> &middot; {session.total_sets} sets &middot; {session.total_volume.toFixed(0)}kg volume</>
                            )}
                        </p>
                    )}
                </div>

                <div className="flex-1 overflow-y-auto p-4">
                    {loading ? (
                        <div className="text-center text-gray-400 py-8">Loading...</div>
                    ) : error ? (
                        <div className="flex flex-col items-center justify-center py-12 text-gray-400">
                            <AlertCircle size={40} className="mb-3 text-gray-600" />
                            <p>{error}</p>
                            <p className="text-sm mt-1">Complete a workout to see your history here</p>
                        </div>
                    ) : session ? (
                        <div className="space-y-4">
                            {session.exercise_groups.map((group) => (
                                <div key={group.exercise_id} className="bg-gray-800 rounded-xl border border-gray-700 overflow-hidden">
                                    <div className="flex items-center gap-2 px-4 py-3 bg-gray-800/50 border-b border-gray-700">
                                        <Dumbbell size={16} className="text-blue-400 shrink-0" />
                                        <span className="text-white font-medium">{group.exercise_name}</span>
                                        <span className="text-gray-500 text-xs ml-auto">{group.set_count} sets</span>
                                    </div>

                                    <div className="divide-y divide-gray-700/50">
                                        {group.sets.map((set) => (
                                            <div key={set.set_id} className="flex items-center gap-4 px-4 py-2.5 text-sm">
                                                <span className="text-gray-500 w-12 shrink-0">
                                                    Set {set.set_number}
                                                </span>
                                                <span className="text-white font-medium tabular-nums">
                                                    {set.effective_load}kg
                                                </span>
                                                <span className="text-gray-400">
                                                    x {set.reps}
                                                </span>
                                                {set.reps_in_reserve !== undefined && (
                                                    <span className="text-gray-500 text-xs ml-auto">
                                                        RIR: {set.reps_in_reserve}
                                                    </span>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : null}
                </div>

                <div className="p-4 border-t border-gray-700 flex gap-3">
                    <button
                        onClick={onStartSession}
                        disabled={!session}
                        className="flex-1 py-3 px-4 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 disabled:text-gray-500 text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
                    >
                        <Play size={18} />
                        Start New Session
                    </button>
                    <button
                        onClick={onClose}
                        className="flex-1 py-3 px-4 bg-gray-800 hover:bg-gray-700 text-white font-medium rounded-lg transition-colors"
                    >
                        Close
                    </button>
                </div>
            </div>
        </div>
    );
};
