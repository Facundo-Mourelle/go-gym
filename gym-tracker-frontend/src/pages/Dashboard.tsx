import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { Play, PlusCircle, Clock, List, Dumbbell, LogOut } from 'lucide-react';

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
                    <LogOut size={20} />
                </button>
            </div>

            {/* Action buttons grid */}
            <div className="flex-1 px-4 pb-4 grid grid-cols-4 gap-3 min-h-0">
                {/* Select From Workout */}
                <button
                    onClick={() => navigate('/workouts')}
                    className="bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                >
                    <div className="p-2.5 bg-night-orange/15 rounded-full group-hover:bg-night-orange/25 transition-colors">
                        <List size={22} className="text-night-orange" />
                    </div>
                    <span className="text-sm font-semibold text-white">Select From Workout</span>
                </button>

                {/* Freestyle */}
                <button
                    onClick={() => navigate('/session')}
                    className="bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                >
                    <div className="p-2.5 bg-night-blue/15 rounded-full group-hover:bg-night-blue/25 transition-colors">
                        <Play size={22} className="text-night-blue" />
                    </div>
                    <span className="text-sm font-semibold text-white">Freestyle</span>
                </button>

                {/* Create Workout */}
                <button
                    onClick={() => navigate('/create-workout')}
                    className="bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                >
                    <div className="p-2.5 bg-night-teal/15 rounded-full group-hover:bg-night-teal/25 transition-colors">
                        <PlusCircle size={22} className="text-night-teal" />
                    </div>
                    <span className="text-sm font-semibold text-white">Create Workout</span>
                </button>

                {/* Session History */}
                <button
                    onClick={() => navigate('/history')}
                    className="bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                >
                    <div className="p-2.5 bg-night-purple/15 rounded-full group-hover:bg-night-purple/25 transition-colors">
                        <Clock size={22} className="text-night-purple" />
                    </div>
                    <span className="text-sm font-semibold text-white">Session History</span>
                </button>

                {/* View Workouts */}
                <button
                    onClick={() => navigate('/workouts')}
                    className="bg-night-surface hover:bg-night-surfaceAlt border border-night-border rounded-xl flex flex-col items-center justify-center gap-2 p-4 transition-all hover:scale-[1.02] active:scale-95 group"
                >
                    <div className="p-2.5 bg-night-green/15 rounded-full group-hover:bg-night-green/25 transition-colors">
                        <Dumbbell size={22} className="text-night-green" />
                    </div>
                    <span className="text-sm font-semibold text-white">View Workouts</span>
                </button>
            </div>
        </div>
    );
};
