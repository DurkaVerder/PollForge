import { Link, useNavigate } from 'react-router-dom';
import { useState } from 'react';

export default function LoginPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:80/api/auth/logging', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: formData.email,
          password: formData.password
        })
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Ошибка при входе');
      }

      // Сохраняем токен (например, в localStorage)
      localStorage.setItem('authToken', data.token);
      localStorage.setItem('userId', data.id);
      // Перенаправляем пользователя
      navigate('/stream-line');
    } catch (err) {
      setError(err.message);
      console.error('Login error:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-start pt-20 justify-center p-4">
      <div className="max-w-md mx-auto bg-white shadow-lg rounded-lg p-6"> 
        <div className="space-y-4 max-w-xs">
          {/* Заголовок */}
          <div className="text-center pb-4">
            <h1 className="text-3xl font-bold text-primary-600">PollForge</h1>
            <h2 className="mt-2 text-gray-600">Вход в аккаунт</h2>
          </div>

          {/* Отображение ошибки */}
          {error && (
            <div className="p-3 bg-red-50 text-red-600 rounded-md text-sm">
              {error}
            </div>
          )}

          {/* Форма входа */}
          <form className="space-y-6" onSubmit={handleSubmit}>
            {/* Email */}
            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                Email
              </label>
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                required
                value={formData.email}
                onChange={handleChange}
                className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                placeholder="your@email.com"
              />
            </div>

            {/* Пароль */}
            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                Пароль
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                value={formData.password}
                onChange={handleChange}
                className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                placeholder="••••••••"
              />
            </div>

            {/* Кнопка входа */}
            <div className="pt-4">
              <button
                type="submit"
                disabled={isLoading}
                className={`w-full h-12 py-2 px-4 bg-white border-2 border-primary-600 text-primary-600 text-base font-medium rounded-full shadow-md hover:bg-primary-50 transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                  isLoading ? 'opacity-70 cursor-not-allowed' : ''
                }`}
              >
                {isLoading ? 'Вход...' : 'Войти'}
              </button>
            </div>
          </form>

          {/* Ссылка "Забыли пароль" */}
          <div className="pt-3 text-center">
            <Link to="/forgot-password" className="text-sm text-primary-600 hover:text-primary-500">
              Забыли пароль?
            </Link>
          </div>

          {/* Ссылка на регистрацию */}
          <div className="text-center text-sm text-gray-600 pt-6 border-t border-gray-200 mt-4">
            Нет аккаунта?{' '}
            <Link to="/register" className="font-medium text-primary-600 hover:text-primary-500">
              Зарегистрироваться
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}