import { useSelector } from 'react-redux';
import { selectCurrentUser, selectIsAuthenticated, selectCurrentToken } from '../feature/auth/authSlice.js';

const useAuth = () => {
  const user = useSelector(selectCurrentUser);
  const isAuthenticated = useSelector(selectIsAuthenticated);
  const token = useSelector(selectCurrentToken);
  
  return { user, isAuthenticated, token };
};

export default useAuth;

