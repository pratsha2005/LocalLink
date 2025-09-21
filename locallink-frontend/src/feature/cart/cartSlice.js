import { createSlice } from '@reduxjs/toolkit';
import toast from 'react-hot-toast';

const initialState = {
  cartItems: localStorage.getItem('cartItems') 
    ? JSON.parse(localStorage.getItem('cartItems')) 
    : [],
};

const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    addToCart: (state, action) => {
      const newItem = action.payload;
      const existItem = state.cartItems.find(item => item.id === newItem.id);

      // Prevent adding items from different producers
      if (state.cartItems.length > 0 && state.cartItems[0].producerId !== newItem.producerId) {
        toast.error('You can only order from one producer at a time. Please clear your cart first.');
        return;
      }

      if (existItem) {
        toast.error('Item is already in your cart.');
        return;
      } else {
        state.cartItems.push(newItem);
        toast.success(`${newItem.name} added to cart!`);
      }
      localStorage.setItem('cartItems', JSON.stringify(state.cartItems));
    },
    removeFromCart: (state, action) => {
        state.cartItems = state.cartItems.filter(item => item.id !== action.payload.id);
        localStorage.setItem('cartItems', JSON.stringify(state.cartItems));
    },
    clearCart: (state) => {
        state.cartItems = [];
        localStorage.removeItem('cartItems');
    }
  },
});

export const { addToCart, removeFromCart, clearCart } = cartSlice.actions;
export default cartSlice.reducer;