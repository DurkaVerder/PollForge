import { useState, useEffect } from 'react';
import defaultAvatar from '../static/img/default-avatar.png';
import { Link } from 'react-router-dom';

export default function AnotherUserProfile( {id}) {
  const [profile, setProfile] = useState({
    id: id || '',
    name: 'Загрузка...',
    email: '',
    bio: '',
    avatar_url: ''
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch(`http://localhost:80/api/profile/user/${id}`, {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken')
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка загрузки профиля');
        }

        const data = await response.json();
        setProfile(data);
      } catch (err) {
        setError(err.message);
        console.error('Profile fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [id]);
  const EditProfileURL = `/profile/${profile.id}/edit`;
  const topics = [
    {
      title: "Технологии",
      pollsCount: 8,
      description: "Языки программирования, фреймворки и технические предпочтения",
      votes: "2.4k"
    },
    {
      title: "Рабочая среда",
      pollsCount: 5,
      description: "Удалённая работа, офисные предпочтения и продуктивность",
      votes: "1.2k"
    },
    {
      title: "UX/UI Дизайн",
      pollsCount: 4,
      description: "Тренды дизайна, инструменты и паттерны пользовательского опыта",
      votes: "950"
    },
    {
      title: "Разработка ПО",
      pollsCount: 7,
      description: "Практики разработки, методологии и команды",
      votes: "1.8k"
    }
  ];

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
          
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">24</div>
          <div className="text-gray-500">Создано опросов</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">4.8k</div>
          <div className="text-gray-500">Всего голосов</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">256</div>
          <div className="text-gray-500">Комментариев к опросам</div>
        </div>
      </div>

      <div className="text-xl font-bold mb-4">Темы опросов</div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {topics.map((topic, index) => (
          <button
            key={index}
            className="border rounded-lg p-4 hover:shadow-md transition-shadow duration-300 hover:bg-gray-50 cursor-pointer text-left flex flex-col"
          >
            <div className="flex justify-between mb-2">
              <h4 className="font-semibold">{topic.title}</h4>
              <span className="text-sm text-gray-500">{topic.pollsCount} опросов</span>
            </div>
            <p className="text-gray-500 text-sm mb-3">{topic.description}</p>
            <div className="mt-auto flex justify-between items-center">
              <span className="text-sm text-primary-600">{topic.votes} голосов всего</span>
              <span className="material-symbols-outlined text-primary-600">
                chevron_right
              </span>
            </div>
          </button>
        ))}
      </div>

      <div className="mt-8">
        <h3 className="text-xl font-bold mb-4">Аналитика по темам</h3>
        <div className="bg-white rounded-lg shadow p-4 mb-6">
          <div className="h-[300px] w-full bg-gray-50 rounded-lg flex items-center justify-center">
            <svg className="w-full h-full" viewBox="0 0 600 300">
              <g transform="translate(50,30)">
                <rect x="0" y="30" width="80" height="200" fill="#818cf8" rx="4"></rect>
                <rect x="100" y="80" width="80" height="150" fill="#a78bfa" rx="4"></rect>
                <rect x="200" y="130" width="80" height="100" fill="#c084fc" rx="4"></rect>
                <rect x="300" y="180" width="80" height="50" fill="#e879f9" rx="4"></rect>
                <rect x="400" y="110" width="80" height="120" fill="#f472b6" rx="4"></rect>

                <text x="40" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Технологии
                </text>
                <text x="140" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Работа
                </text>
                <text x="240" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  UX/UI
                </text>
                <text x="340" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Разработка
                </text>
                <text x="440" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Другое
                </text>

                <text x="40" y="20" textAnchor="middle" fill="#4b5563">8</text>
                <text x="140" y="70" textAnchor="middle" fill="#4b5563">5</text>
                <text x="240" y="120" textAnchor="middle" fill="#4b5563">4</text>
                <text x="340" y="170" textAnchor="middle" fill="#4b5563">2</text>
                <text x="440" y="100" textAnchor="middle" fill="#4b5563">5</text>

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
          <div className="flex justify-center mt-4">
            <button className="text-primary-600 hover:text-primary-700 font-medium hover:underline transition-all duration-300 flex items-center">
              Посмотреть подробную аналитику
              <span className="material-symbols-outlined ml-1">arrow_forward</span>
            </button>
          </div>
          
        </div>
        
      </div>
    </section>
  );
}