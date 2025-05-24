
import { Navigate } from 'react-router-dom';

const ProtectedRoute = ({ children }) => {
  const token = localStorage.getItem('authToken');

  if (!token) {
    // Если токен отсутствует — редирект на страницу логина
    return <Navigate to="/login" replace />;
  }

  return children;
};

export default ProtectedRoute;
