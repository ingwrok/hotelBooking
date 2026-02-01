import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { fetchRooms } from '../features/roomSlice.js';
import { getRatePlansByRoomType } from '../api/index.js';
import { updateBookingDetails } from '../features/bookingSlice.js';
import { Check, User, Info, Wifi, Coffee, Star } from 'lucide-react';

import SearchWidget from '../components/booking/SearchWidget.jsx';
import ImageCarousel from '../components/common/ImageCarousel.jsx';

const SearchPage = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const bookingState = useSelector((state) => state.booking.currentBooking);
    const { rooms, loading } = useSelector((state) => state.rooms);

    const [selectedRoomId, setSelectedRoomId] = useState(null);
    const [ratePlans, setRatePlans] = useState([]);
    const [loadingRates, setLoadingRates] = useState(false);
    const [cols, setCols] = useState(4); // Default to 4 for desktop

    // Responsive columns listener
    useEffect(() => {
        const handleResize = () => {
            if (window.innerWidth < 768) setCols(1);      // Mobile
            else if (window.innerWidth < 1024) setCols(3); // Tablet (md)
            else setCols(4);                               // Desktop (lg)
        };

        handleResize(); // Init
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    useEffect(() => {
        // Dispatch with current selected dates from booking state
        if (bookingState.checkIn && bookingState.checkOut) {
            // Format dates to YYYY-MM-DD for backend
            const formatDate = (date) => {
                if (!date) return "";
                if (typeof date === 'string') return date; // Already string?
                // Date object
                const year = date.getFullYear();
                const month = String(date.getMonth() + 1).padStart(2, '0');
                const day = String(date.getDate()).padStart(2, '0');
                return `${year}-${month}-${day}`;
            };

            dispatch(fetchRooms({
                checkIn: formatDate(bookingState.checkIn),
                checkOut: formatDate(bookingState.checkOut)
            }));
        } else {
            dispatch(fetchRooms());
        }
    }, [dispatch, bookingState.checkIn, bookingState.checkOut]);

    // Fetch rates when selected room changes
    useEffect(() => {
        const fetchRates = async () => {
            if (selectedRoomId) {
                setLoadingRates(true);
                try {
                    const rates = await getRatePlansByRoomType(selectedRoomId);
                    setRatePlans(rates);
                } catch (e) {
                    console.error("Failed to load rates", e);
                } finally {
                    setLoadingRates(false);
                }
            }
        };
        fetchRates();
    }, [selectedRoomId]);

    const handleBookRate = (room, rate) => {
        dispatch(updateBookingDetails({
            roomId: room.roomTypeId,
            roomName: room.name,
            ratePlanId: rate.ratePlanId,
            ratePlanName: rate.name,
            pricePerNight: rate.price,
            roomImage: room.pictureUrl && room.pictureUrl.length > 0 ? room.pictureUrl[0] : ''
        }));
        navigate('/booking');
    };

    const handleRoomClick = (id) => {
        if (selectedRoomId === id) {
            setSelectedRoomId(null); // Toggle off
        } else {
            setSelectedRoomId(id); // Toggle on
        }
    };

    // Helper to chunk rooms into rows
    const chunkedRooms = [];
    for (let i = 0; i < rooms.length; i += cols) {
        chunkedRooms.push(rooms.slice(i, i + cols));
    }

    const selectedRoom = rooms.find(r => r.roomTypeId === selectedRoomId);

    return (
        <div className="min-h-screen bg-white pb-20 font-sans">
            {/* Search Widget - Compact Mode could be passed here, reusing standard for now */}
            <SearchWidget />

            <div className="container mx-auto px-4 mt-8">

                {/* Room Grid Selection */}
                <div className="mb-12">
                    <div className="flex justify-between items-end mb-6 border-b border-gray-100 pb-2">
                        <div className="flex gap-4">
                            <button className="bg-primary text-white px-4 py-2 text-xs font-bold uppercase tracking-wider">Rooms</button>
                            <button className="text-gray-400 hover:text-primary px-4 py-2 text-xs font-bold uppercase tracking-wider">Rates</button>
                        </div>
                        <div className="flex items-center gap-2 text-xs text-gray-500">
                            <span>Sort by: </span>
                            <select className="border-none bg-transparent font-bold text-primary focus:outline-none">
                                <option>Recommended</option>
                                <option>Price: Low to High</option>
                            </select>
                        </div>
                    </div>

                    {loading ? (
                        <div className="h-64 flex items-center justify-center">Loading Rooms...</div>
                    ) : (
                        <div className="flex flex-col gap-6">
                            {chunkedRooms.map((rowRooms, rowIndex) => {
                                // Check if any room in this row is selected
                                const isRowActive = rowRooms.some(r => r.roomTypeId === selectedRoomId);

                                return (
                                    <React.Fragment key={rowIndex}>
                                        {/* Grid Row */}
                                        <div className={`grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6`}>
                                            {rowRooms.map(room => {
                                                const isSelected = selectedRoomId === room.roomTypeId;
                                                return (
                                                    <div
                                                        key={room.roomTypeId}
                                                        onClick={() => handleRoomClick(room.roomTypeId)}
                                                        className={`cursor-pointer transition-all duration-300 relative group bg-white
                                                            ${isSelected ? 'ring-2 ring-primary shadow-lg' : 'hover:opacity-90'}
                                                        `}
                                                    >
                                                        <div className="aspect-[4/3] overflow-hidden bg-gray-100 mb-3 relative">
                                                            <img
                                                                src={room.pictureUrl && room.pictureUrl.length > 0 ? room.pictureUrl[0] : 'https://via.placeholder.com/400x300?text=No+Image'}
                                                                alt={room.name}
                                                                className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
                                                            />
                                                            {room.availableCount === 1 && (
                                                                <div className="absolute top-2 left-2 bg-red-600 text-white text-[10px] font-bold uppercase px-2 py-1 rounded shadow-md animate-pulse">
                                                                    Only 1 Left!
                                                                </div>
                                                            )}
                                                        </div>
                                                        <div className="p-3 pt-0">
                                                            <h3 className="text-xs font-bold uppercase tracking-wider text-primary mb-1 truncate">{room.name}</h3>
                                                            <div className="flex justify-between items-end">
                                                                <div>
                                                                    <div className="text-[10px] text-gray-400">From</div>
                                                                    <div className="text-sm font-serif text-primary">THB {room.price ? room.price.toLocaleString() : 'Check Dates'}</div>
                                                                </div>
                                                                <button className={`text-[10px] font-bold uppercase px-3 py-2 transition-colors ${isSelected ? 'bg-primary text-white' : 'bg-gray-200 text-gray-600'}`}>
                                                                    {isSelected ? 'Viewing' : 'View Rates'}
                                                                </button>
                                                            </div>
                                                        </div>
                                                        {isSelected && (
                                                            <div className="absolute top-2 right-2 bg-primary text-white p-1 rounded-full">
                                                                <Check className="w-3 h-3" />
                                                            </div>
                                                        )}
                                                    </div>
                                                );
                                            })}
                                        </div>

                                        {/* Details Panel - Injected after the active row */}
                                        <div className={`transition-all duration-500 ease-in-out overflow-hidden ${isRowActive ? 'max-h-[2000px] opacity-100 mb-6' : 'max-h-0 opacity-0'}`}>
                                            {selectedRoom && (
                                                <div
                                                    key={selectedRoom.roomTypeId}
                                                    className="bg-white border border-gray-200 shadow-inner p-6 relative mt-2 min-h-[500px] animate-fade-in"
                                                >
                                                    <button onClick={() => setSelectedRoomId(null)} className="absolute top-4 right-4 text-gray-400 hover:text-black z-10">
                                                        X
                                                    </button>
                                                    <div className="flex items-center justify-between mb-6">
                                                        <h2 className="text-2xl font-serif text-primary uppercase">{selectedRoom.name}</h2>
                                                        <div className="flex gap-4 text-xs text-gray-500">
                                                            <span className="flex items-center gap-1"><User className="w-3 h-3" /> {selectedRoom.capacity} Adults</span>
                                                            <span className="flex items-center gap-1">|</span>
                                                            <span className="flex items-center gap-1">{selectedRoom.sizeSqm} m²</span>
                                                        </div>
                                                    </div>

                                                    <div className="flex flex-col lg:flex-row gap-8">
                                                        {/* Left: Carousel */}
                                                        <div className="w-full lg:w-5/12 aspect-video bg-gray-100">
                                                            <ImageCarousel images={selectedRoom.pictureUrl} alt={selectedRoom.name} />
                                                        </div>

                                                        {/* Right: Info & Rates */}
                                                        <div className="flex-1">
                                                            <div className="mb-6 border-b border-gray-100 pb-6">
                                                                <h4 className="font-bold text-sm uppercase mb-3 text-gray-700">Room Amenities</h4>
                                                                <div className="flex flex-wrap gap-2">
                                                                    {selectedRoom.amenities && selectedRoom.amenities.map(am => (
                                                                        <span key={am} className="text-[10px] border border-gray-200 px-2 py-1 text-gray-500 uppercase">{am}</span>
                                                                    ))}
                                                                </div>
                                                                <p className="text-sm text-gray-500 mt-4 leading-relaxed">{selectedRoom.description}</p>
                                                            </div>

                                                            <div className="space-y-3">
                                                                {loadingRates ? (
                                                                    <div className="h-40 flex items-center justify-center text-gray-400">Loading Rates...</div>
                                                                ) : ratePlans.length > 0 ? ratePlans.map(rate => (
                                                                    <div key={rate.ratePlanId} className="flex justify-between items-center group hover:bg-gray-50 p-3 rounded transition-colors border-b border-gray-50 last:border-0">
                                                                        <div>
                                                                            <div className="font-bold text-sm text-primary uppercase">{rate.name}</div>
                                                                            <div className="text-xs text-gray-500">{rate.description}</div>
                                                                            <div className="flex gap-3 mt-1 text-[10px]">
                                                                                {rate.allowFreeCancel && <span className="text-green-600 flex items-center gap-1"><Check className="w-2 h-2" /> Free Cancel</span>}
                                                                                {!rate.allowPayLater && <span className="text-gray-400 flex items-center gap-1"><Info className="w-2 h-2" /> Non-Refundable</span>}
                                                                            </div>
                                                                        </div>
                                                                        <div className="text-right flex items-center gap-4">
                                                                            <div>
                                                                                <div className="text-lg font-serif text-primary">THB {rate.price.toLocaleString()}</div>
                                                                                <div className="text-[10px] text-gray-400">per night</div>
                                                                            </div>
                                                                            <button
                                                                                onClick={() => handleBookRate(selectedRoom, rate)}
                                                                                className="bg-primary hover:bg-gray-800 text-white px-6 py-3 text-xs font-bold uppercase tracking-widest transition-all"
                                                                            >
                                                                                Reserve
                                                                            </button>
                                                                        </div>
                                                                    </div>
                                                                )) : (
                                                                    <div className="text-center text-gray-400 py-4">No rates available for these dates</div>
                                                                )}
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                    </React.Fragment>
                                );
                            })}
                        </div>
                    )}
                </div>

            </div>
        </div>
    );
};

export default SearchPage;

