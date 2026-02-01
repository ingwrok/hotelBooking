import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { CheckCircle } from 'lucide-react';
import Button from '../components/common/Button.jsx';
import { useDispatch, useSelector } from 'react-redux';
import { payBooking } from '../features/bookingSlice.js';

const ConfirmationPage = () => {
    const location = useLocation();
    const booking = location.state?.booking;
    const [timeLeft, setTimeLeft] = useState('');
    const [isExpired, setIsExpired] = useState(false);

    // Payment Logic
    const dispatch = useDispatch();
    const { paymentStatus } = useSelector((state) => state.booking);
    const [isPaymentLoading, setIsPaymentLoading] = useState(false);

    const handlePayment = async () => {
        setIsPaymentLoading(true);
        await dispatch(payBooking(booking.bookingId));
        setIsPaymentLoading(false);
        booking.status = 'confirmed';
    };

    useEffect(() => {
        if (!booking || !booking.expiredAt) return;

        const interval = setInterval(() => {
            const now = new Date().getTime();
            const expiration = new Date(booking.expiredAt).getTime();
            const distance = expiration - now;

            if (distance < 0) {
                clearInterval(interval);
                setTimeLeft("EXPIRED");
                setIsExpired(true);
            } else {
                const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
                const seconds = Math.floor((distance % (1000 * 60)) / 1000);
                setTimeLeft(`${minutes}m ${seconds}s`);
            }
        }, 1000);

        return () => clearInterval(interval);
    }, [booking]);

    if (!booking) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="text-center">
                    <h1 className="text-2xl font-bold text-gray-800">Booking Not Found</h1>
                    <Link to="/" className="text-primary hover:underline mt-4 block">Return Home</Link>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
            <div className="bg-white p-8 rounded-lg shadow-xl text-center max-w-lg w-full">
                <div className="flex justify-center mb-6">
                    <CheckCircle className="w-16 h-16 text-green-500" />
                </div>
                <h1 className="text-3xl font-serif font-bold text-gray-800 mb-2">Booking Received!</h1>
                <p className="text-gray-500 mb-6">
                    Your booking ID is <span className="font-bold text-black">#{booking.bookingId}</span>.
                    <br />A confirmation email has been sent.
                </p>

                <div className={`mb-8 p-4 rounded-lg ${isExpired ? 'bg-red-50 border border-red-200' : 'bg-green-50 border border-green-200'}`}>
                    <p className="text-xs uppercase tracking-widest text-gray-500 mb-1">Time to Complete Payment</p>
                    <div className={`text-3xl font-bold font-mono ${isExpired ? 'text-red-600' : 'text-green-700'}`}>
                        {timeLeft}
                    </div>
                    {isExpired && <p className="text-red-500 text-xs mt-2">This booking has expired. Please make a new booking.</p>}
                </div>

                <div className="bg-gray-50 p-6 rounded-lg mb-8 text-left">
                    <h3 className="font-bold text-sm uppercase tracking-wider mb-4 border-b border-gray-200 pb-2">Summary</h3>
                    <div className="flex justify-between mb-2">
                        <span className="text-gray-500 text-sm">Room</span>
                        <span className="font-bold text-sm">{booking.roomTypeName}</span>
                    </div>
                    <div className="flex justify-between mb-2">
                        <span className="text-gray-500 text-sm">Dates</span>
                        <span className="font-bold text-sm">{new Date(booking.checkInDate).toLocaleDateString()} - {new Date(booking.checkOutDate).toLocaleDateString()}</span>
                    </div>
                    <div className="flex justify-between mt-4 pt-4 border-t border-gray-200">
                        <span className="text-gray-800 font-bold">Total Price</span>
                        <span className="font-serif text-xl text-primary">THB {booking.totalPrice.toLocaleString()}</span>
                    </div>
                </div>

                {/* Payment Section */}
                {booking.status === 'pending' && !isExpired && (
                    <div className="mb-6">
                        <Button fullWidth onClick={handlePayment} disabled={isPaymentLoading} className="bg-green-600 hover:bg-green-700">
                            {isPaymentLoading ? 'Processing...' : 'Secure Pay Now (Simulate)'}
                        </Button>
                        <p className="text-xs text-gray-400 mt-2">
                            This is a simulation. No real money will be charged.
                        </p>
                    </div>
                )}

                {booking.status === 'confirmed' && (
                    <div className="mb-6 p-4 bg-green-50 border border-green-200 text-green-700 rounded-lg">
                        <div className="font-bold flex items-center justify-center gap-2">
                            <CheckCircle className="w-5 h-5" /> Payment Successful
                        </div>
                        <p className="text-sm mt-1">Your booking is fully confirmed!</p>
                    </div>
                )}

                <div className="space-y-4">
                    <Link to="/">
                        <Button fullWidth variant="outline">Return to Home</Button>
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default ConfirmationPage;
