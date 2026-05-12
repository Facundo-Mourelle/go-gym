import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { sessionsApi } from '../api/sessions';
import type { Session } from '../types/session';
import { Play, History, Calendar } from 'lucide-react';

export const Workouts: React.FC = () => {
    const navigate = useNavigate();
    const [sessions, setSessions] = useState<Session[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchSessions = async () => {
            try {
                const data = await sessionsApi.list();
                setSessions(data);
            } catch (error) {
                console.error("Failed to fetch sessions", error);
            } finally {
                setLoading(false);
            }
        };
        fetchSessions();
    }, []);

    return (
        <div className="min-h-screen bg-gray-900 text-white p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold flex items-center gap-2">
                    <History className="text-blue-500" />
                    Workout History
                </h1>
                <button
                    onClick={() => navigate('/session')}
                    className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-semibold transition-colors"
                >
                    <Play size={18} />
                    <span>Start New Workout</span>
                </button>
            </div>
            
            {loading ? (
                <div className="text-center text-gray-400 mt-10">Loading history...</div>
            ) : sessions.length === 0 ? (
                <div className="bg-gray-800 rounded-lg p-6 flex flex-col items-center justify-center min-h-[300px]">
                    <p className="text-gray-400 mb-4">No recent workouts found.</p>
                </div>
            ) : (
                <div className="space-y-4">
                    {sessions.map((session) => (
                        <div key={session.session_id} className="bg-gray-800 p-5 rounded-lg border border-gray-700">
                            <div className="flex justify-between items-start">
                                <div>
                                    <h2 className="text-xl font-semibold">
                                        {session.workout_plan_id ? `Workout Plan: ${session.workout_plan_id}` : 'Freestyle Workout'}
                                    </h2>
                                    <div className="flex items-center text-gray-400 mt-2 gap-2 text-sm">
                                        <Calendar size={14} />
                                        <span>{new Date(session.started_at).toLocaleDateString()}</span>
                                        <span>•</span>
                                        <span>{session.total_sets || 0} Sets</span>
                                        <span>•</span>
                                        <span>{session.total_volume || 0} kg Volume</span>
                                    </div>
                                </div>
                                {session.completed_at ? (
                                    <span className="bg-green-900 text-green-300 text-xs px-2 py-1 rounded-full border border-green-700">Completed</span>
                                ) : (
                                    <span className="bg-yellow-900 text-yellow-300 text-xs px-2 py-1 rounded-full border border-yellow-700">In Progress</span>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};
