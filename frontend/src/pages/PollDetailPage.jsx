import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

export default function PollDetailPage() {
  const { link } = useParams();
  const [poll, setPoll] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('questions');

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

  const calculateTotalVotes = () => {
    if (!poll) return 0;
    return poll.questions.reduce((total, question) => {
      return total + question.answers.reduce((sum, answer) => sum + answer.count, 0);
    }, 0);
  };

  const getMostPopularAnswer = (question) => {
    if (!question.answers.length) return null;
    return question.answers.reduce((prev, current) => 
      (prev.count > current.count) ? prev : current
    );
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
          <div className="flex items-start mb-4">
            <span className="material-symbols-outlined text-4xl text-primary-500 mr-3">poll</span>
            <div className="flex-1">
              <h1 className="text-3xl font-bold mb-2">{poll.title}</h1>
              <p className="text-gray-500 mb-4">{poll.description}</p>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex items-center mb-2">
                    <span className="material-symbols-outlined text-gray-500 mr-2">event</span>
                    <span className="text-gray-500">Дата создания:</span>
                  </div>
                  <p className="font-medium">{formatDate(poll.created_at)}</p>
                </div>
                
                <div className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex items-center mb-2">
                    <span className="material-symbols-outlined text-gray-500 mr-2">schedule</span>
                    <span className="text-gray-500">Статус:</span>
                  </div>
                  <p className={`font-medium ${
                    isActive(poll.expires_at) ? 'text-green-600' : 'text-gray-600'
                  }`}>
                    {isActive(poll.expires_at) ? (
                      <span className="flex items-center">
                        <span className="material-symbols-outlined text-green-500 mr-1">check_circle</span>
                        Активен
                      </span>
                    ) : (
                      <span className="flex items-center">
                        <span className="material-symbols-outlined text-gray-500 mr-1">cancel</span>
                        Завершен
                      </span>
                    )}
                  </p>
                </div>
                
                <div className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex items-center mb-2">
                    <span className="material-symbols-outlined text-gray-500 mr-2">groups</span>
                    <span className="text-gray-500">Всего ответов:</span>
                  </div>
                  <p className="font-medium">{calculateTotalVotes()}</p>
                </div>
              </div>
              
              {poll.expires_at && (
                <div className="bg-blue-50 border-l-4 border-blue-500 p-4 mb-6">
                  <div className="flex items-center">
                    <span className="material-symbols-outlined text-blue-500 mr-2">info</span>
                    <div>
                      <p className="text-blue-700">
                        Опрос {isActive(poll.expires_at) ? 'завершится' : 'завершился'} {formatDate(poll.expires_at)}
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="border-b border-gray-200 mb-6">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => setActiveTab('questions')}
                className={`whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'questions'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                Вопросы и ответы
              </button>
              <button
                onClick={() => setActiveTab('stats')}
                className={`whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'stats'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                Статистика
              </button>
              <button
                onClick={() => setActiveTab('settings')}
                className={`whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'settings'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                Настройки
              </button>
            </nav>
          </div>

          {activeTab === 'questions' && (
            <div className="space-y-6">
              {poll.questions.map((question) => {
                const mostPopular = getMostPopularAnswer(question);
                const totalVotes = question.answers.reduce((sum, answer) => sum + answer.count, 0);
                
                return (
                  <div key={question.id} className="bg-gray-50 rounded-lg p-6 shadow-sm">
                    <div className="flex justify-between items-start mb-4">
                      <h3 className="text-xl font-semibold">{question.title}</h3>
                      {mostPopular && (
                        <div className="flex items-center bg-green-50 text-green-800 px-3 py-1 rounded-full text-sm">
                          <span className="material-symbols-outlined text-green-500 mr-1">star</span>
                          Популярный ответ: {mostPopular.title}
                        </div>
                      )}
                    </div>
                    
                    <div className="space-y-4">
                      {question.answers.map((answer) => (
                        <div key={answer.id} className="flex items-center justify-between">
                          <div className="flex items-center">
                            <span className="mr-3">{answer.title}</span>
                            {answer.id === mostPopular?.id && (
                              <span className="material-symbols-outlined text-yellow-500">star</span>
                            )}
                          </div>
                          <div className="flex items-center">
                            <span className="text-gray-500 mr-2">{answer.count}</span>
                            <span className="text-gray-400 text-sm">({totalVotes > 0 ? Math.round((answer.count / totalVotes) * 100) : 0}%)</span>
                          </div>
                        </div>
                      ))}
                    </div>
                    
                    <div className="mt-4 pt-4 border-t border-gray-200 text-sm text-gray-500">
                      Всего ответов на вопрос: {totalVotes}
                    </div>
                  </div>
                );
              })}
            </div>
          )}

          {activeTab === 'stats' && (
            <div className="bg-gray-50 rounded-lg p-6">
              <h3 className="text-xl font-semibold mb-4">Общая статистика</h3>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                <div className="bg-white p-4 rounded-lg shadow-sm">
                  <h4 className="font-medium mb-3 flex items-center">
                    <span className="material-symbols-outlined text-primary-500 mr-2">query_stats</span>
                    Основные метрики
                  </h4>
                  <ul className="space-y-3">
                    <li className="flex justify-between">
                      <span className="text-gray-500">Всего вопросов:</span>
                      <span className="font-medium">{poll.questions.length}</span>
                    </li>
                    <li className="flex justify-between">
                      <span className="text-gray-500">Всего ответов:</span>
                      <span className="font-medium">{calculateTotalVotes()}</span>
                    </li>
                    <li className="flex justify-between">
                      <span className="text-gray-500">Среднее ответов на вопрос:</span>
                      <span className="font-medium">
                        {poll.questions.length > 0 
                          ? Math.round(calculateTotalVotes() / poll.questions.length) 
                          : 0}
                      </span>
                    </li>
                  </ul>
                </div>
                
                <div className="bg-white p-4 rounded-lg shadow-sm">
                  <h4 className="font-medium mb-3 flex items-center">
                    <span className="material-symbols-outlined text-primary-500 mr-2">trending_up</span>
                    Активность
                  </h4>
                  <ul className="space-y-3">
                    <li className="flex justify-between">
                      <span className="text-gray-500">Дата создания:</span>
                      <span className="font-medium">{formatDate(poll.created_at)}</span>
                    </li>
                    <li className="flex justify-between">
                      <span className="text-gray-500">Статус:</span>
                      <span className={`font-medium ${
                        isActive(poll.expires_at) ? 'text-green-600' : 'text-gray-600'
                      }`}>
                        {isActive(poll.expires_at) ? 'Активен' : 'Завершен'}
                      </span>
                    </li>
                    {poll.expires_at && (
                      <li className="flex justify-between">
                        <span className="text-gray-500">Дата завершения:</span>
                        <span className="font-medium">{formatDate(poll.expires_at)}</span>
                      </li>
                    )}
                  </ul>
                </div>
              </div>
              
              <h4 className="font-medium mb-3 flex items-center">
                <span className="material-symbols-outlined text-primary-500 mr-2">bar_chart</span>
                Распределение ответов
              </h4>
              
              <div className="bg-white p-4 rounded-lg shadow-sm">
                {poll.questions.map((question) => {
                  const totalVotes = question.answers.reduce((sum, answer) => sum + answer.count, 0);
                  
                  return (
                    <div key={question.id} className="mb-6 last:mb-0">
                      <h5 className="font-medium mb-3">{question.title}</h5>
                      <div className="space-y-2">
                        {question.answers.map((answer) => (
                          <div key={answer.id} className="flex items-center">
                            <div className="w-1/3 mr-2">{answer.title}</div>
                            <div className="flex-1 flex items-center">
                              <div className="w-full bg-gray-200 rounded-full h-2.5">
                                <div
                                  className="bg-primary-500 h-2.5 rounded-full"
                                  style={{ 
                                    width: `${totalVotes > 0 ? (answer.count / totalVotes) * 100 : 0}%` 
                                  }}
                                ></div>
                              </div>
                              <span className="ml-2 text-sm text-gray-500 w-16 text-right">
                                {totalVotes > 0 ? Math.round((answer.count / totalVotes) * 100) : 0}%
                              </span>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {activeTab === 'settings' && (
            <div className="bg-gray-50 rounded-lg p-6">
              <h3 className="text-xl font-semibold mb-4">Настройки опроса</h3>
              
              <div className="bg-white p-6 rounded-lg shadow-sm">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h4 className="font-medium mb-3 flex items-center">
                      <span className="material-symbols-outlined text-primary-500 mr-2">settings</span>
                      Основные настройки
                    </h4>
                    <ul className="space-y-3">
                      <li className="flex justify-between">
                        <span className="text-gray-500">Доступ:</span>
                        <span className="font-medium">Публичный</span>
                      </li>
                    </ul>
                  </div>
                  
                  <div>
                    <h4 className="font-medium mb-3 flex items-center">
                      <span className="material-symbols-outlined text-primary-500 mr-2">link</span>
                      Ссылка на опрос
                    </h4>
                    <div className="flex">
                      <input
                        type="text"
                        readOnly
                        value={`http://localhost:3000/poll/vote/${link}`}
                        className="flex-1 border border-gray-300 rounded-l-lg px-3 py-2 text-sm"
                      />
                      <button className="bg-primary-500 text-white px-4 py-2 rounded-r-lg hover:bg-primary-600 transition-colors">
                        Копировать
                      </button>
                    </div>
                  </div>
                </div>
                
                <div className="mt-6 pt-6 border-t border-gray-200 flex justify-end space-x-3">
                  <Link 
                    to={`/poll/edit/${poll.id}`}
                    className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                  >
                    Редактировать
                  </Link>
                  <button className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors">
                    Удалить опрос
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </main>
  );
}