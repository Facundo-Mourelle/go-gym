import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { sessionsApi } from '../api/sessions';
import type { Session } from '../types/session';
import { Clock, Dumbbell, ChevronRight, Activity } from 'lucide-react';

interface SessionSummary {
    session_id: string;
    started_at: string;
    completed_at?: string;
    duration?: number;
    total_sets: number;
    workout_plan_id?: string;
}

export const History: React.FC = () => {
    const navigate = useNavigate();
    const [sessions, setSessions] = useState<SessionSummary[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedSession, setSelectedSession] = useState<Session | null>(null);
    const [loadingDetail, setLoadingDetail] = useState(false);

    useEffect(() => {
        const fetchSessions = async () => {
            try {
                const data = await sessionsApi.list(50) as unknown as SessionSummary[];
                const completed = data.filter(s => s.completed_at);
                setSessions(completed);
            } catch (error) {
                console.error('Failed to fetch sessions', error);
            } finally {
                setLoading(false);
            }
        };
        fetchSessions();
    }, []);

    const loadSessionDetail = async (sessionId: string) => {
        setLoadingDetail(true);
        try {
            const full = await sessionsApi.get(sessionId);
            setSelectedSession(full);
        } catch (error) {
            console.error('Failed to load session detail', error);
        } finally {
            setLoadingDetail(false);
        }
    };

    const formatDate = (iso: string) => {
        const d = new Date(iso);
        return d.toLocaleDateString('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
        });
    };

    const formatTime = (iso: string) => {
        const d = new Date(iso);
        return d.toLocaleTimeString('en-US', {
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const formatDuration = (nanos?: number): string => {
        if (!nanos) return '—';
        const seconds = Math.floor(nanos / 1_000_000_000);
        if (seconds < 60) return `${seconds}s`;
        const mins = Math.floor(seconds / 60);
        const hrs = Math.floor(mins / 60);
        if (hrs > 0) return `${hrs}h ${mins % 60}m`;
        return `${mins}m`;
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-gray-900 text-white p-6">
                <div className="text-center text-gray-400 mt-10">Loading...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-900 text-white p-6">
            <div className="flex items-center gap-2 mb-6">
                <Activity className="text-blue-500" size={28} />
                <h1 className="text-3xl font-bold">History</h1>
            </div>

            {sessions.length === 0 ? (
                <div className="bg-gray-800 rounded-lg p-6 flex flex-col items-center justify-center min-h-[300px]">
                    <div className="w-16 h-16 rounded-full bg-blue-500/20 flex items-center justify-center mb-4">
                        <Dumbbell className="text-blue-500" size={32} />
                    </div>
                    <h2 className="text-xl font-semibold mb-2">No Sessions Yet</h2>
                    <p className="text-gray-400 mb-6 text-center">
                        Complete a workout to see your session history here
                    </p>
                    <button
                        onClick={() => navigate('/session')}
                        className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors"
                    >
                        <Activity size={18} />
                        <span>Start a Session</span>
                    </button>
                </div>
            ) : (
                <div className="space-y-3">
                    {sessions.map((session) => (
                        <div key={session.session_id}>
                            <button
                                onClick={() =>
                                    selectedSession?.session_id === session.session_id
                                        ? setSelectedSession(null)
                                        : loadSessionDetail(session.session_id)
                                }
                                className="w-full bg-gray-800 hover:bg-gray-750 rounded-lg p-4 flex items-center justify-between transition-colors text-left"
                            >
                                <div className="flex items-center gap-4">
                                    <div className="p-2.5 bg-blue-600/15 rounded-full">
                                        <Dumbbell size={20} className="text-blue-500" />
                                    </div>
                                    <div>
                                        <div className="flex items-center gap-2">
                                            <span className="text-base font-semibold">
                                                {formatDate(session.started_at)}
                                            </span>
                                            <span className="text-sm text-gray-400">
                                                {formatTime(session.started_at)}
                                            </span>
                                        </div>
                                        <div className="flex items-center gap-3 mt-1 text-sm text-gray-400">
                                            <span className="flex items-center gap-1">
                                                <Clock size={14} />
                                                {formatDuration(session.duration)}
                                            </span>
                                            <span>{session.total_sets} sets</span>
                                        </div>
                                    </div>
                                </div>
                                <ChevronRight
                                    size={20}
                                    className={`text-gray-500 transition-transform ${
                                        selectedSession?.session_id === session.session_id
                                            ? 'rotate-90'
                                            : ''
                                    }`}
                                />
                            </button>

                            {selectedSession?.session_id === session.session_id && (
                                <div className="bg-gray-800/50 rounded-b-lg px-4 pb-4 -mt-2 pt-4 border-t border-gray-700/50">
                                    {loadingDetail ? (
                                        <div className="text-center text-gray-400 py-4">
                                            Loading details...
                                        </div>
                                    ) : selectedSession ? (
                                        <div className="space-y-3">
                                            {selectedSession.exercise_groups?.length > 0 ? (
                                                selectedSession.exercise_groups.map((group) => (
                                                    <div
                                                        key={group.exercise_id}
                                                        className="bg-gray-800 rounded-lg p-3"
                                                    >
                                                        <div className="flex items-center justify-between mb-2">
                                                            <span className="font-medium text-white">
                                                                {group.exercise_name || 'Exercise'}
                                                            </span>
                                                            <span className="text-sm text-gray-400">
                                                                {group.set_count} sets
                                                            </span>
                                                        </div>
                                                        <div className="space-y-1">
                                                            {group.sets.map((set) => (
                                                                <div
                                                                    key={set.set_id}
                                                                    className="flex items-center justify-between text-sm text-gray-400 pl-2"
                                                                >
                                                                    <span>
                                                                        Set {set.set_number}:{' '}
                                                                        {set.reps} reps @{' '}
                                                                        {set.raw_load}kg
                                                                    </span>
                                                                    {set.reps_in_reserve !==
                                                                        undefined && (
                                                                        <span className="text-xs text-gray-500">
                                                                            RIR{' '}
                                                                            {set.reps_in_reserve}
                                                                        </span>
                                                                    )}
                                                                </div>
                                                            ))}
                                                        </div>
                                                    </div>
                                                ))
                                            ) : (
                                                <div className="text-center text-gray-500 py-4">
                                                    No exercise data for this session
                                                </div>
                                            )}
                                        </div>
                                    ) : null}
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};
