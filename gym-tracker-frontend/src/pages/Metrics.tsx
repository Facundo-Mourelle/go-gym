import React, { useState, useEffect } from 'react';
import { metricsApi } from '../api/metrics';
import { exercisesApi } from '../api/exercises';
import type { Exercise } from '../types/exercise';
import type { ProgressResponse } from '../api/metrics';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { TrendingUp } from 'lucide-react';

export const Metrics: React.FC = () => {
    const [exercises, setExercises] = useState<Exercise[]>([]);
    const [selectedExercise, setSelectedExercise] = useState<string>('');
    const [progress, setProgress] = useState<ProgressResponse | null>(null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const fetchExercises = async () => {
            try {
                const data = await exercisesApi.list();
                setExercises(data);
                if (data.length > 0) {
                    setSelectedExercise(data[0].id);
                }
            } catch (err) {
                console.error("Failed to fetch exercises", err);
            }
        };
        fetchExercises();
    }, []);

    useEffect(() => {
        if (!selectedExercise) return;

        const fetchProgress = async () => {
            setLoading(true);
            try {
                const data = await metricsApi.getExerciseProgress(selectedExercise);
                setProgress(data);
            } catch (err) {
                console.error("Failed to fetch progress", err);
                setProgress(null);
            } finally {
                setLoading(false);
            }
        };
        fetchProgress();
    }, [selectedExercise]);

    return (
        <div className="min-h-screen bg-gray-900 text-white p-6">
            <h1 className="text-3xl font-bold flex items-center gap-2 mb-6">
                <TrendingUp className="text-green-500" />
                Progress Metrics
            </h1>

            <div className="bg-gray-800 p-6 rounded-lg mb-6">
                <label className="block text-sm font-medium text-gray-400 mb-2">Select Exercise</label>
                <select
                    className="w-full bg-gray-700 text-white rounded p-3 border border-gray-600 focus:border-blue-500 outline-none"
                    value={selectedExercise}
                    onChange={(e) => setSelectedExercise(e.target.value)}
                >
                    {exercises.map(ex => (
                        <option key={ex.id} value={ex.id}>{ex.name}</option>
                    ))}
                </select>
            </div>

            {loading ? (
                <div className="text-center text-gray-400">Loading metrics...</div>
            ) : progress && progress.Data && progress.Data.length > 0 ? (
                <div className="space-y-6">
                    <div className="grid grid-cols-2 gap-4">
                        <div className="bg-gray-800 p-4 rounded-lg">
                            <div className="text-gray-400 text-sm">Overall Trend</div>
                            <div className="text-2xl font-bold capitalize text-green-400">{progress.OverallTrend}</div>
                        </div>
                        <div className="bg-gray-800 p-4 rounded-lg">
                            <div className="text-gray-400 text-sm">Improvement</div>
                            <div className="text-2xl font-bold text-blue-400">{progress.ImprovementPercentage.toFixed(1)}%</div>
                        </div>
                    </div>

                    <div className="bg-gray-800 p-6 rounded-lg">
                        <h2 className="text-xl font-semibold mb-4 text-center">Estimated 1RM (kg)</h2>
                        <div className="h-64 w-full">
                            <ResponsiveContainer width="100%" height="100%">
                                <LineChart data={progress.Data} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                                    <XAxis 
                                        dataKey="Date" 
                                        stroke="#9CA3AF" 
                                        tickFormatter={(val) => new Date(val).toLocaleDateString()}
                                    />
                                    <YAxis stroke="#9CA3AF" />
                                    <Tooltip 
                                        contentStyle={{ backgroundColor: '#1F2937', border: 'none', color: '#fff' }}
                                        labelFormatter={(val) => new Date(val).toLocaleDateString()}
                                    />
                                    <Legend />
                                    <Line type="monotone" dataKey="Estimated1RM" name="1RM (Epley)" stroke="#3B82F6" strokeWidth={3} dot={{ r: 4 }} activeDot={{ r: 6 }} />
                                    <Line type="monotone" dataKey="MaxWeight" name="Max Weight" stroke="#10B981" strokeWidth={2} />
                                </LineChart>
                            </ResponsiveContainer>
                        </div>
                    </div>
                </div>
            ) : (
                <div className="text-center text-gray-400 py-10 bg-gray-800 rounded-lg">
                    No data recorded yet for this exercise.
                </div>
            )}
        </div>
    );
};
