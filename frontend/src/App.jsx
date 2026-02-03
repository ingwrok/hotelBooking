import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout.jsx';
import HomePage from './pages/HomePage.jsx';
import SearchPage from './pages/SearchPage.jsx';
import RoomDetailsPage from './pages/RoomDetailsPage.jsx';
import BookingPage from './pages/BookingPage.jsx';
import ConfirmationPage from './pages/ConfirmationPage.jsx';

import LoginPage from './pages/LoginPage.jsx';
import RegisterPage from './pages/RegisterPage.jsx';
import MyHistoryPage from './pages/MyHistoryPage.jsx';

// Admin Components
import AdminLayout from './components/admin/AdminLayout.jsx';
import AdminDashboardPage from './pages/admin/AdminDashboardPage.jsx';
import AdminBookingsPage from './pages/admin/AdminBookingsPage.jsx';
import AdminRoomsPage from './pages/admin/AdminRoomsPage.jsx';
import AdminRoomTypesPage from './pages/admin/AdminRoomTypesPage';
import AdminAddonsPage from './pages/admin/AdminAddonsPage.jsx';
import AdminUsersPage from './pages/admin/AdminUsersPage';
import AdminRoomTypePage from './pages/AdminRoomTypePage'; // Creation Form
import AdminRatePlansPage from './pages/admin/AdminRatePlansPage';
import AdminRatePlanPage from './pages/AdminRatePlanPage';

import RequireAdmin from './components/admin/RequireAdmin.jsx';

// Placeholder for missing pages
const Placeholder = ({ title }) => <div className="p-10 text-center text-2xl">{title}</div>;

import ScrollToTop from './components/common/ScrollToTop.jsx';

const App = () => {
    return (
        <>
            <ScrollToTop />
            <Routes>
                {/* Main Application Routes */}
                <Route path="/" element={<Layout />}>
                    <Route index element={<HomePage />} />
                    <Route path="search" element={<SearchPage />} />
                    <Route path="room/:id" element={<RoomDetailsPage />} />

                    {/* Auth Routes */}
                    <Route path="login" element={<LoginPage />} />
                    <Route path="register" element={<RegisterPage />} />
                    <Route path="my-history" element={<MyHistoryPage />} />

                    {/* Booking Flow */}
                    <Route path="booking" element={<BookingPage />} />
                    <Route path="confirmation" element={<ConfirmationPage />} />

                    {/* Placeholders */}
                    <Route path="my-bookings" element={<Placeholder title="My Bookings" />} />
                    <Route path="about" element={<Placeholder title="About Us" />} />
                    <Route path="contact" element={<Placeholder title="Contact" />} />
                </Route>

                {/* Admin Routes (Separate Layout) */}
                <Route path="/admin" element={
                    <RequireAdmin>
                        <AdminLayout />
                    </RequireAdmin>
                }>
                    <Route index element={<Navigate to="dashboard" replace />} />
                    <Route path="dashboard" element={<AdminDashboardPage />} />
                    <Route path="bookings" element={<AdminBookingsPage />} />
                    <Route path="rooms" element={<AdminRoomsPage />} />

                    {/* Room Types */}
                    <Route path="room-types" element={<AdminRoomTypesPage />} />
                    <Route path="room-types/new" element={<AdminRoomTypePage />} />
                    <Route path="room-types/edit/:id" element={<AdminRoomTypePage />} />

                    {/* Rate Plans */}
                    <Route path="rate-plans" element={<AdminRatePlansPage />} />
                    <Route path="rate-plans/new" element={<AdminRatePlanPage />} />
                    <Route path="rate-plans/edit/:id" element={<AdminRatePlanPage />} />

                    {/* Users (New) */}
                    <Route path="users" element={<AdminUsersPage />} />

                    {/* Addons */}
                    <Route path="addons" element={<AdminAddonsPage />} />
                </Route>
            </Routes>
        </>
    );
};

export default App;
