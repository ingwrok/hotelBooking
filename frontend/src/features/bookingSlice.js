import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { createBooking, checkAvailability, payBooking as apiPayBooking } from '../api';

// Async Thunks
export const submitBooking = createAsyncThunk('booking/submitBooking', async (bookingData) => {
    const response = await createBooking(bookingData);
    return response;
});

export const checkRoomAvailability = createAsyncThunk('booking/checkAvailability', async ({ checkIn, checkOut, roomTypeId, count }) => {
    const isAvailable = await checkAvailability(checkIn, checkOut, roomTypeId, count);
    if (!isAvailable) {
        throw new Error("Room not available for selected dates");
    }
    return isAvailable;
});

export const payBooking = createAsyncThunk('booking/payBooking', async (bookingId) => {
    await apiPayBooking(bookingId);
    return bookingId;
});

const initialState = {
    currentBooking: {
        roomId: null,
        checkIn: null,
        checkOut: null,
        adults: 1,
        children: 0,
        roomCount: 1,
        selectedAddons: {},
        guestDetails: {
            title: '',
            firstName: '',
            lastName: '',
            email: '',
            phone: '',
            country: ''
        },
        paymentMethod: 'qr',
    },
    submissionStatus: 'idle', // idle, loading, succeeded, failed
    availabilityStatus: 'idle', // idle, loading, available, unavailable
    paymentStatus: 'idle', // idle, loading, succeeded, failed
    error: null,
};

const bookingSlice = createSlice({
    name: 'booking',
    initialState,
    reducers: {
        updateBookingDetails: (state, action) => {
            state.currentBooking = { ...state.currentBooking, ...action.payload };
            if (action.payload.checkIn || action.payload.checkOut || action.payload.roomId) {
                state.availabilityStatus = 'idle';
            }
        },
        toggleAddon: (state, action) => {
            const addonId = action.payload;
            const currentQty = state.currentBooking.selectedAddons[addonId] || 0;
            const newQty = currentQty > 0 ? 0 : 1;

            state.currentBooking.selectedAddons = {
                ...state.currentBooking.selectedAddons,
                [addonId]: newQty
            };
        },
        updateGuestDetails: (state, action) => {
            state.currentBooking.guestDetails = {
                ...state.currentBooking.guestDetails,
                ...action.payload
            };
        },
        setPaymentMethod: (state, action) => {
            state.currentBooking.paymentMethod = action.payload;
        },
        resetBooking: (state) => {
            state.currentBooking = initialState.currentBooking;
            state.submissionStatus = 'idle';
            state.availabilityStatus = 'idle';
            state.paymentStatus = 'idle';
            state.error = null;
        },
        resetPaymentStatus: (state) => {
            state.paymentStatus = 'idle';
        },
    },
    extraReducers: (builder) => {
        builder
            // Check Availability
            .addCase(checkRoomAvailability.pending, (state) => {
                state.availabilityStatus = 'loading';
                state.error = null;
            })
            .addCase(checkRoomAvailability.fulfilled, (state) => {
                state.availabilityStatus = 'available';
            })
            .addCase(checkRoomAvailability.rejected, (state, action) => {
                state.availabilityStatus = 'unavailable';
                state.error = action.error.message;
            })
            // Submit Booking
            .addCase(submitBooking.pending, (state) => {
                state.submissionStatus = 'loading';
                state.error = null;
            })
            .addCase(submitBooking.fulfilled, (state) => {
                state.submissionStatus = 'succeeded';
            })
            .addCase(submitBooking.rejected, (state, action) => {
                state.submissionStatus = 'failed';
                state.error = action.error.message;
            })
            // Pay Booking
            .addCase(payBooking.pending, (state) => {
                state.paymentStatus = 'loading';
            })
            .addCase(payBooking.fulfilled, (state) => {
                state.paymentStatus = 'succeeded';
            })
            .addCase(payBooking.rejected, (state, action) => {
                state.paymentStatus = 'failed';
                state.error = action.error.message;
            });
    },
});

export const { updateBookingDetails, toggleAddon, updateGuestDetails, setPaymentMethod, resetBooking, resetPaymentStatus } = bookingSlice.actions;
export default bookingSlice.reducer;
