import apiClient from './apiClient';

// Auth
export const registerUser = (userData) => apiClient.post('/register', userData);
export const loginUser = (credentials) => apiClient.post('/login', credentials);

// User Profile
export const fetchUserProfile = () => apiClient.get('/users/me');

// Products
export const fetchNearbyProducts = (lat, lon, radius) => 
  apiClient.get(`/products/nearby?lat=${lat}&lon=${lon}&radius=${radius}`);
export const createProduct = (productData) => apiClient.post('/products', productData);

// Orders
export const createOrder = (orderData) => apiClient.post('/orders', orderData);
export const fetchUserOrders = () => apiClient.get('/orders');

