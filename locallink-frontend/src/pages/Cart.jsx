import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { removeFromCart, clearCart } from '../features/cart/cartSlice';
import { createOrder } from '../api';
import toast from 'react-hot-toast';

const Cart = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();
    const { cartItems } = useSelector((state) => state.cart);
    const { isAuthenticated } = useSelector((state) => state.auth);

    const total = cartItems.reduce((acc, item) => acc + item.price, 0);

    const handlePlaceOrder = async () => {
        if (!isAuthenticated) {
            toast.error('Please log in to place an order.');
            navigate('/login');
            return;
        }

        const orderData = {
            producerId: cartItems[0].producerId,
            items: cartItems.map(item => ({
                productId: item.id,
                quantity: 1 // Assuming quantity of 1 for now
            }))
        };
        
        const toastId = toast.loading('Placing your order...');
        try {
            await createOrder(orderData);
            toast.success('Order placed successfully!', { id: toastId });
            dispatch(clearCart());
            navigate('/orders');
        } catch (error) {
            console.error('Failed to place order:', error);
            const errorMessage = error.response?.data?.error || 'Failed to place order. Please try again.';
            toast.error(errorMessage, { id: toastId });
        }
    };

    return (
        <div className="p-6 bg-white rounded-lg shadow-md max-w-4xl mx-auto">
            <h1 className="text-3xl font-bold mb-6 text-gray-800">Your Shopping Cart</h1>
            {cartItems.length === 0 ? (
                <div className="text-center py-10">
                    <p className="text-gray-600 text-lg">Your cart is empty.</p>
                    <Link to="/" className="mt-4 inline-block btn-primary max-w-xs">
                        Find Local Products
                    </Link>
                </div>
            ) : (
                <div>
                    <div className="space-y-4">
                        {cartItems.map(item => (
                            <div key={item.id} className="flex justify-between items-center border-b pb-2">
                                <div>
                                    <h2 className="font-semibold text-lg">{item.name}</h2>
                                    <p className="text-gray-600">₹{item.price.toFixed(2)}</p>
                                </div>
                                <button onClick={() => dispatch(removeFromCart(item))} className="text-red-500 hover:text-red-700 font-medium">
                                    Remove
                                </button>
                            </div>
                        ))}
                    </div>
                    <div className="mt-6 text-right">
                        <h2 className="text-2xl font-bold">Total: ₹{total.toFixed(2)}</h2>
                        <button onClick={handlePlaceOrder} className="btn-primary mt-4 max-w-xs ml-auto">
                            Place Order
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Cart;