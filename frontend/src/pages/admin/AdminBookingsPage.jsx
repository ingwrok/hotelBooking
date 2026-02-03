import React, { useEffect, useState } from 'react';
import api from '../../api';
import { format } from 'date-fns';
import { Check, X, Clock, Eye } from 'lucide-react';

const AdminBookingsPage = () => {
    const [bookings, setBookings] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchBookings = async () => {
        try {
            setLoading(true);
            const res = await api.get('/bookings/all');
            setBookings(res.data);
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchBookings();
    }, []);

    const handleStatusChange = async (id, status) => {
        if (!window.confirm(`Are you sure you want to change status to ${status}?`)) return;
        try {
            await api.patch(`/bookings/${id}/status`, { status });
            fetchBookings(); // Refresh
        } catch (e) {
            alert('Failed to update status');
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div>
            <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4 mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Booking Management</h2>
                <div className="bg-white px-4 py-2 border rounded shadow-sm text-sm">
                    Total: <span className="font-bold">{bookings.length}</span>
                </div>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-sm text-left">
                        <thead className="bg-gray-50 text-gray-500 font-bold uppercase text-xs">
                            <tr>
                                <th className="px-6 py-4">ID</th>
                                <th className="px-6 py-4">Guest</th>
                                <th className="px-6 py-4">Pax</th>
                                <th className="px-6 py-4">Dates</th>
                                <th className="px-6 py-4">Room</th>
                                <th className="px-6 py-4 text-right">Total</th>
                                <th className="px-6 py-4 text-center">Status</th>
                                <th className="px-6 py-4 text-center">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {bookings.map(b => (
                                <tr key={b.bookingId} className="hover:bg-gray-50 transition-colors">
                                    <td className="px-6 py-4 text-gray-400 font-mono">#{b.bookingId}</td>
                                    <td className="px-6 py-4">
                                        <div className="font-bold text-gray-800">{b.guestDetails?.firstName} {b.guestDetails?.lastName}</div>
                                        <div className="text-xs text-gray-400">{b.email || b.userEmail}</div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-1 font-bold text-gray-700">
                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
                                            {b.numAdults}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div>{format(new Date(b.checkInDate), 'dd MMM yyyy')}</div>
                                        <div className="text-xs text-gray-400">to {format(new Date(b.checkOutDate), 'dd MMM yyyy')}</div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="font-bold">{b.roomTypeName}</div>
                                        <div className="text-xs text-gray-400">Room {b.roomNumber || '-'}</div>
                                    </td>
                                    <td className="px-6 py-4 text-right font-serif font-bold">
                                        {b.totalPrice.toLocaleString()}
                                    </td>
                                    <td className="px-6 py-4 text-center">
                                        <span className={`px-3 py-1 rounded-full text-[10px] font-bold uppercase
                                        ${b.status === 'confirmed' ? 'bg-green-100 text-green-700' :
                                                b.status === 'pending' ? 'bg-yellow-100 text-yellow-700' :
                                                    b.status === 'cancelled' ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-700'}
                                    `}>
                                            {b.status}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex justify-center gap-2">
                                            {b.status === 'pending' && (
                                                <>
                                                    <button
                                                        onClick={() => handleStatusChange(b.bookingId, 'confirmed')}
                                                        className="p-1 hover:bg-green-100 text-green-600 rounded" title="Confirm"
                                                    >
                                                        <Check className="w-4 h-4" />
                                                    </button>
                                                    <button
                                                        onClick={() => handleStatusChange(b.bookingId, 'cancelled')}
                                                        className="p-1 hover:bg-red-100 text-red-600 rounded" title="Cancel"
                                                    >
                                                        <X className="w-4 h-4" />
                                                    </button>
                                                </>
                                            )}
                                            {b.status === 'confirmed' && (
                                                <button
                                                    onClick={() => handleStatusChange(b.bookingId, 'checked-in')}
                                                    className="p-1 hover:bg-blue-100 text-blue-600 rounded" title="Check In"
                                                >
                                                    <Clock className="w-4 h-4" />
                                                </button>
                                            )}
                                        </div>
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

export default AdminBookingsPage;
