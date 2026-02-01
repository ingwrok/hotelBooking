import React, { useEffect, useState } from 'react';
import api from '../../api';
import { Plus, Trash, Search, Coffee, Pencil } from 'lucide-react';

const AdminAddonsPage = () => {
    const [addons, setAddons] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isAdding, setIsAdding] = useState(false);
    const [editingId, setEditingId] = useState(null);

    // New Addon Form State
    const [newAddon, setNewAddon] = useState({
        categoryId: '',
        name: '',
        description: '',
        price: '',
        unitName: '',
        pictureUrl: ''
    });
    const [uploading, setUploading] = useState(false);

    const fetchData = async () => {
        try {
            setLoading(true);
            const [addonsRes, catsRes] = await Promise.all([
                api.get('/addons'),
                api.get('/addon-categories')
            ]);
            setAddons(addonsRes.data);
            setCategories(catsRes.data);

            // Set default category if available
            if (catsRes.data.length > 0) {
                setNewAddon(prev => ({ ...prev, categoryId: catsRes.data[0].categoryId || catsRes.data[0].addonCategoryId }));
            }
        } catch (e) {
            console.error("Failed to fetch addons", e);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const handleAddAddon = async (e) => {
        e.preventDefault();

        // Validation
        if (!newAddon.name || !newAddon.price || !newAddon.categoryId) {
            alert('Please fill in all required fields');
            return;
        }

        try {
            const payload = {
                ...newAddon,
                price: parseFloat(newAddon.price),
                categoryId: parseInt(newAddon.categoryId)
            };

            if (editingId) {
                await api.put(`/addons/${editingId}`, payload);
            } else {
                await api.post('/addons', payload);
            }

            setIsAdding(false);
            setEditingId(null);
            setNewAddon({
                categoryId: categories[0]?.categoryId || categories[0]?.addonCategoryId || '',
                name: '',
                description: '',
                price: '',
                unitName: '',
                pictureUrl: ''
            });
            fetchData();
        } catch (e) {
            alert(editingId ? 'Failed to update addon.' : 'Failed to add addon.');
            console.error(e);
        }
    };

    const handleEdit = (addon) => {
        setEditingId(addon.addonId);
        setNewAddon({
            categoryId: addon.categoryId,
            name: addon.name,
            description: addon.description,
            price: addon.price,
            unitName: addon.unitName,
            pictureUrl: addon.pictureUrl || ''
        });
        setIsAdding(true);
    };

    const handleDelete = async (id) => {
        if (!window.confirm('Delete this addon?')) return;
        try {
            await api.delete(`/addons/${id}`);
            fetchData();
        } catch (e) {
            alert('Failed to delete addon.');
        }
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setNewAddon(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleFileChange = async (e) => {
        if (e.target.files && e.target.files[0]) {
            setUploading(true);
            const formData = new FormData();
            formData.append('image', e.target.files[0]);

            try {
                const res = await api.post('/addons/upload', formData, {
                    headers: { 'Content-Type': 'multipart/form-data' }
                });
                setNewAddon(prev => ({ ...prev, pictureUrl: res.data.url }));
            } catch (error) {
                console.error("Upload failed", error);
                alert("Failed to upload image");
            } finally {
                setUploading(false);
            }
        }
    };

    if (loading) return <div className="p-8 text-center text-gray-500">Loading Addons...</div>;

    return (
        <div>
            <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4 mb-6">
                <div className="flex items-center gap-3">
                    <div className="bg-orange-100 p-2 rounded-lg">
                        <Coffee className="w-6 h-6 text-orange-600" />
                    </div>
                    <div>
                        <h2 className="text-2xl font-bold text-gray-800">Addon Management</h2>
                        <p className="text-gray-500 text-sm">Manage extra services and items</p>
                    </div>
                </div>

                <button
                    onClick={() => {
                        setIsAdding(!isAdding);
                        setEditingId(null);
                        setNewAddon({
                            categoryId: categories[0]?.categoryId || categories[0]?.addonCategoryId || '',
                            name: '',
                            description: '',
                            price: '',
                            unitName: '',
                            pictureUrl: ''
                        });
                    }}
                    className="bg-primary hover:bg-gray-800 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors shadow-sm"
                >
                    <Plus className="w-4 h-4" />
                    {isAdding ? 'Cancel' : 'Add New Addon'}
                </button>
            </div>

            {isAdding && (
                <div className="bg-white p-6 rounded-xl shadow-md border border-gray-100 mb-8 animate-fade-in-down">
                    <h4 className="font-bold text-lg mb-4 text-gray-800 border-b pb-2">{editingId ? 'Edit Addon' : 'Add New Addon'}</h4>
                    <form onSubmit={handleAddAddon} className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">

                        <div>
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Name *</label>
                            <input
                                type="text"
                                name="name"
                                value={newAddon.name}
                                onChange={handleChange}
                                className="border border-gray-300 p-2.5 rounded-lg w-full focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all outline-none"
                                placeholder="e.g. Airport Pickup"
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Category *</label>
                            <select
                                name="categoryId"
                                value={newAddon.categoryId}
                                onChange={handleChange}
                                className="border border-gray-300 p-2.5 rounded-lg w-full bg-white focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all outline-none"
                                required
                            >
                                {categories.map(cat => (
                                    <option key={cat.categoryId || cat.addonCategoryId} value={cat.categoryId || cat.addonCategoryId}>
                                        {cat.name}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Price (THB) *</label>
                            <input
                                type="number"
                                name="price"
                                value={newAddon.price}
                                onChange={handleChange}
                                className="border border-gray-300 p-2.5 rounded-lg w-full focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all outline-none"
                                placeholder="0.00"
                                min="0"
                                required
                            />
                        </div>

                        <div className="md:col-span-2">
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Description</label>
                            <input
                                type="text"
                                name="description"
                                value={newAddon.description}
                                onChange={handleChange}
                                className="border border-gray-300 p-2.5 rounded-lg w-full focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all outline-none"
                                placeholder="Short description of the service"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Unit Name</label>
                            <input
                                type="text"
                                name="unitName"
                                value={newAddon.unitName}
                                onChange={handleChange}
                                className="border border-gray-300 p-2.5 rounded-lg w-full focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all outline-none"
                                placeholder="e.g. trip, set, person"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-bold uppercase text-gray-500 mb-1">Image</label>
                            <div className="flex gap-2 items-center">
                                <input
                                    type="file"
                                    accept="image/*"
                                    onChange={handleFileChange}
                                    className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary file:text-white hover:file:bg-opacity-90"
                                    disabled={uploading}
                                />
                                {uploading && <span className="text-xs text-yellow-600">Uploading...</span>}
                            </div>
                            {newAddon.pictureUrl && (
                                <div className="mt-2 text-center">
                                    <img src={newAddon.pictureUrl} alt="Preview" className="h-20 w-auto object-cover rounded shadow-sm mx-auto" />
                                </div>
                            )}
                        </div>

                        <div className="md:col-span-3 flex justify-end gap-3 mt-2">
                            <button
                                type="button"
                                onClick={() => { setIsAdding(false); setEditingId(null); }}
                                className="px-5 py-2.5 rounded-lg border border-gray-300 font-medium hover:bg-gray-50 text-gray-600 transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="bg-primary text-white px-6 py-2.5 rounded-lg font-bold hover:bg-gray-800 transition-transform active:scale-95 shadow-lg shadow-gray-200"
                            >
                                {editingId ? 'Update Addon' : 'Save Addon'}
                            </button>
                        </div>
                    </form>
                </div>
            )}

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                {addons.length === 0 ? (
                    <div className="p-12 text-center text-gray-400 bg-gray-50">
                        <Coffee className="w-12 h-12 mx-auto mb-3 opacity-30" />
                        <p>No addons found. Create one to get started.</p>
                    </div>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm text-left">
                            <thead className="bg-gray-50 text-gray-500 font-bold uppercase text-xs">
                                <tr>
                                    <th className="px-6 py-4">Image</th>
                                    <th className="px-6 py-4">Name</th>
                                    <th className="px-6 py-4">Category</th>
                                    <th className="px-6 py-4">Price</th>
                                    <th className="px-6 py-4">Unit</th>
                                    <th className="px-6 py-4 text-center">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100">
                                {addons.map(addon => {
                                    const cat = categories.find(c => (c.categoryId || c.addonCategoryId) === addon.categoryId);
                                    return (
                                        <tr key={addon.addonId} className="hover:bg-gray-50 transition-colors group">
                                            <td className="px-6 py-4">
                                                {addon.pictureUrl ? (
                                                    <img src={addon.pictureUrl} alt={addon.name} className="h-12 w-12 object-cover rounded shadow-sm" />
                                                ) : (
                                                    <div className="h-12 w-12 bg-gray-100 rounded flex items-center justify-center text-gray-300">
                                                        <Coffee className="w-6 h-6" />
                                                    </div>
                                                )}
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="font-bold text-gray-800">{addon.name}</div>
                                                <div className="text-xs text-gray-500 truncate max-w-[200px]">{addon.description}</div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <span className="bg-blue-50 text-blue-700 px-2.5 py-1 rounded text-xs font-semibold">
                                                    {cat ? cat.name : `ID: ${addon.categoryId}`}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 font-serif font-bold text-gray-800">
                                                {addon.price.toLocaleString()} THB
                                            </td>
                                            <td className="px-6 py-4 text-gray-500">
                                                per {addon.unitName || 'unit'}
                                            </td>
                                            <td className="px-6 py-4 text-center">
                                                <div className="flex justify-center gap-2 opacity-100 md:opacity-0 md:group-hover:opacity-100 transition-opacity">
                                                    <button
                                                        onClick={() => handleEdit(addon)}
                                                        className="text-gray-400 hover:text-blue-500 p-2 transition-colors"
                                                        title="Edit Addon"
                                                    >
                                                        <Pencil className="w-4 h-4" />
                                                    </button>
                                                    <button
                                                        onClick={() => handleDelete(addon.addonId)}
                                                        className="text-gray-400 hover:text-red-500 p-2 transition-colors"
                                                        title="Delete Addon"
                                                    >
                                                        <Trash className="w-4 h-4" />
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
};

export default AdminAddonsPage;
