import { useEffect } from 'react';
import toast from 'react-hot-toast';
import useAuth from '../hooks/useAuth';

const NotificationHandler = () => {
  const { isAuthenticated, token } = useAuth();

  useEffect(() => {
    if (!isAuthenticated) {
      return;
    }

    const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    ws.onopen = () => console.log('WebSocket Connected');
    ws.onclose = () => console.log('WebSocket Disconnected');

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        console.log('WebSocket Message Received:', message);

        if (message.type === 'order_update') {
          toast.success(`Order #${message.orderId} status is now: ${message.status}`);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };
    
    // Cleanup on component unmount or when auth status changes
    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
  }, [isAuthenticated, token]);

  return null; // This component does not render anything
};

export default NotificationHandler;

