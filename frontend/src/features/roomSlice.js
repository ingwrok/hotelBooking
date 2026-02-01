import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getRoomTypes, getRoomTypeDetails, getAvailabilityCounts } from '../api';

export const fetchRooms = createAsyncThunk('rooms/fetchRooms', async (params) => {
    // params = { checkIn, checkOut } (optional)
    const rooms = await getRoomTypes();

    if (params && params.checkIn && params.checkOut) {
        console.log("Fetching rooms with params:", params); // DEBUG
        // Filter by availability
        const counts = await getAvailabilityCounts(params.checkIn, params.checkOut);
        console.log("Availability Counts:", counts); // DEBUG

        // Return only rooms with availability > 0
        // We can also attach the count to the room object if needed for UI
        const filtered = rooms.map(room => ({
            ...room,
            availableCount: counts[room.roomTypeId] || 0
        })).filter(room => room.availableCount > 0 && room.totalRooms > 0);

        console.log("Filtered Rooms:", filtered); // DEBUG
        return filtered;
    }

    // Even without dates, filter out ghost rooms (no physical rooms created)
    return rooms.filter(r => r.totalRooms > 0);
});

export const fetchRoomDetails = createAsyncThunk('rooms/fetchRoomDetails', async (id) => {
    const response = await getRoomTypeDetails(id);
    return response;
});

const initialState = {
    rooms: [],
    currentRoom: null,
    loading: false,
    error: null,
};

const roomSlice = createSlice({
    name: 'rooms',
    initialState,
    reducers: {
        clearCurrentRoom: (state) => {
            state.currentRoom = null;
        }
    },
    extraReducers: (builder) => {
        builder
            // Fetch Rooms
            .addCase(fetchRooms.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchRooms.fulfilled, (state, action) => {
                state.loading = false;
                state.rooms = action.payload;
            })
            .addCase(fetchRooms.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message;
            })
            // Fetch Room Details
            .addCase(fetchRoomDetails.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchRoomDetails.fulfilled, (state, action) => {
                state.loading = false;
                state.currentRoom = action.payload;
            })
            .addCase(fetchRoomDetails.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message;
            });
    },
});

export const { clearCurrentRoom } = roomSlice.actions;
export default roomSlice.reducer;
