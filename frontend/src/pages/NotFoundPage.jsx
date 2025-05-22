import { useNavigate } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

export default function NotFoundPage() {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate('/');
  };

  return (
    <div className="max-w-screen-xl mx-auto flex flex-col lg:flex-row gap-6 ">
      <Sidebar />
      <div className="flex-1">
        <div className="bg-white rounded-lg shadow-md p-8">
        <h1 className="text-4xl font-bold mb-4">404 — Страница не найдена</h1>
        <p className="text-gray-600 mb-8">Страница, которую вы ищете, не существует.</p>
        <button
          onClick={handleClick}
          className="bg-primary-500 w-full text-white text-lg p-2 rounded-lg hover:bg-primary-600 transition-colors"
        >
          Вернуться на главную
        </button>
      </div>
      </div>
    </div>
  );
}
