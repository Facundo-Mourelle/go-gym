import React, { useState, useEffect } from 'react';
import { Wrench, Search, Plus, Pencil } from 'lucide-react';
import { equipmentApi } from '../api/equipment';
import type { EquipmentData } from '../api/equipment';
import { CreateEquipmentModal } from '../components/CreateEquipmentModal';
import { EditEquipmentModal } from '../components/EditEquipmentModal';

export const Equipment: React.FC = () => {
    const [equipment, setEquipment] = useState<EquipmentData[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [showCreate, setShowCreate] = useState(false);
    const [editingEquipment, setEditingEquipment] = useState<EquipmentData | null>(null);

    const fetchEquipment = async () => {
        setLoading(true);
        try {
            const data = await equipmentApi.list();
            setEquipment(data);
        } catch (err) {
            console.error("Failed to fetch equipment", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchEquipment();
    }, []);

    const filteredEquipment = equipment.filter(eq =>
        eq.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const freeWeights = filteredEquipment.filter(eq => eq.type === 'freeweight');
    const cables = filteredEquipment.filter(eq => eq.type === 'cable');
    const machines = filteredEquipment.filter(eq => eq.type === 'machine');

    const renderEquipmentCard = (eq: EquipmentData) => (
        <div key={eq.id} className="bg-gray-800 rounded-2xl p-4 border border-gray-700/50 hover:border-gray-600 transition-colors group">
            <div className="flex items-start justify-between gap-2">
                <div className="flex-1 min-w-0">
                    <h3 className="font-bold text-lg truncate">{eq.name}</h3>
                    {eq.manufacturer && (
                        <p className="text-gray-400 text-sm mt-0.5 truncate">{eq.manufacturer}</p>
                    )}
                </div>
                <button
                    onClick={() => setEditingEquipment(eq)}
                    className="p-1.5 text-gray-500 hover:text-blue-400 hover:bg-blue-400/10 rounded-lg transition-colors shrink-0"
                    title="Edit equipment"
                >
                    <Pencil size={16} />
                </button>
            </div>
            <div className="mt-3 flex flex-wrap gap-2">
                <span className={`text-xs px-2 py-1 rounded-md border ${
                    eq.type === 'freeweight' ? 'bg-orange-900/50 text-orange-300 border-orange-800/50' :
                    eq.type === 'cable' ? 'bg-purple-900/50 text-purple-300 border-purple-800/50' :
                    'bg-blue-900/50 text-blue-300 border-blue-800/50'
                }`}>
                    {eq.type === 'freeweight' ? 'Free Weight' : eq.type === 'cable' ? 'Cable' : 'Machine'}
                </span>
                {eq.type === 'freeweight' && eq.actual_weight && (
                    <span className="bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded-md border border-gray-600">
                        {eq.actual_weight} kg
                    </span>
                )}
                {eq.type === 'cable' && (
                    <>
                        {eq.pulley_type && (
                            <span className="bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded-md border border-gray-600">
                                {eq.pulley_type}
                            </span>
                        )}
                        {eq.stack_weights && eq.stack_weights.length > 0 && (
                            <span className="bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded-md border border-gray-600">
                                {eq.stack_weights.length} weights
                            </span>
                        )}
                    </>
                )}
                {eq.type === 'machine' && eq.resistance_profile_name && (
                    <span className="bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded-md border border-gray-600">
                        {eq.resistance_profile_name}
                    </span>
                )}
            </div>
        </div>
    );

    const renderSection = (title: string, items: EquipmentData[]) => {
        if (items.length === 0) return null;
        return (
            <div className="mb-6">
                <h2 className="text-lg font-bold text-gray-300 mb-3 flex items-center gap-2">
                    {title}
                    <span className="text-sm font-normal text-gray-500">({items.length})</span>
                </h2>
                <div className="grid gap-3">
                    {items.map(renderEquipmentCard)}
                </div>
            </div>
        );
    };

    return (
        <div className="bg-gray-900 h-screen flex flex-col text-white">
            <div className="p-4 bg-gray-800 border-b border-gray-700 shadow-md">
                <div className="flex items-center justify-between mb-4">
                    <h1 className="text-2xl font-bold flex items-center gap-2">
                        <Wrench className="text-blue-500" />
                        Equipment
                    </h1>
                    <button
                        onClick={() => setShowCreate(true)}
                        className="flex items-center gap-1.5 bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-xl font-semibold transition-colors text-sm"
                    >
                        <Plus size={18} />
                        <span>Create</span>
                    </button>
                </div>
                
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Search className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                        type="text"
                        className="bg-gray-900 border border-gray-700 text-white text-sm rounded-xl focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-3 outline-none transition-colors"
                        placeholder="Search equipment..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
            </div>

            <div className="flex-1 overflow-y-auto p-4">
                {loading ? (
                    <div className="flex justify-center p-10">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                    </div>
                ) : filteredEquipment.length === 0 ? (
                    <div className="text-center text-gray-400 mt-10 bg-gray-800 p-8 rounded-2xl border border-gray-700/50">
                        <Wrench className="mx-auto h-12 w-12 text-gray-600 mb-3" />
                        <p className="font-medium">No equipment found.</p>
                        <p className="text-sm mt-1">Try a different search term or create new equipment.</p>
                    </div>
                ) : (
                    <>
                        {renderSection('Free Weights', freeWeights)}
                        {renderSection('Cables', cables)}
                        {renderSection('Machines', machines)}
                    </>
                )}
            </div>

            {showCreate && (
                <CreateEquipmentModal
                    onClose={() => setShowCreate(false)}
                    onCreated={fetchEquipment}
                />
            )}

            {editingEquipment && (
                <EditEquipmentModal
                    equipment={editingEquipment}
                    onClose={() => setEditingEquipment(null)}
                    onUpdated={fetchEquipment}
                />
            )}
        </div>
    );
};