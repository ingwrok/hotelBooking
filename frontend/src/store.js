import { configureStore } from '@reduxjs/toolkit';
import roomReducer from './features/roomSlice.js';
import bookingReducer from './features/bookingSlice.js';
import authReducer from './features/authSlice.js';
import addonReducer from './features/addonSlice.js';

const store = configureStore({
    reducer: {
        rooms: roomReducer,
        booking: bookingReducer,
        auth: authReducer,
        addons: addonReducer,
    },
});

export default store;
