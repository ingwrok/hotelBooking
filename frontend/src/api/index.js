import axios from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api',
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
});

const handleApiError = (error, functionName) => {
    const message = error.response?.data?.error || error.response?.data?.message || error.message;
    console.error(`API Error (${functionName}):`, message);
    throw new Error(message);
};

export const getRoomTypes = async () => {
    try {
        const response = await api.get('/room_types');
        return response.data;
    } catch (e) {
        handleApiError(e, "getRoomTypes");
    }
};

export const getRoomTypeDetails = async (id) => {
    try {
        const response = await api.get(`/room_types/${id}`);
        return response.data;
    } catch (e) {
        handleApiError(e, "getRoomTypeDetails");
    }
}

export const getRatePlansByRoomType = async (roomTypeId) => {
    try {
        const response = await api.get(`/rate_plans/room-types/${roomTypeId}`);
        return response.data;
    } catch (e) {
        handleApiError(e, "getRatePlansByRoomType");
    }
}

export const getAddonCategories = async () => {
    try {
        const response = await api.get('/addon-categories');
        return response.data;
    } catch (e) {
        handleApiError(e, "getAddonCategories");
    }
}

export const getAddonsByCategory = async (catId) => {
    try {
        const response = await api.get(`/addons/category/${catId}`);
        return response.data;
    } catch (e) {
        handleApiError(e, "getAddonsByCategory");
    }
}

export const createBooking = async (bookingData) => {
    try {
        const response = await api.post('/bookings', bookingData);
        return response.data;
    } catch (e) {
        handleApiError(e, "createBooking");
    }
}

export const checkAvailability = async (checkIn, checkOut, roomTypeId, count) => {
    try {
        const response = await api.post('/rooms/availability/count', {
            checkIn: checkIn,
            checkOut: checkOut
        });

        const availableCount = response.data[roomTypeId] || 0;
        return availableCount >= count;
    } catch (e) {
        handleApiError(e, "checkAvailability");
    }
}

export const getAvailabilityCounts = async (checkIn, checkOut) => {
    try {
        const response = await api.post('/rooms/availability/count', {
            checkIn: checkIn,
            checkOut: checkOut
        });
        return response.data; 
    } catch (e) {
        handleApiError(e, "getAvailabilityCounts");
    }
}

export const payBooking = async (bookingId) => {
    try {
        const response = await api.post(`/bookings/${bookingId}/pay`);
        return response.data;
    } catch (e) {
        handleApiError(e, "payBooking");
    }
}

export default api;