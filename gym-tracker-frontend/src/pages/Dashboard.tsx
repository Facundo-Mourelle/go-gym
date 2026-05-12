import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { Play, PlusCircle, Clock, Settings } from 'lucide-react';

export const Dashboard: React.FC = () => {
    const navigate = useNavigate();
    const { user, logout } = useAuthStore();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <div className="min-h-full flex flex-col">
            {/* Greeting Header */}
            <div className="px-4 pt-4 pb-3 flex items-center justify-between">
                <h1 className="text-2xl font-bold text-white">Hi, {user?.name} 👋</h1>
                <button
                    onClick={handleLogout}
                    className="p-2 text-night-muted hover:text-night-red transition-colors"
                >
                    <Settings size={20} />
                </button>
            </div>

            {/* Two-column layout */}
            <div className="flex-1 px-4 pb-4 grid grid-cols-3 gap-4 min-h-0">
                {/* Left column: Action buttons (1/3) */}
                <div className="col-span-1 flex flex-col gap-3">
                    <button
                        onClick={() => navigate('/session')}
                        className="flex-1 bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                    >
                        <div className="p-2.5 bg-night-blue/15 rounded-full group-hover:bg-night-blue/25 transition-colors">
                            <Play size={22} className="text-night-blue" />
                        </div>
                        <span className="text-sm font-semibold text-white">New Session</span>
                    </button>

                    <button
                        onClick={() => navigate('/create-workout')}
                        className="flex-1 bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                    >
                        <div className="p-2.5 bg-night-teal/15 rounded-full group-hover:bg-night-teal/25 transition-colors">
                            <PlusCircle size={22} className="text-night-teal" />
                        </div>
                        <span className="text-sm font-semibold text-white">Create Workout</span>
                    </button>

                    <button
                        onClick={() => navigate('/history')}
                        className="flex-1 bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                    >
                        <div className="p-2.5 bg-night-purple/15 rounded-full group-hover:bg-night-purple/25 transition-colors">
                            <Clock size={22} className="text-night-purple" />
                        </div>
                        <span className="text-sm font-semibold text-white">Session History</span>
                    </button>
                </div>

                {/* Right column: Exercise Progress placeholder (2/3) */}
                <div className="col-span-2 bg-night-surface border border-night-border rounded-xl flex flex-col items-center justify-center p-6">
                    <div className="p-3 bg-night-surfaceAlt rounded-full mb-3">
                        <Settings size={24} className="text-night-muted" />
                    </div>
                    <p className="text-night-muted text-sm text-center">
                        Exercise Progress
                    </p>
                    <p className="text-night-muted/50 text-xs text-center mt-1">
                        Coming soon
                    </p>
                </div>
            </div>
        </div>
    );
};
