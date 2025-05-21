import { useNavigate } from 'react-router-dom';

export default function NotFoundPage() {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate('/');
  };

  return (
    <div className="min-h-screen flex justify-center ">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">404 — Страница не найдена</h1>
        <p className="text-gray-600 mb-8">Страница, которую вы ищете, не существует.</p>
        <button
          onClick={handleClick}
          className="bg-primary-500 w-full text-white text-lg px-8 py-4 rounded-lg hover:bg-primary-600 transition-colors"
        >
          Вернуться на главную
        </button>
      </div>
    </div>
  );
}
