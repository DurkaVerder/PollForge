import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

export default function PollDetailPage() {
  const { link } = useParams();
  const [poll, setPoll] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchPoll = async () => {
      try {
        const response = await fetch(`http://localhost:80/api/forms/link/${link}`, {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          },
        });

        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('Опрос не найден');
          }
          throw new Error('Ошибка загрузки опроса');
        }

        const data = await response.json();
        setPoll(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPoll();
  }, [link]);

  const formatDate = (dateString) => {
    const options = {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    };
    return new Date(dateString).toLocaleDateString('ru-RU', options);
  };

  const isActive = (expiresAt) => {
    if (!expiresAt) return true;
    return new Date(expiresAt) > new Date();
  };

  if (loading) {
    return (
      <main className="flex flex-col lg:flex-row gap-6">
        <Sidebar />
        <div className="flex-1 flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500"></div>
        </div>
      </main>
    );
  }

  if (error) {
    return (
      <main className="flex flex-col lg:flex-row gap-6">
        <Sidebar />
        <div className="flex-1">
          <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-6">
            <div className="flex">
              <div className="flex-shrink-0">
                <span className="material-symbols-outlined text-red-500">error</span>
              </div>
              <div className="ml-3">
                <p className="text-sm text-red-700">{error}</p>
              </div>
            </div>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      <div className="flex-1">
        <div className="mb-4">
          <Link
            to="/my-polls"
            className="text-primary-600 hover:text-primary-700 flex items-center"
          >
            <span className="material-symbols-outlined mr-1">arrow_back</span>
            Назад к моим опросам
          </Link>
        </div>
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center mb-4">
            <span className="material-symbols-outlined text-4xl text-primary-500 mr-2">poll</span>
            <div>
              <h1 className="text-3xl font-bold">{poll.title}</h1>
              <p className="text-gray-500">{poll.description}</p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 mb-6">
            <div className="flex items-center">
              <span className="material-symbols-outlined text-gray-500 mr-2">event</span>
              <div>
                <span className="text-gray-500">Завершается:</span>
                <span className="ml-2">
                  {poll.expires_at ? formatDate(poll.expires_at) : 'Без ограничений'}
                </span>
              </div>
            </div>
            <div className="flex items-center">
              <span
                className={`material-symbols-outlined mr-2 ${
                  isActive(poll.expires_at) ? 'text-green-500' : 'text-gray-500'
                }`}
              >
                {isActive(poll.expires_at) ? 'check_circle' : 'cancel'}
              </span>
              <div>
                <span className="text-gray-500">Статус:</span>
                <span
                  className={`ml-2 ${
                    isActive(poll.expires_at) ? 'text-green-600' : 'text-gray-600'
                  }`}
                >
                  {isActive(poll.expires_at) ? 'Активен' : 'Завершен'}
                </span>
              </div>
            </div>
          </div>

          <div className="space-y-6">
            {poll.questions.map((question) => (
              <div key={question.id} className="bg-gray-50 rounded-lg p-4 shadow-sm">
                <h3 className="text-xl font-semibold mb-2">{question.title}</h3>
                <div className="space-y-3">
                  {question.answers.map((answer) => (
                    <div key={answer.id} className="flex items-center justify-between">
                      <span>{answer.title}</span>
                      <div className="w-1/2 bg-gray-200 rounded-full h-2.5">
                        <div
                          className="bg-primary-500 h-2.5 rounded-full"
                          style={{ width: `${(answer.count / 100) * 100}%` }}
                        ></div>
                      </div>
                      <span className="text-gray-600 ml-4">{answer.count} голосов</span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>

          <div className="mt-6 flex justify-end space-x-3">
            <button className="flex items-center px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 transition-colors">
              <span className="material-symbols-outlined mr-1">bar_chart</span>
              Статистика
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}