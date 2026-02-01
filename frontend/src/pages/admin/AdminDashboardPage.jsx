import React, { useEffect, useState } from 'react';
import api from '../../api';

// Simple Dashboard Overview
const AdminDashboardPage = () => {
    const [stats, setStats] = useState({
        totalBookings: 0,
        todayBookings: 0,
        revenue: 0,
        pending: 0
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                const res = await api.get('/bookings/all');
                const bookings = res.data;

                // Calculate basic stats
                const total = bookings.length;
                const revenue = bookings.reduce((sum, b) => sum + (b.totalPrice || 0), 0);
                const pending = bookings.filter(b => b.status === 'pending').length;
                const today = bookings.filter(b => {
                    const d = new Date(b.createdAt);
                    const now = new Date();
                    return d.getDate() === now.getDate() && d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear();
                }).length;

                setStats({ totalBookings: total, todayBookings: today, revenue, pending });
            } catch (error) {
                console.error("Failed to load dashboard stats", error);
            } finally {
                setLoading(false);
            }
        };

        fetchStats();
    }, []);

    if (loading) return <div>Loading Stats...</div>;

    const cards = [
        { label: 'Total Revenue', value: stats.revenue.toLocaleString() + ' THB', color: 'bg-green-50 text-green-700' },
        { label: 'Total Bookings', value: stats.totalBookings, color: 'bg-blue-50 text-blue-700' },
        { label: 'Pending Bookings', value: stats.pending, color: 'bg-orange-50 text-orange-700' },
        { label: 'Bookings Today', value: stats.todayBookings, color: 'bg-purple-50 text-purple-700' },
    ];

    return (
        <div>
            <h2 className="text-2xl font-bold mb-6 text-gray-800">Dashboard Overview</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {cards.map((card, i) => (
                    <div key={i} className={`p-6 rounded-xl shadow-sm border border-gray-100 ${card.color}`}>
                        <div className="text-sm font-bold uppercase tracking-wider mb-2 opacity-70">{card.label}</div>
                        <div className="text-3xl font-serif font-bold">{card.value}</div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default AdminDashboardPage;
