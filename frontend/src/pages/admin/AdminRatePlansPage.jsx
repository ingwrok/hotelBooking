import React, { useEffect, useState } from 'react';
import api from '../../api';
import { Link } from 'react-router-dom';
import { Plus, Edit, Trash, DollarSign } from 'lucide-react';

const AdminRatePlansPage = () => {
    const [ratePlans, setRatePlans] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchRatePlans = async () => {
        try {
            setLoading(true);
            const res = await api.get('/rate_plans');
            setRatePlans(res.data);
        } catch (e) {
            console.error(e);
            alert("Failed to fetch rate plans");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRatePlans();
    }, []);

    const handleDelete = async (id) => {
        if (!window.confirm('Are you sure you want to delete this rate plan?')) return;
        try {
            await api.delete(`/rate_plans/${id}`);
            fetchRatePlans();
        } catch (e) {
            console.error(e);
            alert('Failed to delete rate plan. It might be in use.');
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-gray-800">Rate Plans</h2>
                <Link to="new" className="bg-primary hover:bg-gray-800 text-white px-4 py-2 rounded flex items-center gap-2 transition-colors">
                    <Plus className="w-4 h-4" /> New Rate Plan
                </Link>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {ratePlans.map(rp => (
                    <div key={rp.ratePlanId} className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
                        <div className="flex justify-between items-start mb-4">
                            <div>
                                <h3 className="font-bold text-lg text-primary">{rp.name}</h3>
                                <div className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded inline-block mt-1">
                                    ID: {rp.ratePlanId}
                                </div>
                            </div>
                            <div className="bg-green-100 p-2 rounded-full text-green-600">
                                <DollarSign className="w-5 h-5" />
                            </div>
                        </div>

                        <p className="text-sm text-gray-500 mb-4 h-12 line-clamp-2">{rp.description}</p>

                        <div className="space-y-2 mb-6">
                            <div className="flex items-center text-sm">
                                <span className={`w-2 h-2 rounded-full mr-2 ${rp.isSpecialPackage ? 'bg-green-500' : 'bg-gray-300'}`}></span>
                                Special Package: {rp.isSpecialPackage ? 'Yes' : 'No'}
                            </div>
                            <div className="flex items-center text-sm">
                                <span className={`w-2 h-2 rounded-full mr-2 ${rp.allowFreeCancel ? 'bg-green-500' : 'bg-gray-300'}`}></span>
                                Free Cancellation: {rp.allowFreeCancel ? 'Yes' : 'No'}
                            </div>
                            <div className="flex items-center text-sm">
                                <span className={`w-2 h-2 rounded-full mr-2 ${rp.allowPayLater ? 'bg-green-500' : 'bg-gray-300'}`}></span>
                                Pay Later: {rp.allowPayLater ? 'Yes' : 'No'}
                            </div>
                        </div>

                        <div className="flex justify-end gap-2 border-t pt-4">
                            <Link to={`edit/${rp.ratePlanId}`} className="p-2 hover:bg-gray-100 rounded text-blue-600">
                                <Edit className="w-4 h-4" />
                            </Link>
                            <button onClick={() => handleDelete(rp.ratePlanId)} className="p-2 hover:bg-red-50 rounded text-red-500">
                                <Trash className="w-4 h-4" />
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default AdminRatePlansPage;
