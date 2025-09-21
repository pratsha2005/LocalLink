import React, { useState, useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { fetchNearbyProducts } from '../api';
import ProductCard from '../components/ProductCard';

const Home = () => {
  const [position, setPosition] = useState([25.18, 75.83]); // Default to Ranpur, Rajasthan
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        setPosition([latitude, longitude]);
      },
      () => {
        console.warn("Could not get location, using default.");
      },
      { enableHighAccuracy: true }
    );
  }, []);

  useEffect(() => {
    const getProducts = async () => {
      try {
        setLoading(true);
        const res = await fetchNearbyProducts(position[0], position[1], 20000); // 20km radius
        setProducts(res.data);
      } catch (error) {
        console.error("Failed to fetch products:", error);
      } finally {
        setLoading(false);
      }
    };
    if (position) {
      getProducts();
    }
  }, [position]);

  return (
    <div className="h-[calc(100vh-120px)] w-full rounded-lg shadow-lg overflow-hidden relative">
       <MapContainer center={position} zoom={13} scrollWheelZoom={true} style={{ height: '100%', width: '100%' }}>
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {products.map(product => (
          <Marker key={product.id} position={[product.latitude, product.longitude]}>
            <Popup minWidth={260}>
              <ProductCard product={product} />
            </Popup>
          </Marker>
        ))}
      </MapContainer>
      {loading && <div className="absolute top-4 right-4 z-[1000] bg-white p-3 rounded-lg shadow-xl text-gray-700">Loading products...</div>}
    </div>
  );
};

export default Home;

