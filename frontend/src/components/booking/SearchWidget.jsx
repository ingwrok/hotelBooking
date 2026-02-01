import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { DateRange } from 'react-date-range';
import { format, addDays } from 'date-fns';
import { Calendar as CalendarIcon, User, Search } from 'lucide-react';
import { updateBookingDetails } from '../../features/bookingSlice.js';
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';

const SearchWidget = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();
    const [openOptions, setOpenOptions] = useState(false);

    const bookingFromRedux = useSelector((state) => state.booking.currentBooking);

    const [dates, setDates] = useState([
        {
            startDate: bookingFromRedux.checkIn ? new Date(bookingFromRedux.checkIn) : new Date(),
            endDate: bookingFromRedux.checkOut ? new Date(bookingFromRedux.checkOut) : addDays(new Date(), 2),
            key: 'selection'
        }
    ]);
    const [openDate, setOpenDate] = useState(false);

    // Guests State
    const [adults, setAdults] = useState(bookingFromRedux.adults || 2); 
    const [children, setChildren] = useState(bookingFromRedux.children || 0);
    const [promoCode, setPromoCode] = useState(bookingFromRedux.promoCode || '');

    const handleSearch = () => {
        const checkIn = dates[0].startDate;
        const checkOut = dates[0].endDate;
        const duration = Math.ceil((checkOut - checkIn) / (1000 * 60 * 60 * 24));

        if (duration < 1) {
            alert("Minimum stay is 1 night");
            return;
        }

        dispatch(updateBookingDetails({
            checkIn: format(checkIn, 'yyyy-MM-dd'),
            checkOut: format(checkOut, 'yyyy-MM-dd'),
            adults,
            children,
            roomCount: 1, // Default
            promoCode
        }));
        navigate('/search');
    };

    return (
        <div className="w-full bg-white shadow-sm border-b border-gray-100 py-6">
            <div className="container mx-auto px-4">
                <div className="flex flex-col lg:flex-row items-center justify-center gap-6 text-sm">

                    {/* Date Picker Trigger */}
                    <div className="relative group cursor-pointer" onClick={() => setOpenDate(!openDate)}>
                        <div className="flex items-center gap-2">
                            <div className="text-center">
                                <span className="block text-xs text-gray-400 uppercase tracking-widest font-bold mb-1">Travel Date</span>
                                <div className="flex items-baseline gap-2">
                                    <span className="text-3xl font-serif text-primary">{format(dates[0].startDate, 'd')}</span>
                                    <span className="text-xs uppercase text-gray-500">{format(dates[0].startDate, 'MMM yyyy')}</span>
                                    <span className="mx-2 text-gray-300">|</span>
                                    <span className="text-3xl font-serif text-primary">{format(dates[0].endDate, 'd')}</span>
                                    <span className="text-xs uppercase text-gray-500">{format(dates[0].endDate, 'MMM yyyy')}</span>
                                </div>
                            </div>
                            <CalendarIcon className="w-5 h-5 text-accent ml-2 group-hover:text-primary transition-colors" />
                        </div>

                        {/* Date Range Picker Popup */}
                        {openDate && (
                            <div className="absolute top-full mt-4 left-1/2 -translate-x-1/2 z-50 shadow-2xl rounded-sm overflow-hidden" onClick={(e) => e.stopPropagation()}>
                                <DateRange
                                    editableDateInputs={true}
                                    onChange={item => setDates([item.selection])}
                                    moveRangeOnFirstSelection={false}
                                    ranges={dates}
                                    minDate={new Date()}
                                    months={2} // Show 2 months
                                    direction="horizontal"
                                    className="border border-gray-100"
                                    rangeColors={['#222222']} // Primary Color
                                />
                                <div className="bg-white p-3 text-right border-t border-gray-100">
                                    <button
                                        onClick={(e) => { e.stopPropagation(); setOpenDate(false); }}
                                        className="text-xs font-bold uppercase tracking-wider text-primary hover:text-accent"
                                    >
                                        Close
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>

                    <div className="h-10 w-px bg-gray-200 hidden lg:block mx-4"></div>

                    {/* Guests */}
                    <div className="relative cursor-pointer" onClick={() => setOpenOptions(!openOptions)}>
                        <span className="text-xs text-gray-400 uppercase tracking-widest font-bold mb-1 block text-center">Guests</span>
                        <div className="flex items-center gap-4">
                            <span className="text-xl font-serif text-primary">{adults + children} Persons</span>
                        </div>

                        {openOptions && (
                            <div className="absolute top-full mt-4 right-0 z-50 bg-white shadow-xl p-6 border border-gray-100 min-w-[200px]" onClick={e => e.stopPropagation()}>
                                <div className="flex justify-between items-center mb-4">
                                    <span className="text-sm font-bold">Adults</span>
                                    <div className="flex items-center gap-3">
                                        <button onClick={() => setAdults(Math.max(1, adults - 1))} className="border px-2">-</button>
                                        <span>{adults}</span>
                                        <button onClick={() => setAdults(adults + 1)} className="border px-2">+</button>
                                    </div>
                                </div>
                                <div className="flex justify-between items-center">
                                    <span className="text-sm font-bold">Children</span>
                                    <div className="flex items-center gap-3">
                                        <button onClick={() => setChildren(Math.max(0, children - 1))} className="border px-2">-</button>
                                        <span>{children}</span>
                                        <button onClick={() => setChildren(children + 1)} className="border px-2">+</button>
                                    </div>
                                </div>
                                <button
                                    onClick={() => setOpenOptions(false)}
                                    className="w-full mt-4 text-xs font-bold uppercase text-primary border-t pt-2"
                                >Done</button>
                            </div>
                        )}
                    </div>

                    <div className="h-10 w-px bg-gray-200 hidden lg:block mx-4"></div>

                    {/* Promo Code */}
                    <div className="flex flex-col items-start w-32">
                        <label className="text-xs text-gray-400 uppercase tracking-widest font-bold mb-1 w-full text-center">Code</label>
                        <input
                            type="text"
                            placeholder="Promo Code"
                            value={promoCode}
                            onChange={(e) => setPromoCode(e.target.value)}
                            className="w-full text-center border-b border-gray-300 pb-1 focus:outline-none focus:border-accent text-primary placeholder-gray-300 font-serif"
                        />
                    </div>

                    {/* Search Button */}
                    <button
                        onClick={handleSearch}
                        className="bg-primary hover:bg-gray-800 text-white px-8 py-3 uppercase text-xs tracking-[0.2em] font-bold transition-all ml-4"
                    >
                        Search
                    </button>
                </div>
            </div>

            {/* Backdrop for handling click outside */}
            {openDate && (
                <div className="fixed inset-0 z-40 bg-black/5" onClick={() => setOpenDate(false)}></div>
            )}
        </div>
    );
};

export default SearchWidget;
