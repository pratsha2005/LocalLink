import React, { useState, useEffect } from 'react';
import { fetchUserProfile } from '../api';
import { Link } from 'react-router-dom';

const Dashboard = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const getProfile = async () => {
      try {
        const res = await fetchUserProfile();
        setProfile(res.data);
      } catch (error) {
        console.error("Failed to fetch profile", error);
      } finally {
        setLoading(false);
      }
    };
    getProfile();
  }, []);

  if (loading) return <div className="text-center mt-10">Loading dashboard...</div>;
  if (!profile) return <div className="text-center mt-10 text-red-500">Could not load profile.</div>;


  return (
    <div className="p-6 bg-white rounded-lg shadow-md max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold mb-4 text-gray-800">Welcome, {profile.name}!</h1>
      <div className="space-y-2">
        <p><span className="font-semibold">Email:</span> {profile.email}</p>
        <p><span className="font-semibold">Role:</span> <span className="capitalize px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">{profile.role}</span></p>
      </div>
      
      {profile.role === 'producer' && (
        <div className="mt-8 border-t pt-6">
          <h2 className="text-xl font-semibold mb-3">Producer Actions</h2>
          <Link to="/add-product" className="btn-primary max-w-xs">
            + Add New Product
          </Link>
        </div>
      )}
    </div>
  );
};

export default Dashboard;

