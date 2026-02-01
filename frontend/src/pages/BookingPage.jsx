import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { toggleAddon, updateGuestDetails, setPaymentMethod, submitBooking } from '../features/bookingSlice.js';
import { fetchAddonCategories, fetchAddonsByCategory } from '../features/addonSlice.js';
import { Check, Star, ShieldCheck, CreditCard as CardIcon, ChevronUp, ChevronDown } from 'lucide-react';

const steps = ['Enhance Stay', 'Guest Details', 'Payment'];

const BookingPage = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const bookingState = useSelector((state) => state.booking.currentBooking);
    // State for Addon Category
    const { addons, categories } = useSelector((state) => state.addons);
    const { user } = useSelector((state) => state.auth);
    const [activeCategory, setActiveCategory] = useState(null);
    const [showMobileSummary, setShowMobileSummary] = useState(false);

    // Initial State from navigation (room and rate)
    const { roomId, roomName, ratePlanId, ratePlanName, pricePerNight, roomImage } = bookingState;

    const [currentStep, setCurrentStep] = useState(0);

    // Redirect if not logged in & Auto-fill Guest Details
    useEffect(() => {
        if (!user) {
            navigate('/login');
        } else {
            // Auto-fill if empty
            const userData = user.user || user; // Handle potential structure variations
            if (!guestDetails.firstName && !guestDetails.email) {
                dispatch(updateGuestDetails({
                    firstName: userData.username,
                    email: userData.email
                }));
            }
        }
    }, [user, navigate, dispatch]);

    // Fetch Categories on Mount
    useEffect(() => {
        dispatch(fetchAddonCategories());
    }, [dispatch]);

    // Select first category by default when loaded
    useEffect(() => {
        if (categories.length > 0 && !activeCategory) {
            const firstId = categories[0].categoryId;
            setActiveCategory(firstId);
            dispatch(fetchAddonsByCategory(firstId));
        }
    }, [categories, activeCategory, dispatch]);

    const handleCategoryChange = (id) => {
        if (activeCategory === id) return;
        setActiveCategory(id);
        dispatch(fetchAddonsByCategory(id));
    };

    // Redux State
    const selectedAddons = bookingState.selectedAddons || {};
    const guestDetails = bookingState.guestDetails || {};
    const paymentMethod = bookingState.paymentMethod || 'qr';

    useEffect(() => {
        if (!roomId) {
            navigate('/');
            return;
        }
    }, [roomId, navigate]);


    // Calculations
    const roomTotal = (pricePerNight || 0) * (bookingState.roomCount || 1);
    // Assuming 1 night for simplicity unless dates are parsed
    const addonsTotal = addons.reduce((sum, addon) => {
        return sum + (selectedAddons[addon.addonId] ? addon.price : 0);
    }, 0);
    const grandTotal = roomTotal + addonsTotal;

    const handleNext = () => {
        if (currentStep < steps.length - 1) {
            setCurrentStep(prev => prev + 1);
        } else {
            handleBookingSubmit();
        }
    };

    const handleBack = () => {
        if (currentStep > 0) {
            setCurrentStep(prev => prev - 1);
        } else {
            navigate(-1);
        }
    };

    const handleBookingSubmit = async () => {
        const bookingPayload = {
            userId: user.user?.id || user.id, // Dynamic User ID
            roomTypeId: bookingState.roomId, // SearchPage sets roomId as roomTypeId
            ratePlanId: bookingState.ratePlanId,
            checkInDate: bookingState.checkIn,
            checkOutDate: bookingState.checkOut,
            numAdults: bookingState.adults,
            email: guestDetails.email,
            bookingAddon: Object.entries(selectedAddons)
                .filter(([_, qty]) => qty > 0)
                .map(([id, qty]) => ({ addonId: parseInt(id), quantity: qty }))
        };

        try {
            const resultAction = await dispatch(submitBooking(bookingPayload));
            if (submitBooking.fulfilled.match(resultAction)) {
                navigate('/confirmation', { state: { booking: resultAction.payload } });
            } else {
                alert("Booking failed: " + resultAction.error.message);
            }
        } catch (error) {
            alert("Booking failed. Please try again.");
        }
    };

    return (
        <div className="min-h-screen bg-white font-sans text-primary">
            {/* Header / Steps Indicator */}
            <div className="border-b border-gray-100 py-6">
                <div className="container mx-auto px-4 flex justify-between items-center">
                    <button onClick={handleBack} className="text-xs uppercase tracking-widest text-gray-400 hover:text-primary">
                        ← Back
                    </button>
                    <div className="flex gap-2">
                        {steps.map((s, idx) => (
                            <div key={s} className="flex items-center">
                                <span className={`w-6 h-6 flex items-center justify-center text-[10px] font-bold border ${idx === currentStep ? 'bg-primary text-white border-primary' : 'bg-white text-gray-300 border-gray-200'} rounded-sm mr-2 transition-colors`}>
                                    {idx + 1}
                                </span>
                                <span className={`text-xs uppercase tracking-wider ${idx === currentStep ? 'text-primary' : 'text-gray-300'} hidden md:block`}>{s}</span>
                                {idx < steps.length - 1 && <div className="w-8 h-px bg-gray-100 mx-3 hidden md:block"></div>}
                            </div>
                        ))}
                    </div>
                    <div className="text-xs text-gray-400">Step {currentStep + 1} of {steps.length}</div>
                </div>
            </div>

            <div className="container mx-auto px-4 py-12 lg:pb-12 pb-32">
                <div className="flex flex-col lg:flex-row gap-12">

                    {/* Main Content Area */}
                    <div className="flex-1">
                        <h2 className="text-3xl font-serif mb-8 text-primary">
                            {currentStep === 0 && "Select to enhance your stay"}
                            {currentStep === 1 && "Enter your details"}
                            {currentStep === 2 && "Select a payment method"}
                        </h2>

                        {/* Step 1: Addons */}
                        {currentStep === 0 && (
                            <div className="space-y-6">
                                {/* Category Tabs */}
                                {categories && categories.length > 0 && (
                                    <div className="flex flex-wrap gap-2 mb-8">
                                        {categories.map((cat) => (
                                            <button
                                                key={cat.categoryId}
                                                onClick={() => handleCategoryChange(cat.categoryId)}
                                                className={`px-6 py-3 text-xs font-bold uppercase tracking-widest border transition-all ${activeCategory === cat.categoryId
                                                    ? 'bg-primary text-white border-primary'
                                                    : 'bg-gray-100 text-gray-500 border-transparent hover:bg-gray-200'
                                                    }`}
                                            >
                                                {cat.name}
                                            </button>
                                        ))}
                                    </div>
                                )}

                                {addons.length > 0 ? (
                                    addons.map((addon) => (
                                        <div key={addon.addonId} className="flex border border-gray-100 p-4 hover:border-accent transition-colors">
                                            <div className="w-24 h-24 bg-gray-100 flex-shrink-0 mr-6">
                                                {addon.pictureUrl ? (
                                                    <img src={addon.pictureUrl} alt={addon.name} className="w-full h-full object-cover" />
                                                ) : (
                                                    <div className="w-full h-full bg-gray-200 flex items-center justify-center text-gray-400 text-[10px]">
                                                        {addon.name.substring(0, 2)}
                                                    </div>
                                                )}
                                            </div>
                                            <div className="flex-1">
                                                <h4 className="font-bold uppercase tracking-wider text-sm mb-1">{addon.name}</h4>
                                                <p className="text-xs text-gray-500 mb-2">{addon.description}</p>
                                                <div className="text-sm font-serif text-accent">THB {addon.price.toLocaleString()} <span className="text-xs text-gray-300 font-sans">/ {addon.unitName}</span></div>
                                            </div>
                                            <div className="flex items-center">
                                                <button
                                                    onClick={() => dispatch(toggleAddon(addon.addonId))}
                                                    className={`px-6 py-2 text-xs font-bold uppercase tracking-widest border transition-all ${selectedAddons[addon.addonId]
                                                        ? 'bg-primary text-white border-primary'
                                                        : 'bg-white text-primary border-gray-200 hover:border-primary'
                                                        }`}
                                                >
                                                    {selectedAddons[addon.addonId] ? 'Selected' : 'Select'}
                                                </button>
                                            </div>
                                        </div>
                                    ))
                                ) : (
                                    <div className="py-8 text-center text-gray-400 text-sm border border-dashed border-gray-200">
                                        No items available in this category.
                                    </div>
                                )}
                            </div>
                        )}

                        {/* Step 2: Guest Details - Refined */}
                        {currentStep === 1 && (
                            <div>
                                <h2 className="text-2xl font-serif mb-6 text-primary">Enter your details</h2>
                                <div className="max-w-3xl">
                                    <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
                                        <div className="md:col-span-1">
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Title</label>
                                            <select
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary bg-white"
                                                value={guestDetails.title}
                                                onChange={(e) => dispatch(updateGuestDetails({ title: e.target.value }))}
                                            >
                                                <option value="">Select</option>
                                                <option value="Mr">Mr.</option>
                                                <option value="Ms">Ms.</option>
                                                <option value="Mrs">Mrs.</option>
                                            </select>
                                        </div>
                                        <div className="md:col-span-1">
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">First Name</label>
                                            <input
                                                type="text"
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary"
                                                value={guestDetails.firstName}
                                                onChange={(e) => dispatch(updateGuestDetails({ firstName: e.target.value }))}
                                            />
                                        </div>
                                        <div className="md:col-span-2">
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Last Name</label>
                                            <input
                                                type="text"
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary"
                                                value={guestDetails.lastName}
                                                onChange={(e) => dispatch(updateGuestDetails({ lastName: e.target.value }))}
                                            />
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                                        <div>
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Email</label>
                                            <input
                                                type="email"
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary"
                                                value={guestDetails.email}
                                                onChange={(e) => dispatch(updateGuestDetails({ email: e.target.value }))}
                                            />
                                        </div>
                                        <div>
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Email (re-confirm)</label>
                                            <input
                                                type="email"
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary"
                                                placeholder="Confirm Email"
                                            />
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                                        <div>
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Country of Passport</label>
                                            <select
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary bg-white"
                                                value={guestDetails.country}
                                                onChange={(e) => dispatch(updateGuestDetails({ country: e.target.value }))}
                                            >
                                                <option value="">- Please Select -</option>
                                                <option value="Thailand">Thailand</option>
                                                <option value="United States">United States</option>
                                                <option value="United Kingdom">United Kingdom</option>
                                                {/* Add more countries */}
                                            </select>
                                        </div>
                                        <div>
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Phone Number</label>
                                            <input
                                                type="tel"
                                                className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary"
                                                value={guestDetails.phone}
                                                onChange={(e) => dispatch(updateGuestDetails({ phone: e.target.value }))}
                                            />
                                        </div>
                                    </div>

                                    <div className="mb-8">
                                        <label className="flex items-center gap-2 text-sm text-gray-500">
                                            <input type="checkbox" className="rounded border-gray-300 text-primary focus:ring-primary" />
                                            I want to add all guest names
                                        </label>
                                    </div>

                                    {/* Optional Information */}
                                    <h3 className="text-lg font-serif mb-4 text-primary">Optional information</h3>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                                        <div>
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Arrive by</label>
                                            <select className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary bg-white">
                                                <option>------------</option>
                                                <option>Plane</option>
                                                <option>Car</option>
                                            </select>
                                        </div>
                                        <div className="md:row-span-2">
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Your special request</label>
                                            <textarea
                                                className="w-full border border-gray-200 rounded p-3 text-sm focus:outline-none focus:border-accent h-32"
                                                placeholder="e.g. bed preference, pick-up or drop-off location"
                                            ></textarea>
                                            <p className="text-[10px] text-gray-400 mt-1">Special request is not guaranteed. Your request will be taken care where possible.</p>
                                        </div>
                                        <div>
                                            <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-2">Arrival details</label>
                                            <input type="text" className="w-full border-b border-gray-200 py-2 focus:outline-none focus:border-accent text-primary" />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Step 3: Payment - Refined */}
                        {currentStep === 2 && (
                            <div>
                                <h2 className="text-2xl font-serif mb-6 flex items-center gap-2 text-primary">
                                    <span className="text-green-500"><ShieldCheck className="w-6 h-6" /></span>
                                    Select a payment method
                                </h2>

                                <div className="space-y-6 max-w-2xl">
                                    {/* Credit Card Option - Expanded */}
                                    <div className={`border rounded-lg overflow-hidden transition-all ${paymentMethod === 'card' ? 'border-green-500 shadow-sm' : 'border-gray-200'}`}>
                                        <div
                                            className="p-4 flex items-center justify-between cursor-pointer bg-white"
                                            onClick={() => dispatch(setPaymentMethod('card'))}
                                        >
                                            <div className="flex items-center gap-3">
                                                <CardIcon className={`w-5 h-5 ${paymentMethod === 'card' ? 'text-green-600' : 'text-gray-400'}`} />
                                                <span className={`font-bold ${paymentMethod === 'card' ? 'text-green-700' : 'text-gray-700'}`}>Credit Card</span>
                                            </div>
                                            {paymentMethod === 'card' && <Check className="w-5 h-5 text-green-600" />}
                                        </div>

                                        {/* Card Details Form - Only show if selected */}
                                        {paymentMethod === 'card' && (
                                            <div className="p-6 border-t border-gray-100 bg-white">
                                                <div className="flex justify-between items-start mb-6">
                                                    <div className="text-xs text-gray-500 font-bold uppercase tracking-wider mb-2">We accept the following payment methods</div>
                                                    <div className="flex gap-2">
                                                        <div className="w-10 h-6 bg-blue-600 text-white text-[8px] flex items-center justify-center font-bold italic">VISA</div>
                                                        <div className="w-10 h-6 bg-red-500 text-white text-[8px] flex items-center justify-center font-bold">MC</div>
                                                    </div>
                                                </div>

                                                <div className="space-y-4 max-w-md">
                                                    <div>
                                                        <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-1">Card number</label>
                                                        <input type="text" className="w-full border border-gray-200 rounded p-2 focus:ring-1 focus:ring-green-500 outline-none" placeholder="0000 0000 0000 0000" />
                                                    </div>
                                                    <div>
                                                        <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-1">Card holder name</label>
                                                        <input type="text" className="w-full border border-gray-200 rounded p-2 focus:ring-1 focus:ring-green-500 outline-none" />
                                                    </div>
                                                    <div className="w-32">
                                                        <label className="text-xs font-bold uppercase tracking-wider text-gray-400 block mb-1">Expiry date</label>
                                                        <input type="text" className="w-full border border-gray-200 rounded p-2 focus:ring-1 focus:ring-green-500 outline-none" placeholder="MM/YY" />
                                                    </div>

                                                    <div className="bg-yellow-50 border border-yellow-100 p-3 rounded flex items-center gap-2 mt-4">
                                                        <div className="w-4 h-4 bg-yellow-400 rounded-full flex items-center justify-center text-[10px] text-white font-bold">!</div>
                                                        <div className="text-xs text-gray-600">
                                                            <strong>Booking with confidence.</strong><br />
                                                            Your credit card details are safe via this secured payment form.
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        )}
                                    </div>

                                    {/* Placeholder for QR */}
                                    <div className={`border rounded-lg p-4 cursor-pointer flex items-center gap-3 ${paymentMethod === 'qr' ? 'border-green-500' : 'border-gray-200'}`} onClick={() => dispatch(setPaymentMethod('qr'))}>
                                        <div className="w-5 h-5 border-2 border-gray-300 rounded-full flex items-center justify-center">
                                            {paymentMethod === 'qr' && <div className="w-3 h-3 bg-green-500 rounded-full"></div>}
                                        </div>
                                        <span className="font-bold text-gray-700">QR Payment</span>
                                    </div>
                                </div>
                            </div>
                        )}

                        <div className="mt-12 pt-6 border-t border-gray-100 flex justify-end">
                            <button
                                onClick={handleNext}
                                className="bg-primary hover:bg-gray-800 text-white px-10 py-4 text-xs font-bold uppercase tracking-[0.2em] transition-all"
                            >
                                {currentStep === steps.length - 1 ? "Complete Reservation" : "Continue"}
                            </button>
                        </div>
                    </div>

                    {/* Sticky Sidebar */}
                    <div className="hidden lg:block w-96 relative">
                        <div className="sticky top-8 bg-gray-50 border border-gray-100 p-6">
                            <h3 className="font-serif text-xl border-b border-gray-200 pb-4 mb-6">Reservation Summary</h3>

                            <div className="flex gap-4 mb-6">
                                <div className="w-20 h-20 bg-gray-200 overflow-hidden">
                                    {roomImage ? (
                                        <img src={roomImage} alt={roomName} className="w-full h-full object-cover" />
                                    ) : (
                                        <div className="w-full h-full flex items-center justify-center text-xs text-gray-400">No Image</div>
                                    )}
                                </div>
                                <div>
                                    <div className="text-xs font-bold uppercase tracking-wider text-primary mb-1">{roomName}</div>
                                    <div className="flex text-[10px] text-gray-500 gap-1 mb-1">
                                        <Star className="w-3 h-3 text-accent" />
                                        <Star className="w-3 h-3 text-accent" />
                                        <Star className="w-3 h-3 text-accent" />
                                        <Star className="w-3 h-3 text-accent" />
                                        <Star className="w-3 h-3 text-accent" />
                                    </div>
                                    <div className="text-xs text-gray-400">{bookingState.adults} Adults, {bookingState.children} Children</div>
                                </div>
                            </div>

                            <div className="space-y-3 text-sm text-gray-600 border-b border-gray-200 pb-6 mb-6">
                                <div className="flex justify-between">
                                    <span>Check-In</span>
                                    <span className="font-bold text-primary">{bookingState.checkIn}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span>Check-Out</span>
                                    <span className="font-bold text-primary">{bookingState.checkOut}</span>
                                </div>
                            </div>

                            <div className="space-y-3 text-sm border-b border-gray-200 pb-6 mb-6">
                                <div className="flex justify-between text-gray-500">
                                    <span>{ratePlanName} (1 Night)</span>
                                    <span>THB {roomTotal.toLocaleString()}</span>
                                </div>
                                {Object.entries(selectedAddons).map(([id, qty]) => {
                                    if (qty === 0) return null;
                                    const addon = addons.find(a => a.addonId === parseInt(id));
                                    if (!addon) return null;
                                    return (
                                        <div key={id} className="flex justify-between text-gray-500">
                                            <span>{addon.name}</span>
                                            <span>THB {addon.price.toLocaleString()}</span>
                                        </div>
                                    );
                                })}
                                <div className="flex justify-between text-gray-500">
                                    <span>Taxes & Fees</span>
                                    <span>Included</span>
                                </div>
                            </div>

                            <div className="flex justify-between items-baseline">
                                <span className="font-serif text-lg">Total</span>
                                <span className="font-serif text-2xl text-accent">THB {grandTotal.toLocaleString()}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>


            {/* Mobile Bottom Summary Bar */}
            <div className="lg:hidden fixed bottom-0 left-0 w-full bg-white border-t border-gray-200 p-4 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.1)] z-40">
                <div className="flex justify-between items-center">
                    <div onClick={() => setShowMobileSummary(!showMobileSummary)} className="flex flex-col cursor-pointer">
                        <span className="text-xs text-gray-500 flex items-center gap-1">
                            Total (1 Night) {showMobileSummary ? <ChevronDown className="w-3 h-3" /> : <ChevronUp className="w-3 h-3" />}
                        </span>
                        <span className="font-serif text-xl text-accent font-bold">THB {grandTotal.toLocaleString()}</span>
                    </div>
                    <button
                        onClick={handleNext}
                        className="bg-primary text-white px-6 py-3 rounded text-xs font-bold uppercase tracking-wider"
                    >
                        {currentStep === steps.length - 1 ? "Book" : "Next"}
                    </button>
                </div>

                {/* Mobile Summary Details (Expandable) */}
                {showMobileSummary && (
                    <div className="mt-4 pt-4 border-t border-gray-100 max-h-[50vh] overflow-y-auto animate-fade-in-up">
                        <div className="flex gap-4 mb-4">
                            <div className="w-16 h-16 bg-gray-200 rounded overflow-hidden flex-shrink-0">
                                {roomImage ? (
                                    <img src={roomImage} alt={roomName} className="w-full h-full object-cover" />
                                ) : (
                                    <div className="w-full h-full flex items-center justify-center text-[10px] text-gray-400">No Img</div>
                                )}
                            </div>
                            <div>
                                <div className="text-xs font-bold uppercase text-primary mb-1">{roomName}</div>
                                <div className="text-xs text-gray-400">{bookingState.adults} Adults, {bookingState.children} Children</div>
                                <div className="text-xs text-gray-400">{bookingState.checkIn} - {bookingState.checkOut}</div>
                            </div>
                        </div>
                        <div className="space-y-2 text-sm text-gray-600">
                            <div className="flex justify-between">
                                <span>{ratePlanName}</span>
                                <span>THB {roomTotal.toLocaleString()}</span>
                            </div>
                            {Object.entries(selectedAddons).map(([id, qty]) => {
                                if (qty === 0) return null;
                                const addon = addons.find(a => a.addonId === parseInt(id));
                                if (!addon) return null;
                                return (
                                    <div key={id} className="flex justify-between text-gray-500 text-xs">
                                        <span>{addon.name}</span>
                                        <span>THB {addon.price.toLocaleString()}</span>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                )}
            </div>
        </div >
    );
};

export default BookingPage;
