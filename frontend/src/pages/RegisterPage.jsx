import { Link } from 'react-router-dom';

export default function RegisterPage() {
  return (
    <div className="min-h-screen bg-gray-50 flex items-start pt-20 justify-center p-4">
      <div className=" max-w-md mx-auto bg-white shadow-lg rounded-lg p-6"> 
      <div className="space-y-4 max-w-xs">

        {/* Заголовок */}
        <div className="text-center pb-4">
          <h1 className="text-3xl font-bold text-primary-600">PollForge</h1>
          <h2 className="mt-2 text-gray-600">Создайте аккаунт</h2>
          
        </div>

        {/* Форма регистрации */}
        <form className="space-y-6">
          {/* Имя */}
          <div className="space-y-2">
            <label htmlFor="name" className="block text-sm font-medium text-gray-700">
              Имя
            </label>
            <input
              id="name"
              name="name"
              type="text"
              autoComplete="name"
              required
              className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              placeholder="Ваше имя"
            />
          </div>

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
              autoComplete="new-password"
              required
              className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              placeholder="••••••••"
            />
          </div>

          {/* Подтверждение пароля */}
          <div className="space-y-2">
            <label htmlFor="confirm-password" className="block text-sm font-medium text-gray-700">
              Подтвердите пароль
            </label>
            <input
              id="confirm-password"
              name="confirm-password"
              type="password"
              autoComplete="new-password"
              required
              className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              placeholder="••••••••"
            />
          </div>

          {/* Кнопка регистрации */}
          <div className="pt-4">
            <button
              type="submit"
              className="w-full h-12 py-2 px-4 bg-white border-2 border-primary-600 text-primary-600 text-base font-medium rounded-full shadow-md hover:bg-primary-50 transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            >
              Зарегистрироваться
            </button>
          </div>
        </form>

        <div className="text-center pb-4">
        <p className="mt-1 text-sm text-gray-600">
            Уже есть аккаунт?{' '}
            <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500">
              Войдите
            </Link>
          </p>
        </div>
        
      </div>
      </div>
    </div>
  );
}