import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { authApi } from '../../api/auth';
import { Dumbbell, Mail, Lock, User, AlertCircle, CheckCircle } from 'lucide-react';

export const Register: React.FC = () => {
    const navigate = useNavigate();
    const { setAuth } = useAuthStore();

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        // Validation
        if (!name || !email || !password || !confirmPassword) {
            setError('Please fill in all fields');
            return;
        }

        if (password.length < 8) {
            setError('Password must be at least 8 characters');
            return;
        }

        if (password !== confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        setIsLoading(true);

        try {
            console.log('📝 Registering user...');
            const response = await authApi.register({ name, email, password });
            console.log('✅ Registration successful:', response);
            console.log('Token received:', response.token);
            localStorage.setItem('auth_token', response.token);
            console.log('Token saved to localStorage');

            const savedToken = localStorage.getItem('auth_token');
            console.log('Verification - token in localStorage:', savedToken ? 'YES ✅' : 'NO ❌');
            setAuth(
                {
                    user_id: response.user_id,
                    email: response.email,
                    name: response.name,
                },
                response.token
            );
            console.log('Auth store updated');
            console.log('Navigating to dashboard...');
            navigate('/dashboard');
        } catch (err: any) {
            console.error('Registration error:', err);
            setError(
                err.response?.data?.message || 'Registration failed. Please try again.'
            );
        } finally {
            setIsLoading(false);
        }
    };

    const passwordStrength = password.length >= 8 ? 'strong' : password.length >= 6 ? 'medium' : 'weak';

    return (
        <div className="min-h-screen bg-night-bg flex items-center justify-center p-4">
            <div className="w-full max-w-md">
                {/* Logo/Header */}
                <div className="text-center mb-8">
                    <div className="flex items-center justify-center gap-3 mb-4">
                        <Dumbbell className="text-night-blue" size={40} />
                        <h1 className="text-3xl font-bold text-night-text">Gym Tracker</h1>
                    </div>
                    <p className="text-night-muted">Create your account to get started</p>
                </div>

                {/* Register Card */}
                <div className="bg-night-surface rounded-lg p-8 shadow-xl">
                    <h2 className="text-2xl font-semibold text-night-text mb-6">Register</h2>

                    {/* Error Message */}
                    {error && (
                        <div className="mb-4 p-4 bg-red-900 bg-opacity-50 border border-red-500 rounded-lg flex items-center gap-3">
                            <AlertCircle className="text-red-400" size={20} />
                            <p className="text-red-200 text-sm">{error}</p>
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        {/* Name Input */}
                        <div>
                            <label className="block text-sm font-medium text-night-text mb-2">
                                Name
                            </label>
                            <div className="relative">
                                <User
                                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-night-muted"
                                    size={20}
                                />
                                <input
                                    type="text"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    className="w-full pl-10 pr-4 py-3 bg-night-surfaceAlt border border-night-border rounded-lg text-night-text placeholder-night-muted focus:outline-none focus:ring-2 focus:ring-night-blue"
                                    placeholder="John Doe"
                                    disabled={isLoading}
                                    autoComplete="name"
                                />
                            </div>
                        </div>

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
                                    autoComplete="new-password"
                                />
                            </div>
                            {password && (
                                <div className="mt-2 flex items-center gap-2 text-xs">
                                    {passwordStrength === 'strong' && (
                                        <>
                                            <CheckCircle className="text-green-400" size={16} />
                                            <span className="text-green-400">Strong password</span>
                                        </>
                                    )}
                                    {passwordStrength === 'medium' && (
                                        <span className="text-yellow-400">Medium strength</span>
                                    )}
                                    {passwordStrength === 'weak' && (
                                        <span className="text-red-400">Weak password (min 8 chars)</span>
                                    )}
                                </div>
                            )}
                        </div>

                        {/* Confirm Password Input */}
                        <div>
                            <label className="block text-sm font-medium text-night-text mb-2">
                                Confirm Password
                            </label>
                            <div className="relative">
                                <Lock
                                    className="absolute left-3 top-1/2 transform -translate-y-1/2 text-night-muted"
                                    size={20}
                                />
                                <input
                                    type="password"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    className="w-full pl-10 pr-4 py-3 bg-night-surfaceAlt border border-night-border rounded-lg text-night-text placeholder-night-muted focus:outline-none focus:ring-2 focus:ring-night-blue"
                                    placeholder="••••••••"
                                    disabled={isLoading}
                                    autoComplete="new-password"
                                />
                            </div>
                            {confirmPassword && password !== confirmPassword && (
                                <p className="mt-2 text-xs text-red-400">Passwords do not match</p>
                            )}
                        </div>

                        {/* Submit Button */}
                        <button
                            type="submit"
                            disabled={isLoading || password !== confirmPassword}
                            className="w-full py-3 bg-night-blue hover:bg-night-blue/80 disabled:bg-night-surfaceAlt text-night-text font-semibold rounded-lg transition-colors"
                        >
                            {isLoading ? 'Creating account...' : 'Register'}
                        </button>
                    </form>

                    {/* Login Link */}
                    <div className="mt-6 text-center">
                        <p className="text-night-muted text-sm">
                            Already have an account?{' '}
                            <Link
                                to="/login"
                                className="text-night-blue hover:text-night-blue/80 font-medium"
                            >
                                Login
                            </Link>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};
