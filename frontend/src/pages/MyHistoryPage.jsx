import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import api from '../api';
import { Bed, Calendar } from 'lucide-react';

const MyHistoryPage = () => {
    const navigate = useNavigate();
    const user = useSelector((state) => state.auth.user);
    const [bookings, setBookings] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!user) {
            navigate('/login');
            return;
        }

        const fetchBookings = async () => {
            try {
                const token = user.token;
                const config = {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                };
                const response = await api.get('/bookings/my', config);
                setBookings(response.data);
            } catch (error) {
                console.error("Error fetching bookings:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchBookings();
    }, [user, navigate]);

    if (loading) {
        return <div className="min-h-screen text-center pt-20">Loading history...</div>;
    }

    return (
        <div className="min-h-screen bg-gray-50 py-12">
            <div className="container mx-auto px-4 max-w-4xl">
                <h1 className="text-3xl font-serif font-bold text-gray-800 mb-8">My Booking History</h1>

                {bookings.length === 0 ? (
                    <div className="bg-white p-8 rounded-lg shadow text-center">
                        <p className="text-gray-500 mb-4">You haven't made any bookings yet.</p>
                    </div>
                ) : (
                    <div className="space-y-6">
                        {bookings.map((booking) => (
                            <div key={booking.bookingId} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
                                <div className="p-6">
                                    <div className="flex flex-col md:flex-row justify-between md:items-center mb-4 border-b border-gray-100 pb-4">
                                        <div>
                                            <span className={`inline-block px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider mb-2 ${booking.status === 'confirmed' ? 'bg-green-100 text-green-800' :
                                                booking.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                                    'bg-gray-100 text-gray-800'
                                                }`}>
                                                {booking.status}
                                            </span>
                                            <h3 className="text-xl font-bold text-primary">Booking #{booking.bookingId}</h3>
                                        </div>
                                        <div className="text-right mt-2 md:mt-0 flex flex-col items-end">
                                            <p className="text-sm text-gray-500">Total Price</p>
                                            <p className="text-xl font-serif font-bold mb-2">THB {booking.totalPrice.toLocaleString()}</p>

                                            {booking.status === 'pending' && (
                                                <button
                                                    onClick={() => navigate('/confirmation', { state: { booking: booking } })}
                                                    className="text-xs bg-primary text-white px-3 py-1 rounded hover:bg-primary-dark transition-colors"
                                                >
                                                    Pay Now
                                                </button>
                                            )}
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        <div className="flex items-center text-gray-600">
                                            <Bed className="w-5 h-5 mr-3 text-accent" />
                                            <div>
                                                <p className="text-xs text-gray-400 uppercase">Room Type</p>
                                                <p className="font-semibold">{booking.roomTypeName} ({booking.roomNumber})</p>
                                            </div>
                                        </div>
                                        <div className="flex items-center text-gray-600">
                                            <Calendar className="w-5 h-5 mr-3 text-accent" />
                                            <div>
                                                <p className="text-xs text-gray-400 uppercase">Dates</p>
                                                <p className="font-semibold">
                                                    {new Date(booking.checkInDate).toLocaleDateString()} - {new Date(booking.checkOutDate).toLocaleDateString()}
                                                </p>
                                            </div>
                                        </div>
                                    </div>

                                    {booking.bookingAddon && booking.bookingAddon.length > 0 && (
                                        <div className="mt-4 pt-4 border-t border-gray-50">
                                            <p className="text-xs text-gray-400 uppercase mb-2">Add-ons</p>
                                            <div className="flex flex-wrap gap-2">
                                                {booking.bookingAddon.map((addon) => (
                                                    <span key={addon.bookingAddonId} className="bg-gray-50 text-gray-600 px-2 py-1 rounded text-xs border border-gray-200">
                                                        {addon.addonName} (x{addon.quantity})
                                                    </span>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default MyHistoryPage;
