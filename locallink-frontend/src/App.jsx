import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import AddProduct from './pages/AddProduct';
import Orders from './pages/Orders';
import ProtectedRoute from './components/ProtectedRoute';
import NotificationHandler from './components/NotificationHandler';
import { Toaster } from 'react-hot-toast';
import Cart from './pages/Cart'; // <-- IMPORT

function App() {
  return (
    <>
      <Toaster position="top-right" reverseOrder={false} />
      <NotificationHandler />
      <Routes>
        <Route path="/" element={<Layout />}>
          {/* Public routes */}
          <Route index element={<Home />} />
          <Route path="login" element={<Login />} />
          <Route path="register" element={<Register />} />
          <Route path="cart" element={<Cart />} /> {/* <-- ADD THIS ROUTE */}

          {/* Protected routes */}
          <Route element={<ProtectedRoute />}>
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="add-product" element={<AddProduct />} />
            <Route path="orders" element={<Orders />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;