import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { registerUser } from '../api';
import toast from 'react-hot-toast';

const Register = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('buyer');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    const toastId = toast.loading('Registering...');
    try {
      await registerUser({ name, email, password, role });
      toast.success('Registration successful! Please log in.', { id: toastId });
      navigate('/login');
    } catch (error) {
      console.error('Registration failed:', error);
      toast.error('Registration failed. Please try again.', { id: toastId });
    }
  };

  return (
    <div className="flex justify-center items-start mt-10">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center text-gray-800">Create an Account</h1>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="label">Full Name</label>
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} required className="input" placeholder="Your Name" />
          </div>
          <div>
            <label className="label">Email Address</label>
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required className="input" placeholder="you@example.com" />
          </div>
          <div>
            <label className="label">Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required className="input" placeholder="••••••••" />
          </div>
          <div>
            <label className="label">I want to register as a...</label>
            <select value={role} onChange={(e) => setRole(e.target.value)} className="input">
              <option value="buyer">Buyer</option>
              <option value="producer">Producer</option>
            </select>
          </div>
          <button type="submit" className="btn-primary">Register</button>
        </form>
         <p className="text-center text-sm text-gray-600">Already have an account? <Link to="/login" className="font-medium text-green-600 hover:underline">Login here</Link></p>
      </div>
    </div>
  );
};

export default Register;

