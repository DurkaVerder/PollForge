import { Link } from 'react-router-dom';

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-gray-50 flex items-start pt-20 justify-center p-4">
     <div className=" max-w-md mx-auto bg-white shadow-lg rounded-lg p-6"> 
      <div className="space-y-4 max-w-xs"> {/* Увеличено space-y и добавлена максимальная ширина */}

        {/* Заголовок с увеличенными отступами */}
        <div className="text-center pb-4"> {/* Добавлен отступ снизу */}
          <h1 className="text-3xl font-bold text-primary-600">PollForge</h1>
          <h2 className="mt-2 text-gray-600">Вход в аккаунт</h2> {/* Увеличен mt */}
        </div>

        {/* Форма входа с увеличенными отступами */}
        <form className="space-y-6"> {/* Увеличено space-y */}
          {/* Email */}
          <div className="space-y-2"> {/* Добавлен внутренний отступ */}
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
          <div className="space-y-2"> {/* Добавлен внутренний отступ и отступ сверху */}
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
              Пароль
            </label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              required
              className="w-full h-10 px-4 text-base rounded-full border border-gray-300 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              placeholder="••••••••"
            />
          </div>

          {/* Кнопка с увеличенными отступами */}
          <div className="pt-4"> {/* Большой отступ сверху */}
            <button
              type="submit"
              className="w-full h-12 py-2 px-4 bg-white border-2 border-primary-600 text-primary-600 text-base font-medium rounded-full shadow-md hover:bg-primary-50 transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            >
              Войти
            </button>
          </div>
        </form>

        {/* Ссылка "Забыли пароль" с отступом */}
        <div className="pt-3 text-center"> {/* Добавлен отступ сверху */}
          <Link to="/forgot-password" className="text-sm text-primary-600 hover:text-primary-500">
            Забыли пароль?
          </Link>
        </div>

        {/* Ссылка на регистрацию с увеличенными отступами */}
        <div className="text-center text-sm text-gray-600 pt-6 border-t border-gray-200 mt-4"> {/* Увеличен pt и добавлен border */}
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