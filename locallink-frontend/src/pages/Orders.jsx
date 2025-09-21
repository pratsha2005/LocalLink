import React, { useState, useEffect } from 'react';
import { fetchUserOrders } from '../api';

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const getOrders = async () => {
      try {
        const res = await fetchUserOrders();
        // Handle cases where backend returns null for no orders
        setOrders(res.data || []);
      } catch (error) {
        console.error("Failed to fetch orders:", error);
      } finally {
        setLoading(false);
      }
    };
    getOrders();
  }, []);

  if (loading) return <div className="text-center mt-10">Loading orders...</div>;

  return (
    <div className="p-6 bg-white rounded-lg shadow-md max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6 text-gray-800">My Orders</h1>
      {orders.length === 0 ? (
        <p className="text-gray-600">You have not placed or received any orders yet.</p>
      ) : (
        <div className="space-y-6">
          {orders.map(order => (
            <div key={order.id} className="border border-gray-200 p-4 rounded-lg">
              <div className="flex justify-between items-center mb-3">
                <h2 className="font-bold text-lg text-gray-700">Order #{order.id}</h2>
                <span className={`capitalize px-3 py-1 text-sm font-medium rounded-full 
                  ${order.status === 'completed' ? 'bg-green-100 text-green-800' : 
                   order.status === 'cancelled' ? 'bg-red-100 text-red-800' : 
                   'bg-blue-100 text-blue-800'}`}>
                  {order.status}
                </span>
              </div>
              <p className="text-gray-600">Total Price: <span className="font-semibold">₹{order.totalPrice.toFixed(2)}</span></p>
              <p className="text-gray-600 text-sm">Date: {new Date(order.createdAt).toLocaleDateString()}</p>
              <div className="mt-4 border-t pt-3">
                <h3 className="font-semibold text-gray-700">Items:</h3>
                <ul className="list-disc pl-5 mt-2 space-y-1 text-gray-600">
                  {order.items.map(item => (
                    <li key={item.id}>{item.quantity} x (Product ID: {item.productId}) @ ₹{item.price.toFixed(2)} each</li>
                  ))}
                </ul>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Orders;

