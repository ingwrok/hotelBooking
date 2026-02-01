import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { LayoutDashboard, BedDouble, CalendarRange, Users, LogOut, Coffee, Tag, Percent, Menu, X } from 'lucide-react';
import { useDispatch } from 'react-redux';
import { logoutUser } from '../../features/authSlice';

const AdminLayout = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const handleLogout = () => {
        dispatch(logoutUser());
        navigate('/login');
    };

    const navItems = [
        { path: '/admin/dashboard', label: 'Dashboard', icon: LayoutDashboard },
        { path: '/admin/bookings', label: 'Bookings', icon: CalendarRange },
        { path: '/admin/rooms', label: 'Rooms', icon: BedDouble },
        { path: '/admin/room-types', label: 'Room Types', icon: Tag },
        { path: '/admin/rate-plans', label: 'Rate Plans', icon: Percent },
        { path: '/admin/addons', label: 'Addons', icon: Coffee },
        { path: '/admin/users', label: 'Users', icon: Users },
    ];

    const [isMobileMenuOpen, setIsMobileMenuOpen] = React.useState(false);

    const toggleMobileMenu = () => {
        setIsMobileMenuOpen(!isMobileMenuOpen);
    };

    return (
        <div className="flex min-h-screen bg-gray-100 font-sans overflow-x-hidden w-full">
            {/* Sidebar (Desktop) */}
            <aside className="hidden md:flex w-64 bg-white shadow-xl fixed h-full z-10 flex-col">
                <div className="p-8 border-b border-gray-100 flex items-center justify-center">
                    <h1 className="text-2xl font-serif font-bold text-primary tracking-widest uppercase text-center">
                        Admin
                    </h1>
                </div>

                <nav className="flex-1 overflow-y-auto py-6">
                    <ul className="space-y-1 px-4">
                        {navItems.map((item) => (
                            <li key={item.path}>
                                <NavLink
                                    to={item.path}
                                    className={({ isActive }) =>
                                        `flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors ${isActive
                                            ? 'bg-primary text-white shadow-md'
                                            : 'text-gray-500 hover:bg-gray-50 hover:text-primary'
                                        }`
                                    }
                                >
                                    <item.icon className="w-5 h-5" />
                                    {item.label}
                                </NavLink>
                            </li>
                        ))}
                    </ul>
                </nav>

                <div className="p-4 border-t border-gray-100">
                    <button
                        onClick={handleLogout}
                        className="flex items-center gap-3 px-4 py-3 w-full text-left rounded-lg text-sm font-medium text-red-500 hover:bg-red-50 transition-colors"
                    >
                        <LogOut className="w-5 h-5" />
                        Logout
                    </button>
                </div>
            </aside>

            {/* Mobile Sidebar (Overlay) */}
            {isMobileMenuOpen && (
                <div className="fixed inset-0 z-50 md:hidden flex">
                    <div className="fixed inset-0 bg-black/50" onClick={() => setIsMobileMenuOpen(false)}></div>
                    <aside className="w-64 bg-white shadow-xl h-full flex flex-col relative z-50 animate-slide-in-left">
                        <div className="p-6 border-b border-gray-100 flex items-center justify-between">
                            <h1 className="text-xl font-serif font-bold text-primary tracking-widest uppercase">
                                Admin
                            </h1>
                            <button onClick={() => setIsMobileMenuOpen(false)} className="text-gray-500">
                                <X className="w-6 h-6" />
                            </button>
                        </div>

                        <nav className="flex-1 overflow-y-auto py-6">
                            <ul className="space-y-1 px-4">
                                {navItems.map((item) => (
                                    <li key={item.path}>
                                        <NavLink
                                            to={item.path}
                                            onClick={() => setIsMobileMenuOpen(false)}
                                            className={({ isActive }) =>
                                                `flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors ${isActive
                                                    ? 'bg-primary text-white shadow-md'
                                                    : 'text-gray-500 hover:bg-gray-50 hover:text-primary'
                                                }`
                                            }
                                        >
                                            <item.icon className="w-5 h-5" />
                                            {item.label}
                                        </NavLink>
                                    </li>
                                ))}
                            </ul>
                        </nav>

                        <div className="p-4 border-t border-gray-100">
                            <button
                                onClick={handleLogout}
                                className="flex items-center gap-3 px-4 py-3 w-full text-left rounded-lg text-sm font-medium text-red-500 hover:bg-red-50 transition-colors"
                            >
                                <LogOut className="w-5 h-5" />
                                Logout
                            </button>
                        </div>
                    </aside>
                </div>
            )}

            {/* Mobile Header */}
            <header className="md:hidden fixed top-0 w-full bg-white border-b border-gray-100 z-10 px-4 h-16 flex items-center justify-between shadow-sm">
                <button onClick={toggleMobileMenu} className="text-gray-600 p-2">
                    <Menu className="w-6 h-6" />
                </button>
                <span className="text-lg font-serif font-bold text-primary tracking-widest uppercase">Admin</span>
                <div className="w-6"></div> {/* Spacer */}
            </header>

            {/* Main Content */}
            <main className="flex-1 ml-0 md:ml-64 p-4 md:p-8 pt-20 md:pt-8 min-w-0">
                <div className="max-w-7xl mx-auto">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default AdminLayout;
