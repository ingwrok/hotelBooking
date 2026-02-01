
import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../api';

const AdminRoomTypePage = () => {
    const { id } = useParams();
    const isEditMode = !!id;
    const navigate = useNavigate();

    const [formData, setFormData] = useState({
        name: '',
        description: '',
        sizeSqm: 0,
        bedType: '',
        capacity: 0,
        pictureUrl: [],
        amenityIds: []
    });
    const [uploading, setUploading] = useState(false);
    const [loading, setLoading] = useState(false);

    const [ratePlans, setRatePlans] = useState([]);
    const [prices, setPrices] = useState({});

    useEffect(() => {
        if (isEditMode) {
            const fetchData = async () => {
                try {
                    setLoading(true);

                    // Fetch Room Type Details
                    const rtRes = await api.get(`/room_types/${id}`);
                    setFormData({
                        name: rtRes.data.name,
                        description: rtRes.data.description,
                        sizeSqm: rtRes.data.sizeSqm,
                        bedType: rtRes.data.bedType,
                        capacity: rtRes.data.capacity,
                        pictureUrl: rtRes.data.pictureUrl || [],
                        amenityIds: rtRes.data.amenityIds || []
                    });

                    // Fetch All Rate Plans
                    const rpRes = await api.get('/rate_plans');
                    setRatePlans(rpRes.data);

                    // Fetch Existing Prices for this Room Type
                    try {
                        const priceRes = await api.get(`/rate_plans/room-types/${id}`);
                        const priceMap = {};
                        priceRes.data.forEach(p => {
                            priceMap[p.ratePlanId] = p.price;
                        });
                        setPrices(priceMap);
                    } catch (err) {
                        console.warn("No existing prices found or error fetching prices", err);
                    }

                } catch (error) {
                    console.error("Error fetching data:", error);
                    alert("Failed to load details.");
                    navigate('/admin/room-types');
                } finally {
                    setLoading(false);
                }
            };
            fetchData();
        }
    }, [id, isEditMode, navigate]);

    const handlePriceChange = (ratePlanId, value) => {
        setPrices(prev => ({
            ...prev,
            [ratePlanId]: value
        }));
    };

    const handleSavePrice = async (ratePlanId) => {
        const price = prices[ratePlanId];
        if (!price && price !== 0) {
            alert("Please enter a price");
            return;
        }

        try {
            await api.put(`/rate_plans/${ratePlanId}/room-types/${id}`, {
                price: Number(price)
            });
            alert("Price updated successfully!");
        } catch (error) {
            console.error("Error saving price:", error);
            alert("Failed to save price.");
        }
    };

    const handleDeletePrice = async (ratePlanId) => {
        if (!window.confirm("Are you sure you want to unlink this price?")) return;
        try {
            await api.delete(`/rate_plans/${ratePlanId}/room-types/${id}`);
            setPrices(prev => {
                const newPrices = { ...prev };
                delete newPrices[ratePlanId];
                return newPrices;
            });
            alert("Price removed successfully!");
        } catch (error) {
            console.error("Error removing price:", error);
            alert("Failed to remove price.");
        }
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: name === 'sizeSqm' || name === 'capacity' ? Number(value) : value
        }));
    };

    const handleFileChange = (e) => {
        if (e.target.files && e.target.files[0]) {
            handleUpload(e.target.files[0]);
        }
    };

    const handleUpload = async (file) => {
        setUploading(true);
        const uploadData = new FormData();
        uploadData.append('image', file);

        try {
            const response = await api.post('/room_types/upload', uploadData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            const data = response.data;
            setFormData(prev => ({
                ...prev,
                pictureUrl: [...prev.pictureUrl, data.url]
            }));
        } catch (error) {
            console.error('Error uploading image:', error);
            alert('Failed to upload image.');
        } finally {
            setUploading(false);
        }
    };

    const removeImage = (indexToRemove) => {
        setFormData(prev => ({
            ...prev,
            pictureUrl: prev.pictureUrl.filter((_, index) => index !== indexToRemove)
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (isEditMode) {
                await api.patch(`/room_types/${id}`, formData);
                alert('Room Type updated successfully!');
            } else {
                await api.post('/room_types', formData);
                alert('Room Type created successfully!');
            }
            navigate('/admin/room-types');
        } catch (error) {
            console.error('Error saving room type:', error);
            alert('Error saving room type: ' + (error.response?.data?.message || error.message));
        }
    };

    if (loading) return <div className="p-8 text-center">Loading...</div>;

    return (
        <div className="w-full max-w-4xl mx-auto">
            <h1 className="text-3xl font-serif mb-8">{isEditMode ? 'Edit Room Type' : 'Create New Room Type'}</h1>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Form */}
                <div className="lg:col-span-2">
                    <form onSubmit={handleSubmit} className="space-y-6 bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Name</label>
                            <input type="text" name="name" value={formData.name} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">Description</label>
                            <textarea name="description" value={formData.description} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" rows="3" />
                        </div>

                        <div className="grid grid-cols-2 gap-6">
                            <div>
                                <label className="block text-sm font-medium text-gray-700">Size (SQM)</label>
                                <input type="number" name="sizeSqm" value={formData.sizeSqm} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700">Capacity</label>
                                <input type="number" name="capacity" value={formData.capacity} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" />
                            </div>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">Bed Type</label>
                            <input type="text" name="bedType" value={formData.bedType} onChange={handleChange} required className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2" placeholder="e.g. King Bed" />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Room Images</label>
                            <div className="flex items-center gap-4 mb-4">
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleFileChange}
                                    className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary file:text-white hover:file:bg-opacity-90"
                                    disabled={uploading}
                                />
                                {uploading && <span className="text-sm text-yellow-600 animate-pulse">Uploading...</span>}
                            </div>

                            {formData.pictureUrl.length > 0 ? (
                                <div className="grid grid-cols-3 gap-4">
                                    {formData.pictureUrl.map((url, idx) => (
                                        <div key={idx} className="relative group">
                                            <img src={url} alt={`Room ${idx + 1} `} className="h-24 w-full object-cover rounded shadow-sm" />
                                            <button
                                                type="button"
                                                onClick={() => removeImage(idx)}
                                                className="absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                                                title="Remove Image"
                                            >
                                                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                                                    <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                                                </svg>
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-sm text-gray-400 italic">No images uploaded yet.</p>
                            )}
                        </div>

                        <div className="flex gap-4 pt-4">
                            <button
                                type="button"
                                onClick={() => navigate('/admin/room-types')}
                                className="flex-1 py-3 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="flex-1 py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
                            >
                                {isEditMode ? 'Update Room Type' : 'Create Room Type'}
                            </button>
                        </div>
                    </form>
                </div>

                {/* Pricing Section - Only in Edit Mode */}
                {isEditMode && (
                    <div className="lg:col-span-1">
                        <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 sticky top-6">
                            <h2 className="text-xl font-bold text-gray-800 mb-4 flex items-center gap-2">
                                <span className="p-1 bg-green-100 text-green-600 rounded">
                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                                    </svg>
                                </span>
                                Pricing Mapping
                            </h2>
                            <p className="text-sm text-gray-500 mb-6">Set prices per night for each rate plan.</p>

                            <div className="space-y-4 max-h-[600px] overflow-y-auto pr-2 custom-scrollbar">
                                {ratePlans.length === 0 ? (
                                    <p className="text-gray-400 text-center">No rate plans available.</p>
                                ) : (
                                    ratePlans.map(rp => {
                                        const isLinked = prices[rp.ratePlanId] !== undefined && prices[rp.ratePlanId] !== null;
                                        return (
                                            <div key={rp.ratePlanId} className={`border rounded-lg p-3 transition-colors ${isLinked ? 'border-primary bg-blue-50' : 'border-gray-200 hover:border-gray-300'}`}>
                                                <div className="mb-2 flex justify-between items-start">
                                                    <div>
                                                        <div className="font-medium text-gray-900">{rp.name}</div>
                                                        <div className="text-xs text-gray-500">{rp.description}</div>
                                                    </div>
                                                    {isLinked ? (
                                                        <span className="text-xs font-bold text-primary bg-white px-2 py-0.5 rounded border border-blue-100">
                                                            Linked
                                                        </span>
                                                    ) : (
                                                        <span className="text-xs text-gray-400">
                                                            Unlinked
                                                        </span>
                                                    )}
                                                </div>
                                                <div className="flex flex-wrap gap-2 items-center">
                                                    <input
                                                        type="number"
                                                        placeholder="Price"
                                                        className="flex-1 min-w-[80px] border border-gray-300 rounded px-2 py-1 text-sm"
                                                        value={prices[rp.ratePlanId] || ''}
                                                        onChange={(e) => handlePriceChange(rp.ratePlanId, e.target.value)}
                                                    />
                                                    <div className="flex gap-1 shrink-0">
                                                        <button
                                                            type="button"
                                                            onClick={() => handleSavePrice(rp.ratePlanId)}
                                                            className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-xs font-medium"
                                                        >
                                                            Save
                                                        </button>
                                                        <button
                                                            type="button"
                                                            onClick={() => handleDeletePrice(rp.ratePlanId)}
                                                            className="bg-red-50 hover:bg-red-100 text-red-500 px-2 py-1 rounded border border-red-200"
                                                            title="Unlink Price"
                                                        >
                                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                                                <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                                                            </svg>
                                                        </button>
                                                    </div>
                                                </div>
                                            </div>
                                        );
                                    })
                                )}
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default AdminRoomTypePage;
