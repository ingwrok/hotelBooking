import React, { useEffect, useState } from 'react';
import api from '../../api';
import { Plus, Trash, Settings } from 'lucide-react';

const AdminRoomsPage = () => {
    const [rooms, setRooms] = useState([]);
    const [roomTypes, setRoomTypes] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isAdding, setIsAdding] = useState(false);

    // New Room Form State
    const [newRoomNumber, setNewRoomNumber] = useState('');
    const [newRoomTypeId, setNewRoomTypeId] = useState('');

    const fetchData = async () => {
        try {
            setLoading(true);
            const [roomsRes, typesRes] = await Promise.all([
                api.get('/rooms'),
                api.get('/room_types')
            ]);
            setRooms(roomsRes.data);
            setRoomTypes(typesRes.data);
            if (typesRes.data.length > 0) setNewRoomTypeId(typesRes.data[0].roomTypeId);
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const handleAddRoom = async (e) => {
        e.preventDefault();
        try {
            await api.post(`/rooms/${newRoomTypeId}`, { roomNumber: newRoomNumber });
            setIsAdding(false);
            setNewRoomNumber('');
            fetchData();
        } catch (e) {
            alert('Failed to add room. Room number might be duplicate.');
        }
    };

    const handleDelete = async (id) => {
        if (!window.confirm('Delete this room?')) return;
        try {
            await api.delete(`/rooms/${id}`);
            fetchData();
        } catch (e) {
            alert('Failed to delete room.');
        }
    };

    const handleStatusChange = async (id, currentStatus) => {
        const newStatus = currentStatus === 'available' ? 'maintenance' : 'available'; // Toggle for simplicity or use modal
        try {
            await api.patch(`/rooms/${id}/status`, { status: newStatus });
            fetchData();
        } catch (e) {
            alert('Failed to update status');
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div>
            <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4 mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Room Management</h2>
                <button onClick={() => setIsAdding(!isAdding)} className="bg-primary hover:bg-gray-800 text-white px-4 py-2 rounded flex items-center gap-2 transition-colors">
                    <Plus className="w-4 h-4" /> Add Room
                </button>
            </div>

            {isAdding && (
                <div className="bg-gray-50 p-4 rounded-lg mb-6 border border-gray-200">
                    <h4 className="font-bold mb-3">Add New Room</h4>
                    <form onSubmit={handleAddRoom} className="flex flex-col md:flex-row gap-4 md:items-end">
                        <div>
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Room Number</label>
                            <input
                                type="text"
                                value={newRoomNumber}
                                onChange={(e) => setNewRoomNumber(e.target.value)}
                                className="border p-2 rounded w-full md:w-40"
                                placeholder="e.g. 101"
                                required
                            />
                        </div>
                        <div className="w-full md:w-auto">
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Type</label>
                            <select
                                value={newRoomTypeId}
                                onChange={(e) => setNewRoomTypeId(e.target.value)}
                                className="border p-2 rounded w-full md:w-60 bg-white"
                            >
                                {roomTypes.map(rt => (
                                    <option key={rt.roomTypeId} value={rt.roomTypeId}>{rt.name}</option>
                                ))}
                            </select>
                        </div>
                        <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded font-bold hover:bg-green-700">Save</button>
                    </form>
                </div>
            )}

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-sm text-left">
                        <thead className="bg-gray-50 text-gray-500 font-bold uppercase text-xs">
                            <tr>
                                <th className="px-6 py-4">Room Number</th>
                                <th className="px-6 py-4">Type</th>
                                <th className="px-6 py-4 text-center">Status</th>
                                <th className="px-6 py-4 text-center">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {rooms.map(r => (
                                <tr key={r.roomId} className="hover:bg-gray-50 transition-colors">
                                    <td className="px-6 py-4 font-serif font-bold text-lg bg-gray-50/50">{r.roomNumber}</td>
                                    <td className="px-6 py-4">
                                        {roomTypes.find(rt => rt.roomTypeId === r.roomTypeId)?.name || r.roomTypeId}
                                    </td>
                                    <td className="px-6 py-4 text-center">
                                        <button
                                            onClick={() => handleStatusChange(r.roomId, r.status)}
                                            className={`px-3 py-1 rounded-full text-[10px] font-bold uppercase cursor-pointer hover:opacity-80
                                        ${r.status === 'available' ? 'bg-green-100 text-green-700' :
                                                    r.status === 'dirty' ? 'bg-red-100 text-red-700' :
                                                        r.status === 'maintenance' ? 'bg-gray-800 text-white' : 'bg-gray-100'}
                                    `}>
                                            {r.status}
                                        </button>
                                    </td>
                                    <td className="px-6 py-4 text-center">
                                        <button
                                            onClick={() => handleDelete(r.roomId)}
                                            className="text-red-400 hover:text-red-600 p-2"
                                        >
                                            <Trash className="w-4 h-4" />
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default AdminRoomsPage;
