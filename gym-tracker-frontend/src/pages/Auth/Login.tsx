import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { authApi } from '../../api/auth';
import { Dumbbell, Mail, Lock, AlertCircle } from 'lucide-react';

export const Login: React.FC = () => {
    const navigate = useNavigate();
    const { setAuth } = useAuthStore();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!email || !password) {
            setError('Please fill in all fields');
            return;
        }

        setIsLoading(true);

        try {
            const response = await authApi.login({ email, password });
            setAuth(
                {
                    user_id: response.user_id,
                    email: response.email,
                    name: response.name,
                },
                response.token
            );
            navigate('/dashboard');
        } catch (err: any) {
            console.error('Login error:', err);
            setError(err.response?.data?.message || 'Invalid email or password');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-night-bg flex items-center justify-center p-4">
            <div className="w-full max-w-md">
                {/* Logo/Header */}
                <div className="text-center mb-8">
                    <div className="flex items-center justify-center gap-3 mb-4">
                        <Dumbbell className="text-night-blue" size={40} />
                        <h1 className="text-3xl font-bold text-night-text">Gym Tracker</h1>
                    </div>
                    <p className="text-night-muted">Track your progress, reach your goals</p>
                </div>

                {/* Login Card */}
                <div className="bg-night-surface rounded-lg p-8 shadow-xl">
                    <h2 className="text-2xl font-semibold text-night-text mb-6">Login</h2>

                    {/* Error Message */}
                    {error && (
                        <div className="mb-4 p-4 bg-red-900 bg-opacity-50 border border-red-500 rounded-lg flex items-center gap-3">
                            <AlertCircle className="text-red-400" size={20} />
                            <p className="text-red-200 text-sm">{error}</p>
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        {/* Email Input */}
                        <div>
                            <label className="block text-sm font-medium text-night-text mb-2">
                                Email
                            </label>
                            <div className="relative">
                                <Mail
                                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-night-muted"
                                    size={20}
                                />
                                <input
                                    type="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    className="w-full pl-10 pr-4 py-3 bg-night-surfaceAlt border border-night-border rounded-lg text-night-text placeholder-night-muted focus:outline-none focus:ring-2 focus:ring-night-blue"
                                    placeholder="john@example.com"
                                    disabled={isLoading}
                                    autoComplete="email"
                                />
                            </div>
                        </div>

                        {/* Password Input */}
                        <div>
                            <label className="block text-sm font-medium text-night-text mb-2">
                                Password
                            </label>
                            <div className="relative">
                                <Lock
                                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-night-muted"
                                    size={20}
                                />
                                <input
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    className="w-full pl-10 pr-4 py-3 bg-night-surfaceAlt border border-night-border rounded-lg text-night-text placeholder-night-muted focus:outline-none focus:ring-2 focus:ring-night-blue"
                                    placeholder="••••••••"
                                    disabled={isLoading}
                                    autoComplete="current-password"
                                />
                            </div>
                        </div>

                        {/* Submit Button */}
                        <button
                            type="submit"
                            disabled={isLoading}
                            className="w-full py-3 bg-night-blue hover:bg-night-blue/80 disabled:bg-night-surfaceAlt text-night-text font-semibold rounded-lg transition-colors"
                        >
                            {isLoading ? 'Logging in...' : 'Login'}
                        </button>
                    </form>

                    {/* Register Link */}
                    <div className="mt-6 text-center">
                        <p className="text-night-muted text-sm">
                            Don't have an account?{' '}
                            <Link
                                to="/register"
                                className="text-night-blue hover:text-night-blue/80 font-medium"
                            >
                                Register
                            </Link>
                        </p>
                    </div>
                </div>

                {/* Demo Account Info (Optional) */}
                <div className="mt-6 text-center">
                    <p className="text-night-muted text-xs">
                        Demo: demo@example.com / password123
                    </p>
                </div>
            </div>
        </div>
    );
};
