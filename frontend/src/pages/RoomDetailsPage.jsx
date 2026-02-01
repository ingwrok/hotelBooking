import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { Check, User, Wifi, Wind, Coffee } from 'lucide-react';
import { getRoomTypeDetails, getRatePlansByRoomType } from '../api/index.js';
import Button from '../components/common/Button.jsx';
import { updateBookingDetails } from '../features/bookingSlice.js';

const RoomDetailsPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const [roomType, setRoomType] = useState(null);
    const [ratePlans, setRatePlans] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            if (!id) return;
            try {
                setLoading(true);
                const [roomData, ratesData] = await Promise.all([
                    getRoomTypeDetails(id),
                    getRatePlansByRoomType(id)
                ]);
                setRoomType(roomData);
                setRatePlans(ratesData);
            } catch (error) {
                console.error("Failed to fetch details", error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    const handleReserve = (ratePlan) => {
        if (!roomType || !id) return;

        dispatch(updateBookingDetails({
            roomId: parseInt(id),
        }));

        navigate('/booking', { state: { roomType, ratePlan } });
    };

    if (loading) {
        return <div className="flex justify-center h-screen items-center"><div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div></div>;
    }

    if (!roomType) {
        return <div className="p-8 text-center">Room not found</div>;
    }

    return (
        <div className="min-h-screen bg-gray-50 pb-20">
            {/* Header Image */}
            <div className="h-[400px] w-full relative">
                <img
                    src={roomType.pictureUrl?.[0] || 'https://images.unsplash.com/photo-1611892440504-42a792e24d32?q=80&w=2070&auto=format&fit=crop'}
                    alt={roomType.name}
                    className="w-full h-full object-cover"
                />
                <div className="absolute inset-0 bg-black bg-opacity-30 flex items-end">
                    <div className="container mx-auto px-4 pb-8 text-white">
                        <h1 className="text-4xl font-serif font-bold mb-2">{roomType.name}</h1>
                        <p className="opacity-90 max-w-2xl">{roomType.description}</p>
                    </div>
                </div>
            </div>

            <div className="container mx-auto px-4 -mt-10 relative z-10">
                <div className="bg-white rounded-lg shadow-lg p-6 mb-8">
                    <h2 className="text-2xl font-serif mb-4">Room Amenities</h2>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div className="flex items-center text-gray-600"><Wifi className="w-5 h-5 mr-2" /> Free Wi-Fi</div>
                        <div className="flex items-center text-gray-600"><Wind className="w-5 h-5 mr-2" /> Air Conditioning</div>
                        <div className="flex items-center text-gray-600"><User className="w-5 h-5 mr-2" /> Max {roomType.capacity} Guests</div>
                        <div className="flex items-center text-gray-600"><Coffee className="w-5 h-5 mr-2" /> Breakfast Included</div>
                    </div>
                </div>

                <h2 className="text-2xl font-serif mb-6 text-gray-800">Available Rates</h2>
                <div className="space-y-4">
                    {ratePlans.map((plan) => (
                        <div key={plan.ratePlanId} className="bg-white border rounded-lg p-6 flex flex-col md:flex-row justify-between items-center shadow-sm hover:shadow-md transition">
                            <div className="mb-4 md:mb-0">
                                <h3 className="text-xl font-bold text-gray-800">{plan.name}</h3>
                                <div className="text-sm text-gray-500 mt-1 flex flex-col space-y-1">
                                    <span className="flex items-center">{plan.allowFreeCancel ? <Check className="text-green-500 w-4 h-4 mr-1" /> : null} {plan.allowFreeCancel ? "Free Cancellation" : "Non-Refundable"}</span>
                                    <span className="flex items-center">{plan.allowPayLater ? <Check className="text-green-500 w-4 h-4 mr-1" /> : null} {plan.allowPayLater ? "Pay Later Available" : "Pay Now"}</span>
                                    <span className="flex items-center"><Check className="text-green-500 w-4 h-4 mr-1" /> Instant Confirmation</span>
                                </div>
                            </div>

                            <div className="text-right">
                                <div className="text-3xl font-bold text-gray-900 mb-2">
                                    THB {plan.price ? plan.price.toLocaleString() : '5,000'}
                                    <span className="text-xs font-normal text-gray-500 block">per night</span>
                                </div>
                                <Button onClick={() => handleReserve(plan)} size="lg" className="px-8 min-w-[150px]">
                                    RESERVE
                                </Button>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default RoomDetailsPage;
