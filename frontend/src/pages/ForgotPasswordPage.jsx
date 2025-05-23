import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';

// Страница запроса сброса пароля
export function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:80/api/auth/password_resets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Ошибка при запросе сброса пароля');
      }

      setMessage('Ссылка для сброса пароля отправлена на вашу почту');
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-start pt-20 justify-center p-4">
      <div className="max-w-md mx-auto bg-white shadow-lg rounded-lg p-6">
        <div className="space-y-4 max-w-xs">
          <div className="text-center pb-4">
            <h1 className="text-3xl font-bold text-primary-600">PollForge</h1>
            <h2 className="mt-2 text-gray-600">Восстановление пароля</h2>
          </div>

          {error && (
            <div className="p-3 bg-red-50 text-red-600 rounded-md text-sm">
              {error}
            </div>
          )}

          {message ? (
            <div className="p-3 bg-green-50 text-green-600 rounded-md text-sm">
              {message}
            </div>
          ) : (
            <>
              <p className="text-sm text-gray-600">
                Введите email, связанный с вашим аккаунтом, и мы отправим вам ссылку для сброса пароля.
              </p>

              <form className="space-y-6" onSubmit={handleSubmit}>
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
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="your@email.com"
                  />
                </div>

                <div className="pt-4">
                  <button
                    type="submit"
                    disabled={isLoading}
                    className={`w-full h-12 py-2 px-4 bg-white border-2 border-primary-600 text-primary-600 text-base font-medium rounded-full shadow-md hover:bg-primary-50 transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                      isLoading ? 'opacity-70 cursor-not-allowed' : ''
                    }`}
                  >
                    {isLoading ? 'Отправка...' : 'Отправить ссылку'}
                  </button>
                </div>
              </form>
            </>
          )}

          <div className="text-center text-sm text-gray-600 pt-6 border-t border-gray-200 mt-4">
            <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500">
              Вернуться к входу
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

// Страница установки нового пароля
export function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const [newPassword, setNewPassword] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!token) {
      setError('Недействительная ссылка для сброса пароля');
      return;
    }

    if (newPassword.length < 8) {
      setError('Пароль должен содержать минимум 8 символов');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:80/api/auth/password_resets/confirm', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token,
          new_password: newPassword
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Ошибка при сбросе пароля');
      }

      setMessage('Пароль успешно изменен. Теперь вы можете войти с новым паролем.');
      
      // Перенаправляем на страницу входа через 3 секунды
      setTimeout(() => {
        navigate('/login');
      }, 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-start pt-20 justify-center p-4">
      <div className="max-w-md mx-auto bg-white shadow-lg rounded-lg p-6">
        <div className="space-y-4 max-w-xs">
          <div className="text-center pb-4">
            <h1 className="text-3xl font-bold text-primary-600">PollForge</h1>
            <h2 className="mt-2 text-gray-600">Установка нового пароля</h2>
          </div>

          {error && (
            <div className="p-3 bg-red-50 text-red-600 rounded-md text-sm">
              {error}
            </div>
          )}

          {message ? (
            <div className="p-3 bg-green-50 text-green-600 rounded-md text-sm">
              {message}
            </div>
          ) : (
            <>
              <form className="space-y-6" onSubmit={handleSubmit}>
                <div className="space-y-2">
                  <label htmlFor="newPassword" className="block text-sm font-medium text-gray-700">
                    Новый пароль
                  </label>
                  <input
                    id="newPassword"
                    name="newPassword"
                    type="password"
                    required
                    minLength={8}
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="••••••••"
                  />
                </div>

                <div className="pt-4">
                  <button
                    type="submit"
                    disabled={isLoading}
                    className={`w-full h-12 py-2 px-4 bg-white border-2 border-primary-600 text-primary-600 text-base font-medium rounded-full shadow-md hover:bg-primary-50 transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                      isLoading ? 'opacity-70 cursor-not-allowed' : ''
                    }`}
                  >
                    {isLoading ? 'Сохранение...' : 'Установить новый пароль'}
                  </button>
                </div>
              </form>
            </>
          )}

          <div className="text-center text-sm text-gray-600 pt-6 border-t border-gray-200 mt-4">
            <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500">
              Вернуться к входу
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}