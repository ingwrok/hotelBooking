import React, { useEffect, useState } from 'react';
import { Trash2, User, Shield, ShieldAlert } from 'lucide-react';
import api from '../../api';
import { useSelector } from 'react-redux';

const AdminUsersPage = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const { user: currentUser } = useSelector((state) => state.auth);

    const fetchUsers = async () => {
        try {
            const response = await api.get('/users');
            setUsers(response.data);
        } catch (error) {
            console.error("Failed to fetch users", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    const handleDelete = async (id) => {
        if (!window.confirm("Are you sure you want to delete this user? This action cannot be undone.")) {
            return;
        }

        try {
            await api.delete(`/users/${id}`);
            setUsers(users.filter(u => u.user_id !== id && u.id !== id)); // Handle both casing if dto inconsistent
            // Re-fetch to be safe or filter locally
            fetchUsers();
        } catch (error) {
            alert("Failed to delete user: " + (error.response?.data?.error || error.message));
        }
    };

    if (loading) return <div className="p-8 text-center text-gray-500">Loading users...</div>;

    return (
        <div>
            <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4 mb-8">
                <div>
                    <h1 className="text-2xl font-serif font-bold text-gray-800">User Management</h1>
                    <p className="text-gray-500 text-sm mt-1">View and manage registered users</p>
                </div>
                <div className="bg-blue-50 text-blue-700 px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2">
                    <User className="w-4 h-4" />
                    Total Users: {users.length}
                </div>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-left border-collapse">
                        <thead>
                            <tr className="bg-gray-50 text-gray-500 text-xs uppercase tracking-wider border-b border-gray-100">
                                <th className="p-4 font-medium">ID</th>
                                <th className="p-4 font-medium">Username</th>
                                <th className="p-4 font-medium">Email</th>
                                <th className="p-4 font-medium">Role</th>
                                <th className="p-4 font-medium text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-50">
                            {users.map((u) => {
                                // normalize id from dto
                                const userId = u.user_id || u.id;
                                const isMe = currentUser?.user?.id === userId || currentUser?.id === userId;

                                return (
                                    <tr key={userId} className="hover:bg-gray-50/50 transition-colors">
                                        <td className="p-4 text-gray-400 font-mono text-xs">#{userId}</td>
                                        <td className="p-4 font-medium text-gray-700">{u.username}</td>
                                        <td className="p-4 text-gray-500">{u.email}</td>
                                        <td className="p-4">
                                            {u.is_admin ? (
                                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-bold bg-purple-100 text-purple-700">
                                                    <ShieldAlert className="w-3 h-3" /> Admin
                                                </span>
                                            ) : (
                                                <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-bold bg-gray-100 text-gray-600">
                                                    <User className="w-3 h-3" /> User
                                                </span>
                                            )}
                                        </td>
                                        <td className="p-4 text-right">
                                            {isMe ? (
                                                <span className="text-xs text-gray-300 italic">Current User</span>
                                            ) : (
                                                <button
                                                    onClick={() => handleDelete(userId)}
                                                    className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                                                    title="Delete User"
                                                >
                                                    <Trash2 className="w-4 h-4" />
                                                </button>
                                            )}
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>

                {users.length === 0 && (
                    <div className="p-12 text-center text-gray-400">
                        No users found.
                    </div>
                )}
            </div>
        </div>
    );
};

export default AdminUsersPage;
