import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';

const RequireAdmin = ({ children }) => {
    const { user } = useSelector((state) => state.auth);
    const location = useLocation();

    if (!user) {
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    const isAdmin = user?.user?.is_admin || user?.isAdmin || false; // Handle both structures just in case

    if (!isAdmin) {
        return <Navigate to="/" replace />;
    }

    return children;
};

export default RequireAdmin;
