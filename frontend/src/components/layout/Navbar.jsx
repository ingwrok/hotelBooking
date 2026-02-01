import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Menu, Hotel, User, LogOut, History } from 'lucide-react';
import { useSelector, useDispatch } from 'react-redux';
import { logoutUser, reset } from '../../features/authSlice';

const Navbar = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const { user } = useSelector((state) => state.auth);

    const [isOpen, setIsOpen] = React.useState(false);

    const onLogout = () => {
        dispatch(logoutUser());
        dispatch(reset());
        navigate('/');
    };

    return (
        <nav className="bg-white border-b border-gray-100 sticky top-0 z-50">
            <div className="container mx-auto px-4 h-16 flex items-center justify-between">

                {/* Logo */}
                <Link to="/" className="flex items-center space-x-2 text-primary">
                    <Hotel className="w-8 h-8" />
                    <span className="text-xl font-serif font-bold tracking-tight">LuxStay</span>
                </Link>

                {/* Desktop Menu */}
                <div className="hidden md:flex items-center space-x-8">
                    <Link to="/" className="text-gray-600 hover:text-primary transition">Home</Link>
                    <Link to="/search" className="text-gray-600 hover:text-primary transition">Rooms</Link>

                    {user ? (
                        <div className="flex items-center space-x-6">
                            <Link to="/my-history" className="flex items-center text-gray-600 hover:text-primary transition">
                                <History className="w-4 h-4 mr-1" />
                                My History
                            </Link>
                            <div className="flex items-center space-x-2 text-primary font-medium">
                                <User className="w-4 h-4" />
                                <span>{user.user?.username || user.username}</span>
                            </div>
                            <button onClick={onLogout} className="flex items-center text-gray-500 hover:text-red-500 transition">
                                <LogOut className="w-4 h-4 mr-1" />
                                Logout
                            </button>
                        </div>
                    ) : (
                        <>
                            <Link to="/login" className="text-gray-600 hover:text-primary transition">Login</Link>
                            <Link to="/register" className="bg-primary text-white px-5 py-2 rounded-full hover:bg-primary-dark transition shadow-md shadow-primary/20">
                                Sign Up
                            </Link>
                        </>
                    )}
                </div>

                {/* Mobile Menu Button */}
                <button onClick={() => setIsOpen(!isOpen)} className="md:hidden text-gray-600 p-2">
                    <Menu className="w-6 h-6" />
                </button>
            </div>

            {/* Mobile Menu Dropdown */}
            {isOpen && (
                <div className="md:hidden bg-white border-t border-gray-100 absolute w-full left-0 shadow-lg py-4 px-4 flex flex-col space-y-4">
                    <Link to="/" onClick={() => setIsOpen(false)} className="text-gray-600 hover:text-primary py-2">Home</Link>
                    <Link to="/search" onClick={() => setIsOpen(false)} className="text-gray-600 hover:text-primary py-2">Rooms</Link>

                    {user ? (
                        <>
                            <Link to="/my-history" onClick={() => setIsOpen(false)} className="flex items-center text-gray-600 hover:text-primary py-2">
                                <History className="w-4 h-4 mr-2" />
                                My History
                            </Link>
                            <div className="flex items-center space-x-2 text-primary font-medium py-2 border-t border-gray-50 pt-3">
                                <User className="w-4 h-4" />
                                <span>{user.user?.username || user.username}</span>
                            </div>
                            <button onClick={onLogout} className="flex items-center text-gray-500 hover:text-red-500 py-2">
                                <LogOut className="w-4 h-4 mr-2" />
                                Logout
                            </button>
                        </>
                    ) : (
                        <div className="flex flex-col gap-3 pt-2">
                            <Link to="/login" onClick={() => setIsOpen(false)} className="text-center text-gray-600 hover:text-primary py-2 border border-gray-200 rounded-lg">Login</Link>
                            <Link to="/register" onClick={() => setIsOpen(false)} className="text-center bg-primary text-white px-5 py-2 rounded-lg hover:bg-primary-dark shadow-md">
                                Sign Up
                            </Link>
                        </div>
                    )}
                </div>
            )}
        </nav>
    );
};

export default Navbar;
