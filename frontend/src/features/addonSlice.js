import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getAddonCategories, getAddonsByCategory } from '../api';

export const fetchAddonCategories = createAsyncThunk('addons/fetchCategories', async () => {
    const response = await getAddonCategories();
    return response;
});

export const fetchAddonsByCategory = createAsyncThunk('addons/fetchByCategory', async (categoryId) => {
    const response = await getAddonsByCategory(categoryId);
    return response;
});

const initialState = {
    categories: [],
    addons: [],
    loading: false,
    error: null,
};

const addonSlice = createSlice({
    name: 'addons',
    initialState,
    reducers: {
        clearAddons: (state) => {
            state.addons = [];
        }
    },
    extraReducers: (builder) => {
        builder
            // Categories
            .addCase(fetchAddonCategories.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchAddonCategories.fulfilled, (state, action) => {
                state.loading = false;
                state.categories = action.payload;
            })
            .addCase(fetchAddonCategories.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message;
            })
            // Addons
            .addCase(fetchAddonsByCategory.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchAddonsByCategory.fulfilled, (state, action) => {
                state.loading = false;
                state.addons = action.payload;
            })
            .addCase(fetchAddonsByCategory.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message;
            });
    },
});

export const { clearAddons } = addonSlice.actions;
export default addonSlice.reducer;
