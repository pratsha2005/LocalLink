import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createProduct } from '../api';
import toast from 'react-hot-toast';

const AddProduct = () => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [price, setPrice] = useState('');
    const [quantity, setQuantity] = useState('');
    // For simplicity, we'll use a fixed location. A real app would use a map picker.
    const [latitude] = useState(25.18);
    const [longitude] = useState(75.83);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        const productData = {
            name,
            description,
            price: parseFloat(price),
            quantity: parseInt(quantity, 10),
            latitude,
            longitude
        };
        const toastId = toast.loading('Adding product...');
        try {
            await createProduct(productData);
            toast.success('Product added successfully!', { id: toastId });
            navigate('/dashboard');
        } catch (error) {
            console.error('Failed to add product:', error);
            toast.error('Failed to add product.', { id: toastId });
        }
    };

    return (
        <div className="flex justify-center items-start mt-10">
            <div className="w-full max-w-lg p-8 space-y-6 bg-white rounded-lg shadow-md">
                <h1 className="text-2xl font-bold text-center text-gray-800">Add a New Product</h1>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="label">Product Name</label>
                        <input type="text" value={name} onChange={(e) => setName(e.target.value)} required className="input" />
                    </div>
                    <div>
                        <label className="label">Description</label>
                        <textarea value={description} onChange={(e) => setDescription(e.target.value)} required className="input h-24" placeholder="Describe your product..."></textarea>
                    </div>
                    <div>
                        <label className="label">Price (₹)</label>
                        <input type="number" step="0.01" value={price} onChange={(e) => setPrice(e.target.value)} required className="input" />
                    </div>
                    <div>
                        <label className="label">Quantity Available</label>
                        <input type="number" value={quantity} onChange={(e) => setQuantity(e.target.value)} required className="input" />
                    </div>
                    <button type="submit" className="btn-primary">Add Product</button>
                </form>
            </div>
        </div>
    );
};

export default AddProduct;

