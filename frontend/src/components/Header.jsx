import { Link, useNavigate } from 'react-router-dom';
import { useCallback, useEffect, useState } from 'react';
import defaultAvatar from '../static/img/default-avatar.png';

export default function Header() {
  const navigate = useNavigate();
  const [profile, setProfile] = useState({
    name: 'Загрузка...',
    avatar_url: ''
  });
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch('http://localhost:80/api/profile/', {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch profile');
        }

        const data = await response.json();
        setProfile(data);
      } catch (error) {
        console.error('Profile fetch error:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchProfile();
  }, []);

  const handleLogout = useCallback(() => {
    localStorage.removeItem('authToken');
    navigate('/login');
  }, [navigate]);

  return (
    <header className="bg-white rounded-lg shadow-md p-4 mb-6 flex items-center justify-between sticky top-0 z-10">
      <Link
        to="/"
        className="text-2xl font-bold text-primary-600 hover:text-primary-700 transition-colors duration-300"
      >
        PollForge
      </Link>

      <div className="flex items-center space-x-4">
        <Link
          to="/create-poll"
          className="h-10 w-10 rounded-full bg-primary-500 hover:bg-primary-600 flex items-center justify-center text-white shadow-md transform hover:scale-105 transition-all duration-300"
        >
          <span className="material-symbols-outlined">add</span>
        </Link>

        {!isLoading && (
          <div className="relative">
            <details className="group">
              <summary className="list-none cursor-pointer">
                <div className="flex items-center space-x-2">
                  <div className="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center overflow-hidden border-2 border-primary-300 hover:border-primary-500 transition-all duration-300">
                    <img
                      src={profile.avatar_url || defaultAvatar}
                      alt="Профиль пользователя"
                      className="h-full w-full object-cover"
                      onError={(e) => {
                        e.target.src = 'https://images.unsplash.com/photo-1633332755192-727a05c4013d?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwxfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080';
                      }}
                    />
                  </div>
                  <span className="hidden md:block font-medium">{profile.name}</span>
                  <span className="material-symbols-outlined text-gray-500 group-open:rotate-180 transition-transform duration-300">
                    expand_more
                  </span>
                </div>
              </summary>

              <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg overflow-hidden z-20">
                <div className="py-2">
                  <Link
                    to="/profile"
                    className="block px-4 py-2 hover:bg-gray-100 transition-colors duration-200"
                  >
                    <span className="material-symbols-outlined mr-2 align-middle">person</span>
                    Профиль
                  </Link>
                  <Link
                    to="/settings"
                    className="block px-4 py-2 hover:bg-gray-100 transition-colors duration-200"
                  >
                    <span className="material-symbols-outlined mr-2 align-middle">settings</span>
                    Настройки
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="w-full text-left px-4 py-2 hover:bg-gray-100 transition-colors duration-200 text-red-500 flex items-center"
                  >
                    <span className="material-symbols-outlined mr-2 align-middle">logout</span>
                    Выйти
                  </button>
                </div>
              </div>
            </details>
          </div>
        )}
      </div>
    </header>
  );
}