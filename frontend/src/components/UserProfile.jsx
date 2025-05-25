import { useState, useEffect } from 'react';
import defaultAvatar from '../static/img/default-avatar.png';
import { Link } from 'react-router-dom';

export default function UserProfile() {
  const [profile, setProfile] = useState({
    id: 1,
    name: 'Загрузка...',
    email: '',
    bio: '',
    avatar_url: '',
  });
  const [stats, setStats] = useState({
    count_created_polls: 0,
    count_answered_polls: 0,
    count_comments_by_poll: 0,
    themes: [], // Initialize as empty array to avoid null
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const userId = localStorage.getItem('userId');
        const authToken = localStorage.getItem('authToken');

        if (!userId || !authToken) {
          throw new Error('Пользователь не авторизован');
        }

        const [profileResponse, statsResponse] = await Promise.all([
          fetch('http://localhost:80/api/profile/', {
            headers: {
              Authorization: `Bearer ${authToken}`,
            },
          }),
          fetch(`http://localhost:80/api/stats/profile/${userId}`, {
            headers: {
              Authorization: `Bearer ${authToken}`,
            },
          }),
        ]);

        if (!profileResponse.ok) {
          throw new Error('Ошибка загрузки профиля');
        }
        if (!statsResponse.ok) {
          throw new Error('Ошибка загрузки статистики');
        }

        const profileData = await profileResponse.json();
        const statsData = await statsResponse.json();

        // Ensure themes is an array
        setProfile(profileData);
        setStats({
          ...statsData,
          themes: Array.isArray(statsData.themes) ? statsData.themes : [],
        });
      } catch (err) {
        setError(err.message);
        console.error('Fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  const EditProfileURL = `/profile/${profile.id}/edit`;

  if (loading) {
    return (
      <section className="bg-white rounded-lg shadow-md p-6 mb-8 flex justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
      </section>
    );
  }

  if (error) {
    return (
      <section className="bg-white rounded-lg shadow-md p-6 mb-8">
        <div className="text-red-500">{error}</div>
      </section>
    );
  }

  return (
    <section className="bg-white rounded-lg shadow-md p-6 mb-8" id="profile">
      <div className="flex flex-col md:flex-row md:items-center mb-6 gap-6">
        <div className="flex-shrink-0">
          <div className="h-24 w-24 rounded-full overflow-hidden border-4 border-primary-100 shadow-md">
            <img
              src={profile.avatar_url || defaultAvatar}
              alt="Профиль пользователя"
              className="h-full w-full object-cover"
            />
          </div>
        </div>
        <div className="flex-1">
          <h2 className="text-2xl font-bold">{profile.name}</h2>
          <p className="text-gray-500">{profile.name}</p>
          <p className="mt-2">
            {profile.bio || 'Пользователь пока не добавил информацию о себе'}
          </p>
        </div>
        <div className="flex space-x-3 mt-4 md:mt-0">
          <Link
            to={EditProfileURL}
            className="bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded-lg transition-colors duration-300"
          >
            Редактировать профиль
          </Link>
          
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">{stats.count_created_polls}</div>
          <div className="text-gray-500">Создано опросов</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">{stats.count_answered_polls}</div>
          <div className="text-gray-500">Всего голосов</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">{stats.count_comments_by_poll}</div>
          <div className="text-gray-500">Комментариев к опросам</div>
        </div>
      </div>

      <div className="text-xl font-bold mb-4">Темы опросов</div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {stats.themes.length > 0 ? (
          stats.themes.map((theme, index) => (
            <button
              key={index}
              className="border rounded-lg p-4 hover:shadow-md transition-shadow duration-300 hover:bg-gray-50 cursor-pointer text-left flex flex-col"
            >
              <div className="flex justify-between mb-2">
                <h4 className="font-semibold">{theme.name}</h4>
                <span className="text-sm text-gray-500">{theme.count_polls} опросов</span>
              </div>
              <p className="text-gray-500 text-sm mb-3">{theme.description || 'Без описания'}</p>
              <div className="mt-auto flex justify-between items-center">
                <span className="text-sm text-primary-600">{theme.count_votes} голосов всего</span>
              </div>
            </button>
          ))
        ) : (
          <p className="text-gray-500 col-span-full">Темы опросов отсутствуют</p>
        )}
      </div>

      {stats.themes.length > 0 && (
        <div className="mt-8">
          <h3 className="text-xl font-bold mb-4">Аналитика по темам</h3>
          <div className="bg-white rounded-lg shadow p-4 mb-6">
            <div className="h-[300px] w-full bg-gray-50 rounded-lg flex items-center justify-center">
              <svg className="w-full h-full" viewBox="0 0 600 300">
                <g transform="translate(50,30)">
                  {stats.themes.map((theme, index) => {
                    const maxPolls = Math.max(...stats.themes.map((t) => t.count_polls)) || 1; // Avoid division by zero
                    const barHeight = (theme.count_polls / maxPolls) * 200;
                    const colors = ['#818cf8', '#a78bfa', '#c084fc', '#e879f9', '#f472b6'];
                    return (
                      <g key={index}>
                        <rect
                          x={index * 100}
                          y={230 - barHeight}
                          width="80"
                          height={barHeight}
                          fill={colors[index % colors.length]}
                          rx="4"
                        />
                        <text
                          x={index * 100 + 40}
                          y="250"
                          textAnchor="middle"
                          fill="#4b5563"
                          fontSize="14"
                        >
                          {theme.name.split(' ')[0]}
                        </text>
                        <text
                          x={index * 100 + 40}
                          y={220 - barHeight}
                          textAnchor="middle"
                          fill="#4b5563"
                        >
                          {theme.count_polls}
                        </text>
                      </g>
                    );
                  })}
                  <text
                    x="-40"
                    y="130"
                    textAnchor="middle"
                    transform="rotate(-90 -40 130)"
                    fill="#4b5563"
                    fontSize="14"
                  >
                    Количество опросов
                  </text>
                  <text x="240" y="280" textAnchor="middle" fill="#4b5563" fontSize="14">
                    Темы
                  </text>
                  <text
                    x="240"
                    y="-10"
                    textAnchor="middle"
                    fill="#4b5563"
                    fontSize="16"
                    fontWeight="bold"
                  >
                    Распределение опросов по темам
                  </text>
                </g>
              </svg>
            </div>
          </div>
        </div>
      )}
    </section>
  );
}