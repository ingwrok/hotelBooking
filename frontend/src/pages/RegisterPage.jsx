import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { registerUser, reset } from '../features/authSlice';
import Button from '../components/common/Button.jsx';
import { UserPlus } from 'lucide-react';

const RegisterPage = () => {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
    });

    const { username, email, password, confirmPassword } = formData;
    const navigate = useNavigate();
    const dispatch = useDispatch();

    const { user, isLoading, isError, isSuccess, message } = useSelector(
        (state) => state.auth
    );

    useEffect(() => {
        if (isError) {
            alert(message);
        }

        if (isSuccess) {
            alert("Registration successful! Please login.");
            navigate('/login');
        }

        if (user) {
            navigate('/');
        }

        dispatch(reset());
    }, [user, isError, isSuccess, message, navigate, dispatch]);

    const onChange = (e) => {
        setFormData((prevState) => ({
            ...prevState,
            [e.target.name]: e.target.value,
        }));
    };

    const onSubmit = (e) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            alert('Passwords do not match');
        } else {
            const userData = {
                username,
                email,
                password,
            };
            dispatch(registerUser(userData));
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center -mt-20">
            <div className="bg-white p-8 rounded-lg shadow-xl max-w-md w-full">
                <div className="flex justify-center mb-6 text-primary">
                    <UserPlus className="w-12 h-12" />
                </div>
                <h1 className="text-3xl font-serif font-bold text-center text-gray-800 mb-6">Create Account</h1>

                <form onSubmit={onSubmit} className="space-y-6">
                    <div>
                        <input
                            type="text"
                            className="w-full px-4 py-3 border border-gray-300 rounded focus:outline-none focus:border-primary"
                            id="username"
                            name="username"
                            value={username}
                            placeholder="Username"
                            onChange={onChange}
                        />
                    </div>
                    <div>
                        <input
                            type="email"
                            className="w-full px-4 py-3 border border-gray-300 rounded focus:outline-none focus:border-primary"
                            id="email"
                            name="email"
                            value={email}
                            placeholder="Email"
                            onChange={onChange}
                        />
                    </div>
                    <div>
                        <input
                            type="password"
                            className="w-full px-4 py-3 border border-gray-300 rounded focus:outline-none focus:border-primary"
                            id="password"
                            name="password"
                            value={password}
                            placeholder="Password"
                            onChange={onChange}
                        />
                    </div>
                    <div>
                        <input
                            type="password"
                            className="w-full px-4 py-3 border border-gray-300 rounded focus:outline-none focus:border-primary"
                            id="confirmPassword"
                            name="confirmPassword"
                            value={confirmPassword}
                            placeholder="Confirm Password"
                            onChange={onChange}
                        />
                    </div>
                    <div>
                        <Button type="submit" fullWidth disabled={isLoading}>
                            {isLoading ? 'Creating Account...' : 'Register'}
                        </Button>
                    </div>
                </form>
                <div className="mt-6 text-center text-sm">
                    <span className="text-gray-500">Already have an account? </span>
                    <Link to="/login" className="text-primary font-bold hover:underline">Login</Link>
                </div>
            </div>
        </div>
    );
};

export default RegisterPage;
