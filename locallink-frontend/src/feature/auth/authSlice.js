import { createSlice } from '@reduxjs/toolkit';
import { jwtDecode } from 'jwt-decode';

const token = localStorage.getItem('authToken');

let initialState;
try {
  initialState = {
    token: token || null,
    isAuthenticated: !!token,
    user: token ? jwtDecode(token) : null,
  };
} catch (error) {
  // Handle invalid token in localStorage
  console.error("Invalid token found in localStorage", error);
  localStorage.removeItem('authToken');
  initialState = {
    token: null,
    isAuthenticated: false,
    user: null,
  };
}

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials: (state, action) => {
      const { token } = action.payload;
      const decoded = jwtDecode(token);
      
      state.token = token;
      // Extracting userID from your Go JWT payload
      state.user = { id: decoded.userID, exp: decoded.exp }; 
      state.isAuthenticated = true;

      localStorage.setItem('authToken', token);
    },
    logOut: (state) => {
      state.token = null;
      state.user = null;
      state.isAuthenticated = false;
      localStorage.removeItem('authToken');
    },
  },
});

export const { setCredentials, logOut } = authSlice.actions;

export default authSlice.reducer;

export const selectCurrentUser = (state) => state.auth.user;
export const selectIsAuthenticated = (state) => state.auth.isAuthenticated;
export const selectCurrentToken = (state) => state.auth.token;

