import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { workoutsApi } from '../api/workouts';
import type { WorkoutSummary } from '../types/workout';
import { Plus, Calendar, Dumbbell, ChevronRight, List } from 'lucide-react';

export const Workouts: React.FC = () => {
    const navigate = useNavigate();
    const [workouts, setWorkouts] = useState<WorkoutSummary[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchWorkouts = async () => {
            try {
                const data = await workoutsApi.list();
                setWorkouts(data);
            } catch (error) {
                console.error("Failed to fetch workouts", error);
            } finally {
                setLoading(false);
            }
        };
        fetchWorkouts();
    }, []);

    return (
        <div className="min-h-screen bg-gray-900 text-white p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold flex items-center gap-2">
                    <Dumbbell className="text-blue-500" />
                    Workout Plans
                </h1>
                <button
                    onClick={() => navigate('/create-workout')}
                    className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-semibold transition-colors"
                >
                    <Plus size={18} />
                    <span>Create Workout</span>
                </button>
            </div>
            
            {loading ? (
                <div className="text-center text-gray-400 mt-10">Loading...</div>
            ) : workouts.length === 0 ? (
                <div className="bg-gray-800 rounded-lg p-6 flex flex-col items-center justify-center min-h-[300px]">
                    <div className="w-16 h-16 rounded-full bg-blue-500/20 flex items-center justify-center mb-4">
                        <Dumbbell className="text-blue-500" size={32} />
                    </div>
                    <h2 className="text-xl font-semibold mb-2">No Workout Plans Yet</h2>
                    <p className="text-gray-400 mb-6">Create your first workout plan to get started</p>
                    <button
                        onClick={() => navigate('/create-workout')}
                        className="flex items-center space-x-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-semibold transition-colors"
                    >
                        <Plus size={18} />
                        <span>Create Workout</span>
                    </button>
                </div>
            ) : (
                <div className="space-y-4">
                    {workouts.map((workout) => (
                        <div 
                            key={workout.ID} 
                            onClick={() => navigate(`/create-workout?id=${workout.ID}`)}
                            className="bg-gray-800 p-5 rounded-lg border border-gray-700 hover:border-blue-500/50 hover:bg-gray-750 transition-all cursor-pointer"
                        >
                            <div className="flex justify-between items-start">
                                <div className="flex-1">
                                    <h2 className="text-xl font-semibold">
                                        {workout.Name}
                                    </h2>
                                    <div className="flex items-center text-gray-400 mt-2 gap-4 text-sm">
                                        <div className="flex items-center gap-1">
                                            <Dumbbell size={14} />
                                            <span>{workout.TotalSets} sets</span>
                                        </div>
                                        <div className="flex items-center gap-1">
                                            <List size={14} />
                                            <span>{workout.ExerciseCount} exercises</span>
                                        </div>
                                    </div>
                                    <div className="flex items-center text-gray-500 mt-2 text-sm">
                                        <Calendar size={14} className="mr-1" />
                                        <span>Created {new Date(workout.CreatedAt).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })}</span>
                                    </div>
                                </div>
                                <ChevronRight className="text-gray-500" size={20} />
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};