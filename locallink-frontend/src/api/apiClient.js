import axios from 'axios';
import { store } from '../app/store';
import { logOut } from '../feature/auth/authSlice.js';

const apiClient = axios.create({
  baseURL: 'http://localhost:8080', // Your Go backend URL
});

// Request interceptor to add the auth token to headers
apiClient.interceptors.request.use((config) => {
  const token = store.getState().auth.token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor to handle expired tokens
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // If the token is expired or invalid, log the user out
      store.dispatch(logOut());
    }
    return Promise.reject(error);
  }
);


export default apiClient;

