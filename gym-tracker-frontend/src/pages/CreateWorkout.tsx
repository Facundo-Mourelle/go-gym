import React from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Dumbbell } from 'lucide-react';

export const CreateWorkout: React.FC = () => {
    const navigate = useNavigate();

    return (
        <div className="p-4">
            <button
                onClick={() => navigate('/dashboard')}
                className="flex items-center gap-2 text-night-muted hover:text-night-text transition-colors mb-4"
            >
                <ArrowLeft size={20} />
                <span className="text-sm">Back</span>
            </button>

            <div className="flex flex-col items-center justify-center mt-20 text-center">
                <div className="p-4 bg-night-surfaceAlt rounded-full mb-4">
                    <Dumbbell size={32} className="text-night-blue" />
                </div>
                <h2 className="text-xl font-bold text-night-text mb-2">Create Workout</h2>
                <p className="text-night-muted text-sm max-w-xs">
                    Workout templates coming soon. You'll be able to create and save custom routines here.
                </p>
            </div>
        </div>
    );
};
