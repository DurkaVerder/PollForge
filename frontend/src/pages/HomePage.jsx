import { NavLink } from 'react-router-dom';
import { useState, useEffect } from 'react';
import Sidebar from '../components/Sidebar';
import PollFeed from '../components/PollFeed';

// Базовый URL API
const API_BASE_URL = 'http://localhost:80/api';

export default function HomePage() {
  const [polls, setPolls] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPolls = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/streamline/news`, {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        setPolls(data.polls || []);
      } catch (err) {
        console.error('Error fetching polls:', err);
        setError(err.message);
        // Для демонстрации используем fallback данные при ошибке
        setPolls(getFallbackPolls());
      } finally {
        setLoading(false);
      }
    };

    fetchPolls();
  }, []);

  // Функция с fallback данными на случай ошибки запроса
  const getFallbackPolls = () => {
    return [
      {
        "id": 1,
        "title": "Любимый язык программирования",
        "description": "Выберите ваш любимый и ненавистный язык программирования",
        "link": "524d0bad-3ee6-4608-b082-cb56f18cc72c",
        "likes": {
          "count": 0,
          "is_liked": false
        },
        "count_comments": 0,
        "questions": [
          {
            "id": 1,
            "title": "Какой язык программирования вы предпочитаете",
            "total_count_votes": 0,
            "answers": [
              {
                "id": 1,
                "title": "Python",
                "percent": 0,
                "count_votes": 0,
                "is_selected": false
              },
              // ... остальные данные из вашего JSON
            ]
          }
        ],
        "created_at": "2025-05-18 12:50:46",
        "expires_at": "2025-06-30 23:59:59"
      }
    ];
  };

  return (
    <main className="">
      <div className="max-w-screen-xl mx-auto flex flex-col lg:flex-row gap-6">
        <Sidebar />
        
        <div className="flex-1">
          {/* Мобильные вкладки */}
          <div className="lg:hidden flex bg-white rounded-lg shadow-md mb-6 overflow-hidden">
            <NavLink 
              to="/" 
              className={({isActive}) => 
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'} font-medium`
              }
            >
              Лента
            </NavLink>
            <NavLink
              to="/profile"
              className={({isActive}) => 
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              Профиль
            </NavLink>
            <NavLink
              to="/my-polls"
              className={({isActive}) => 
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              Мои опросы
            </NavLink>
          </div>

          {loading ? (
            <div className="flex justify-center items-center h-64">
              <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
            </div>
          ) : error ? (
            <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-6">
              <div className="flex">
                <div className="flex-shrink-0">
                  <span className="material-symbols-outlined text-red-500">error</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm text-red-700">
                    Ошибка загрузки данных: {error}. Показаны демонстрационные данные.
                  </p>
                </div>
              </div>
            </div>
          ) : null}

          <PollFeed polls={polls} />
        </div>
      </div>
    </main>
  );
}