import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { setCredentials } from '../feature/auth/authSlice.js';
import { loginUser } from '../api';
import toast from 'react-hot-toast';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const toastId = toast.loading('Logging in...');
    try {
      const res = await loginUser({ email, password });
      dispatch(setCredentials({ token: res.data.token }));
      toast.success('Logged in successfully!', { id: toastId });
      navigate('/dashboard');
    } catch (error) {
      console.error('Login failed:', error);
      toast.error('Login failed. Please check your credentials.', { id: toastId });
    }
  };

  return (
    <div className="flex justify-center items-start mt-10">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center text-gray-800">Login to LocalLink</h1>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="label">Email Address</label>
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required className="input" placeholder="you@example.com" />
          </div>
          <div>
            <label className="label">Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required className="input" placeholder="••••••••" />
          </div>
          <button type="submit" className="btn-primary">Login</button>
        </form>
        <p className="text-center text-sm text-gray-600">Don't have an account? <Link to="/register" className="font-medium text-green-600 hover:underline">Register here</Link></p>
      </div>
    </div>
  );
};

export default Login;

