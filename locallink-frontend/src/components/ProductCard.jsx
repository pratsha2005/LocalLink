import React from 'react';
import { useDispatch } from 'react-redux';
import { addToCart } from '../features/cart/cartSlice';
import useAuth from '../hooks/useAuth';

const ProductCard = ({ product }) => {
  const dispatch = useDispatch();
  const { isAuthenticated } = useAuth();

  const handleAddToCart = () => {
    dispatch(addToCart(product));
  };

  return (
    <div className="w-64 p-1">
      <h3 className="text-lg font-bold text-gray-800">{product.name}</h3>
      <p className="text-sm text-gray-600 mt-1 truncate">{product.description}</p>
      <div className="mt-2 flex justify-between items-center">
        <p className="text-lg font-semibold text-green-600">₹{product.price.toFixed(2)}</p>
        <p className="text-sm text-gray-500">In Stock: {product.quantity}</p>
      </div>
      {isAuthenticated && (
          <button 
            onClick={handleAddToCart}
            className="mt-4 w-full bg-green-600 text-white py-2 rounded-md text-sm font-medium hover:bg-green-700 transition"
          >
            Add to Cart
          </button>
      )}
    </div>
  );
};

export default ProductCard;