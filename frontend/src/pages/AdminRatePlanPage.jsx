import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../api';

const AdminRatePlanPage = () => {
    const { id } = useParams();
    const isEditMode = !!id;
    const navigate = useNavigate();

    const [formData, setFormData] = useState({
        name: '',
        description: '',
        isSpecialPackage: false,
        allowFreeCancel: false,
        allowPayLater: false
    });
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (isEditMode) {
            const fetchRatePlan = async () => {
                try {
                    setLoading(true);
                    const response = await api.get(`/rate_plans/${id}`);
                    const data = response.data;
                    setFormData({
                        name: data.name,
                        description: data.description,
                        isSpecialPackage: data.isSpecialPackage,
                        allowFreeCancel: data.allowFreeCancel,
                        allowPayLater: data.allowPayLater
                    });
                } catch (error) {
                    console.error("Error fetching rate plan:", error);
                    alert("Failed to load rate plan details.");
                    navigate('/admin/rate-plans');
                } finally {
                    setLoading(false);
                }
            };
            fetchRatePlan();
        }
    }, [id, isEditMode, navigate]);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (isEditMode) {
                // Backend expects PUT for Rate Plans
                await api.put(`/rate_plans/${id}`, formData);
                alert('Rate Plan updated successfully!');
            } else {
                await api.post('/rate_plans', formData);
                alert('Rate Plan created successfully!');
            }
            navigate('/admin/rate-plans');
        } catch (error) {
            console.error('Error saving rate plan:', error);
            alert('Error saving rate plan: ' + (error.response?.data?.message || error.message));
        }
    };

    if (loading) return <div className="p-8 text-center">Loading...</div>;

    return (
        <div className="container mx-auto p-8 max-w-2xl">
            <h1 className="text-3xl font-serif mb-8">{isEditMode ? 'Edit Rate Plan' : 'Create New Rate Plan'}</h1>
            <form onSubmit={handleSubmit} className="space-y-6 bg-white p-8 rounded-xl shadow-sm border border-gray-100">
                <div>
                    <label className="block text-sm font-medium text-gray-700">Name</label>
                    <input
                        type="text"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        required
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-3 focus:ring-primary focus:border-primary"
                        placeholder="e.g. Standard Rate with Breakfast"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700">Description</label>
                    <textarea
                        name="description"
                        value={formData.description}
                        onChange={handleChange}
                        required
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-3 focus:ring-primary focus:border-primary"
                        rows="4"
                        placeholder="Describe what this rate plan includes..."
                    />
                </div>

                <div className="space-y-4 pt-4 border-t border-gray-100">
                    <h3 className="font-medium text-gray-900">Policies</h3>

                    <div className="flex items-center">
                        <input
                            id="isSpecialPackage"
                            name="isSpecialPackage"
                            type="checkbox"
                            checked={formData.isSpecialPackage}
                            onChange={handleChange}
                            className="h-5 w-5 text-primary focus:ring-primary border-gray-300 rounded"
                        />
                        <label htmlFor="isSpecialPackage" className="ml-3 block text-sm text-gray-700">
                            Is Special Package
                            <span className="block text-xs text-gray-500">Enable this for promotions or bundles</span>
                        </label>
                    </div>

                    <div className="flex items-center">
                        <input
                            id="allowFreeCancel"
                            name="allowFreeCancel"
                            type="checkbox"
                            checked={formData.allowFreeCancel}
                            onChange={handleChange}
                            className="h-5 w-5 text-primary focus:ring-primary border-gray-300 rounded"
                        />
                        <label htmlFor="allowFreeCancel" className="ml-3 block text-sm text-gray-700">
                            Allow Free Cancellation
                            <span className="block text-xs text-gray-500">Guests can cancel without penalty</span>
                        </label>
                    </div>

                    <div className="flex items-center">
                        <input
                            id="allowPayLater"
                            name="allowPayLater"
                            type="checkbox"
                            checked={formData.allowPayLater}
                            onChange={handleChange}
                            className="h-5 w-5 text-primary focus:ring-primary border-gray-300 rounded"
                        />
                        <label htmlFor="allowPayLater" className="ml-3 block text-sm text-gray-700">
                            Allow Pay Later
                            <span className="block text-xs text-gray-500">Payment collected at the hotel</span>
                        </label>
                    </div>
                </div>

                <div className="flex gap-4 pt-6">
                    <button
                        type="button"
                        onClick={() => navigate('/admin/rate-plans')}
                        className="flex-1 py-3 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className="flex-1 py-3 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-colors"
                    >
                        {isEditMode ? 'Update Rate Plan' : 'Create Rate Plan'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default AdminRatePlanPage;
