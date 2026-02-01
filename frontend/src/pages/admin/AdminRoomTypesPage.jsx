import React, { useEffect, useState } from 'react';
import api from '../../api';
import { Link } from 'react-router-dom';
import { Plus, Edit, Trash, Image as ImageIcon } from 'lucide-react';

const AdminRoomTypesPage = () => {
    const [roomTypes, setRoomTypes] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchRoomTypes = async () => {
        try {
            setLoading(true);
            const res = await api.get('/room_types');
            setRoomTypes(res.data);
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRoomTypes();
    }, []);

    const handleDelete = async (id) => {
        if (!window.confirm('Are you sure you want to delete this room type?')) return;
        try {
            await api.delete(`/room_types/${id}`);
            fetchRoomTypes();
        } catch (e) {
            alert('Failed to delete room type. It might have existing bookings or rooms.');
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Room Types</h2>
                <Link to="new" className="bg-primary hover:bg-gray-800 text-white px-4 py-2 rounded flex items-center gap-2 transition-colors">
                    <Plus className="w-4 h-4" /> New Room Type
                </Link>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {roomTypes.map(rt => (
                    <div key={rt.roomTypeId} className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden group">
                        <div className="h-40 bg-gray-200 relative">
                            {rt.pictureUrl && rt.pictureUrl.length > 0 ? (
                                <img src={rt.pictureUrl[0]} alt={rt.name} className="w-full h-full object-cover" />
                            ) : (
                                <div className="flex items-center justify-center h-full text-gray-400">
                                    <ImageIcon className="w-8 h-8" />
                                </div>
                            )}
                            <div className="absolute top-2 right-2 bg-white/90 px-2 py-1 text-xs font-bold rounded">
                                ID: {rt.roomTypeId}
                            </div>
                        </div>
                        <div className="p-4">
                            <h3 className="font-bold text-lg text-primary mb-1">{rt.name}</h3>
                            <div className="text-sm text-gray-500 mb-4 line-clamp-2">{rt.description}</div>

                            <div className="flex justify-between items-center text-sm font-medium mb-4">
                                <div>
                                    <span className="text-gray-400">Price:</span> {rt.price ? rt.price.toLocaleString() : 'N/A'}
                                </div>
                                <div>
                                    <span className="text-gray-400">Capacity:</span> {rt.capacity}
                                </div>
                            </div>

                            <div className="flex justify-end gap-2 border-t pt-3">
                                <Link to={`edit/${rt.roomTypeId}`} className="p-2 hover:bg-gray-100 rounded text-blue-600 disabled:opacity-50">
                                    <Edit className="w-4 h-4" />
                                </Link>
                                <button onClick={() => handleDelete(rt.roomTypeId)} className="p-2 hover:bg-red-50 rounded text-red-500">
                                    <Trash className="w-4 h-4" />
                                </button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default AdminRoomTypesPage;
