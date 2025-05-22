import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import SharePollModal from '../components/SharePollModal';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export default function MyPollsPage() {
  const [polls, setPolls] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isShareModalOpen, setIsShareModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [currentPollLink, setCurrentPollLink] = useState('');
  const [pollToDelete, setPollToDelete] = useState(null);
  const [sortOption, setSortOption] = useState('date_desc');
  const [sortedPolls, setSortedPolls] = useState([]);

  useEffect(() => {
    fetchPolls();
  }, []);

  useEffect(() => {
    sortPolls();
  }, [polls, sortOption]);

  const fetchPolls = async () => {
    try {
      const response = await fetch('http://localhost:80/api/profile/forms', {
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken')
        }
      });

      if (!response.ok) {
        throw new Error('Ошибка загрузки опросов');
      }

      const data = await response.json();
      setPolls(data.forms);
    } catch (err) {
      setError(err.message);
      console.error('Fetch polls error:', err);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString) => {
    const options = { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    };
    return new Date(dateString).toLocaleDateString('ru-RU', options);
  };

  const getPollStatus = (expiresAt) => {
    if (!expiresAt) return 'active';
    return new Date(expiresAt) > new Date() ? 'active' : 'expired';
  };

  const openShareModal = (pollLink) => {
    setCurrentPollLink(`http://localhost:3000/poll/vote/${pollLink}`);
    setIsShareModalOpen(true);
  };

  const closeShareModal = () => {
    setIsShareModalOpen(false);
    setCurrentPollLink('');
  };

  const openDeleteModal = (pollId) => {
    setPollToDelete(pollId);
    setIsDeleteModalOpen(true);
  };

  const closeDeleteModal = () => {
    setIsDeleteModalOpen(false);
    setPollToDelete(null);
  };

  const handleDeletePoll = async () => {
    try {
      const response = await fetch(`http://localhost:80/api/forms/${pollToDelete}`, {
        method: 'DELETE',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken')
        }
      });

      if (!response.ok) {
        throw new Error('Ошибка при удалении опроса');
      }

      setPolls(polls.filter(poll => poll.id !== pollToDelete));
      toast.success('Опрос успешно удален');
    } catch (err) {
      console.error('Delete poll error:', err);
      toast.error(err.message);
    } finally {
      closeDeleteModal();
    }
  };

  const sortPolls = () => {
    let sorted = [...polls];
    switch (sortOption) {
      case 'date_desc':
        sorted.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        break;
      case 'date_asc':
        sorted.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
        break;
      case 'status_active':
        sorted = sorted.filter(poll => getPollStatus(poll.expires_at) === 'active');
        break;
      case 'status_expired':
        sorted = sorted.filter(poll => getPollStatus(poll.expires_at) === 'expired');
        break;
      default:
        break;
    }
    setSortedPolls(sorted);
  };

  if (!polls) {
    return (
      <div className="flex flex-col lg:flex-row gap-6">
        <Sidebar />
        <div className="bg-white rounded-lg shadow-md p-8 text-center w-full">
          <div className="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
            <span className="material-symbols-outlined text-gray-400 text-4xl">poll</span>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">У вас пока нет опросов</h3>
          <p className="text-gray-500 mb-6">Создайте свой первый опрос и начните собирать мнения</p>
          <Link 
            to="/create-poll" 
            className="inline-flex items-center px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600"
          >
            Создать опрос
          </Link>
        </div>
      </div>
    );
  }

  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      
      <div className="flex-1">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Мои опросы</h2>
          <div className="flex items-center space-x-4">
            <select
              value={sortOption}
              onChange={(e) => setSortOption(e.target.value)}
              className="border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="date_desc">Новые</option>
              <option value="date_asc">Старые</option>
              <option value="status_active">Активные</option>
              <option value="status_expired">Завершенные</option>
            </select>
            <Link 
              to="/create-poll" 
              className="bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded-lg transition-colors"
            >
              Создать новый опрос
            </Link>
          </div>
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
                  Ошибка загрузки опросов: {error}
                </p>
              </div>
            </div>
          </div>
        ) : sortedPolls.length === 0 ? (
          <div className="bg-white rounded-lg shadow-md p-8 text-center">
            <div className="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
              <span className="material-symbols-outlined text-gray-400 text-4xl">poll</span>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">У вас пока нет опросов</h3>
            <p className="text-gray-500 mb-6">Создайте свой первый опрос и начните собирать мнения</p>
            <Link 
              to="/create-poll" 
              className="inline-flex items-center px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600"
            >
              Создать опрос
            </Link>
          </div>
        ) : (
          <div className="space-y-6">
            {sortedPolls.map((poll, index) => {
              const status = getPollStatus(poll.expires_at);
              return (
                <div
                  key={poll.id}
                  className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow fade-in"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <h3 className="text-xl font-semibold">{poll.title}</h3>
                      <p className="text-gray-600 mt-1">{poll.description}</p>
                    </div>
                    <div className="flex items-center space-x-2">
                      <div 
                        className={`w-4 h-4 rounded-full ${
                          status === 'active' ? 'bg-green-500' : 'bg-red-500'
                        }`}
                        title={status === 'active' ? 'Активен' : 'Завершен'}
                      ></div>
                      <button 
                        onClick={() => openDeleteModal(poll.id)}
                        className="p-2 text-gray-400 hover:text-red-500 transition-colors"
                        title="Удалить опрос"
                      >
                        <span className="material-symbols-outlined">delete</span>
                      </button>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-4 mb-4">
                    <div className="text-sm">
                      <div className="text-gray-500">Дата создания</div>
                      <div>{formatDate(poll.created_at || new Date().toISOString())}</div>
                    </div>
                    <div className="text-sm">
                      <div className="text-gray-500">Завершается</div>
                      <div>{poll.expires_at ? formatDate(poll.expires_at) : 'Без ограничений'}</div>
                    </div>
                  </div>

                  <div className="flex justify-between items-center pt-4 border-t border-gray-100">
                    <Link 
                      to={`/poll/${poll.link}`} 
                      className="text-primary-600 hover:text-primary-700 font-medium flex items-center"
                    >
                      Открыть опрос
                      <span className="material-symbols-outlined ml-1">chevron_right</span>
                    </Link>
                    <div className="flex space-x-3">
                      <button 
                        onClick={() => openShareModal(poll.link)}
                        className="text-gray-500 hover:text-gray-700 transition-colors"
                        title="Поделиться"
                      >
                        <span className="material-symbols-outlined">share</span>
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Модальное окно для шаринга */}
      <SharePollModal 
        pollLink={currentPollLink} 
        isOpen={isShareModalOpen} 
        onClose={closeShareModal} 
      />

      {/* Модальное окно подтверждения удаления */}
      {isDeleteModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50 modal-fade-in">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
            <div className="flex justify-between items-start mb-4">
              <h3 className="text-xl font-semibold">Удаление опроса</h3>
              <button 
                onClick={closeDeleteModal}
                className="text-gray-400 hover:text-gray-500"
              >
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>
            <p className="mb-6">Вы уверены, что хотите удалить этот опрос? Это действие нельзя отменить.</p>
            <div className="flex justify-end space-x-3">
              <button
                onClick={closeDeleteModal}
                className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Отмена
              </button>
              <button
                onClick={handleDeletePoll}
                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
              >
                Удалить
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}