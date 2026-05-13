import React, { useState, useEffect } from 'react';
import { metricsApi } from '../api/metrics';
import { exercisesApi } from '../api/exercises';
import type { Exercise } from '../types/exercise';
import type { ProgressResponse } from '../api/metrics';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { TrendingUp } from 'lucide-react';

export const Metrics: React.FC = () => {
    const [exercises, setExercises] = useState<Exercise[]>([]);
    const [selectedExercise, setSelectedExercise] = useState<string>('');
    const [progress, setProgress] = useState<ProgressResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [equipmentFilter, setEquipmentFilter] = useState<string>('all');

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

    const filteredData = progress?.all_data_points.filter(
        dp => equipmentFilter === 'all' || dp.equipment_type === equipmentFilter
    ) || [];

    const sortedData = [...filteredData].sort(
        (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
    );

    const chartData = [...filteredData].sort(
        (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
    ).map(dp => ({
        date: dp.date,
        effectiveWeight: dp.weight * (dp.reps + dp.rir),
    }));

    const formatDate = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    };

    return (
        <div className="min-h-screen bg-gray-900 text-white p-6">
            <h1 className="text-3xl font-bold flex items-center gap-2 mb-6">
                <TrendingUp className="text-green-500" />
                Exercise Progress
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
            ) : progress && progress.all_data_points && progress.all_data_points.length > 0 ? (
                <div className="flex gap-6">
                    <div className="w-1/3 bg-gray-800 p-4 rounded-lg max-h-[500px] overflow-y-auto">
                        <h2 className="text-lg font-semibold mb-4 text-gray-300">Set History</h2>
                        <div className="space-y-3">
                            {sortedData.map((dp) => (
                                <div key={dp.set_id} className="bg-gray-700 p-3 rounded-lg">
                                    <div className="text-gray-400 text-sm mb-1">{formatDate(dp.date)}</div>
                                    <div className="text-white font-medium">
                                        {dp.weight}kg × {dp.reps} reps • RIR {dp.rir}
                                    </div>
                                    <div className="text-gray-500 text-xs mt-1 capitalize">{dp.equipment_type}</div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="flex-1 bg-gray-800 p-4 rounded-lg">
                        <div className="flex justify-between items-center mb-4">
                            <h2 className="text-lg font-semibold text-gray-300">Effective Weight</h2>
                            <select
                                className="bg-gray-700 text-white rounded p-2 border border-gray-600 focus:border-blue-500 outline-none text-sm"
                                value={equipmentFilter}
                                onChange={(e) => setEquipmentFilter(e.target.value)}
                            >
                                <option value="all">All Equipment</option>
                                <option value="freeweight">Freeweight</option>
                                <option value="machine">Machine</option>
                                <option value="cable">Cable</option>
                            </select>
                        </div>
                        <div className="h-[420px] w-full">
                            <ResponsiveContainer width="100%" height="100%">
                                <LineChart data={chartData} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                                    <XAxis 
                                        dataKey="date" 
                                        stroke="#9CA3AF" 
                                        tickFormatter={(val) => new Date(val).toLocaleDateString()}
                                    />
                                    <YAxis stroke="#9CA3AF" />
                                    <Tooltip 
                                        contentStyle={{ backgroundColor: '#1F2937', border: 'none', color: '#fff' }}
                                        labelFormatter={(val) => new Date(val).toLocaleDateString()}
                                        formatter={(value) => [`${Number(value).toFixed(1)} kg`, 'Effective Weight']}
                                    />
                                    <Line type="monotone" dataKey="effectiveWeight" stroke="#3B82F6" strokeWidth={2} dot={{ r: 4 }} />
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