import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { loginUser, reset } from '../features/authSlice';
import Button from '../components/common/Button.jsx';
import { Lock } from 'lucide-react';

const LoginPage = () => {
    const [formData, setFormData] = useState({
        username: '',
        password: '',
    });

    const { username, password } = formData;
    const navigate = useNavigate();
    const dispatch = useDispatch();

    const { user, isLoading, isError, isSuccess, message } = useSelector(
        (state) => state.auth
    );

    useEffect(() => {
        if (isError) {
            alert(message);
        }

        if (isSuccess || user) {
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
        const userData = {
            username,
            password,
        };
        dispatch(loginUser(userData));
    };

    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center -mt-20">
            <div className="bg-white p-8 rounded-lg shadow-xl max-w-md w-full">
                <div className="flex justify-center mb-6 text-primary">
                    <Lock className="w-12 h-12" />
                </div>
                <h1 className="text-3xl font-serif font-bold text-center text-gray-800 mb-6">Login</h1>

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
                        <Button type="submit" fullWidth disabled={isLoading}>
                            {isLoading ? 'Logging In...' : 'Login'}
                        </Button>
                    </div>
                </form>
                <div className="mt-6 text-center text-sm">
                    <span className="text-gray-500">Don't have an account? </span>
                    <Link to="/register" className="text-primary font-bold hover:underline">Register</Link>
                </div>
            </div>
        </div>
    );
};

export default LoginPage;
