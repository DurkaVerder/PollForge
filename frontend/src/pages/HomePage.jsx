import { useState, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import PollFeed from '../components/PollFeed';

const API_BASE_URL = 'http://localhost:80/api';

export default function HomePage() {
  const [polls, setPolls] = useState([]);
  const [cursor, setCursor] = useState(null);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [error, setError] = useState(null);

  const fetchPolls = async (cursor = null) => {
    if (cursor) setIsLoadingMore(true);
    else setLoading(true);

    try {
      const url = cursor
        ? `${API_BASE_URL}/streamline/news?limit=10&cursor=${encodeURIComponent(cursor)}`
        : `${API_BASE_URL}/streamline/news?limit=10`;
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setPolls((prevPolls) =>
        cursor ? [...prevPolls, ...data.polls.polls] : data.polls.polls
      );
      setHasMore(data.hasMore);
      setCursor(data.next_cursor);
    } catch (err) {
      console.error('Error fetching polls:', err);
      setError(err.message);
    } finally {
      if (cursor) setIsLoadingMore(false);
      else setLoading(false);
    }
  };

  useEffect(() => {
    fetchPolls();
  }, []);

  return (
    <main className="">
      <div className="max-w-screen-xl mx-auto flex flex-col lg:flex-row gap-6">
        <Sidebar />
        <div className="flex-1">
          {/* Мобильные вкладки */}
          <div className="lg:hidden flex bg-white rounded-lg shadow-md mb-6 overflow-hidden">
            <NavLink
              to="/"
              className={({ isActive }) =>
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'} font-medium`
              }
            >
              Лента
            </NavLink>
            <NavLink
              to="/profile"
              className={({ isActive }) =>
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              Профиль
            </NavLink>
            <NavLink
              to="/my-polls"
              className={({ isActive }) =>
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              Мои опросы
            </NavLink>
          </div>

          {loading && polls.length === 0 ? (
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
                    Произошла ошибка при загрузке опросов. Попробуйте позже.
                  </p>
                </div>
              </div>
            </div>
          ) : (
            <PollFeed
              polls={polls}
              hasMore={hasMore}
              onLoadMore={() => fetchPolls(cursor)}
              isLoadingMore={isLoadingMore}
            />
          )}
        </div>
      </div>
    </main>
  );
}