import React from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Home, History, Dumbbell, Wrench, TrendingUp } from 'lucide-react';

export const Layout: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const navItems = [
        { path: '/dashboard', icon: Home, label: 'Home' },
        { path: '/history', icon: History, label: 'History' },
        { path: '/exercises', icon: Dumbbell, label: 'Exercises' },
        { path: '/equipment', icon: Wrench, label: 'Equipment' },
        { path: '/progress', icon: TrendingUp, label: 'Progress' },
    ];

    return (
        <div className="min-h-screen bg-night-bg flex flex-col">
            {/* Top Navigation */}
            <nav className="fixed top-0 left-0 right-0 bg-night-surface border-b border-night-border z-50">
                <div className="flex justify-around items-center p-3">
                    {navItems.map((item) => {
                        const Icon = item.icon;
                        const isActive = location.pathname.startsWith(item.path);
                        return (
                            <button
                                key={item.path}
                                onClick={() => navigate(item.path)}
                                className={`flex flex-col items-center gap-1 px-4 py-1 rounded-lg transition-colors ${
                                    isActive ? 'text-night-blue' : 'text-night-muted hover:text-night-text'
                                }`}
                            >
                                <Icon size={24} />
                                <span className="text-[10px] font-medium">{item.label}</span>
                            </button>
                        );
                    })}
                </div>
            </nav>
            
            <main className="flex-1 pt-16">
                <Outlet />
            </main>
        </div>
    );
};
