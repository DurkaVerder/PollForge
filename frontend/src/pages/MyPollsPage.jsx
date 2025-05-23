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
      toast.success('Опрос успешно удален', {
        position: "top-right",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        theme: "light",
      });
    } catch (err) {
      console.error('Delete poll error:', err);
      toast.error(err.message, {
        position: "top-right",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        theme: "light",
      });
    } finally {
      closeDeleteModal();
    }
  };

  const sortPolls = () => {
    if (!polls) return;
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

  if (!polls || polls.length === 0) {
    return (
      <div className="flex flex-col lg:flex-row gap-6 p-6">
        <Sidebar />
        <div className="bg-white rounded-2xl shadow-lg p-10 text-center w-full transform transition-all duration-300 hover:shadow-xl">
          <div className="mx-auto w-24 h-24 bg-gradient-to-br from-blue-100 to-blue-200 rounded-full flex items-center justify-center mb-6 animate-pulse">
            <span className="material-symbols-outlined text-blue-500 text-5xl">poll</span>
          </div>
          <h3 className="text-2xl font-semibold text-gray-800 mb-3">У вас пока нет опросов</h3>
          <p className="text-gray-600 mb-8 text-lg">Создайте свой первый опрос и начните собирать мнения</p>
          <Link 
            to="/create-poll" 
            className="inline-flex items-center px-6 py-3  bg-primary-500  text-white rounded-xl  hover:bg-primary-600  shadow-md"
          >
            <span className="material-symbols-outlined mr-2">add</span> Создать опрос
          </Link>
        </div>
      </div>      
    );
  }

  return (
    <main className="flex flex-col lg:flex-row gap-6 p-6 bg-gray-50 min-h-screen">
      <Sidebar />
      
      <div className="flex-1">
        <div className="flex justify-between items-center mb-8">
          <h2 className="text-3xl font-bold text-gray-800">Мои опросы</h2>
          <div className="flex items-center space-x-4">
            <select
              value={sortOption}
              onChange={(e) => setSortOption(e.target.value)}
              className="border border-gray-200 rounded-xl px-4 py-2 bg-white shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200"
            >
              <option value="date_desc">Новые</option>
              <option value="date_asc">Старые</option>
              <option value="status_active">Активные</option>
              <option value="status_expired">Завершенные</option>
            </select>
            <Link 
              to="/create-poll" 
              className=" bg-primary-500 text-white px-3 py-2 rounded-xl hover:bg-primary-600 transition-all duration-300 shadow-md flex items-center"
            >
              <span className="material-symbols-outlined mr-3">add</span> Создать новый опрос
            </Link>
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-blue-500"></div>
          </div>
        ) : error ? (
          <div className="bg-red-50 border-l-4 border-red-500 p-6 mb-8 rounded-xl shadow-sm">
            <div className="flex items-center">
              <span className="material-symbols-outlined text-red-500 text-2xl mr-3">error</span>
              <p className="text-base text-red-700">Ошибка загрузки опросов: {error}</p>
            </div>
          </div>
        ) : sortedPolls.length === 0 ? (
          <div className="bg-white rounded-2xl shadow-lg p-10 text-center transform transition-all duration-300 hover:shadow-xl">
            <div className="mx-auto w-24 h-24 bg-gradient-to-br from-blue-100 to-blue-200 rounded-full flex items-center justify-center mb-6 animate-pulse">
              <span className="material-symbols-outlined text-blue-500 text-5xl">poll</span>
            </div>
            <h3 className="text-2xl font-semibold text-gray-800 mb-3">Нет опросов по выбранному фильтру</h3>
            <p className="text-gray-600 mb-8 text-lg">Попробуйте изменить параметры сортировки</p>
          </div>
        ) : (
          <div className="space-y-6">
            {sortedPolls.map((poll, index) => {
              const status = getPollStatus(poll.expires_at);
              return (
                <div
                  key={poll.id}
                  className="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 animate-fade-in"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <h3 className="text-xl font-semibold text-gray-800">{poll.title}</h3>
                      <p className="text-gray-600 mt-2 text-base line-clamp-2">{poll.description}</p>
                      {poll.theme && (
                        <div className="mt-3">
                          <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                            {poll.theme}
                          </span>
                        </div>
                      )}
                    </div>
                    <div className="flex items-center space-x-3">
                      <div 
                        className={`w-4 h-4 rounded-full ${
                          status === 'active' ? 'bg-green-400' : 'bg-red-400'
                        } animate-pulse`}
                        title={status === 'active' ? 'Активен' : 'Завершен'}
                      ></div>
                      <button 
                        onClick={() => openDeleteModal(poll.id)}
                        className="p-2 text-gray-400 hover:text-red-500 transition-colors duration-200"
                        title="Удалить опрос"
                      >
                        <span className="material-symbols-outlined text-2xl">delete</span>
                      </button>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-6 mb-4">
                    <div className="text-sm">
                      <div className="text-gray-500 font-medium">Дата создания</div>
                      <div className="text-gray-800">{formatDate(poll.created_at)}</div>
                    </div>
                    <div className="text-sm">
                      <div className="text-gray-500 font-medium">Завершается</div>
                      <div className="text-gray-800">{poll.expires_at ? formatDate(poll.expires_at) : 'Без ограничений'}</div>
                    </div>
                  </div>

                  <div className="flex justify-between items-center pt-4 border-t border-gray-100">
                    <Link 
                      to={`/poll/${poll.link}`} 
                      className="text-blue-600 hover:text-blue-700 font-medium flex items-center transition-colors duration-200"
                    >
                      Открыть опрос
                      <span className="material-symbols-outlined ml-2">chevron_right</span>
                    </Link>
                    <div className="flex space-x-4">
                      <button 
                        onClick={() => openShareModal(poll.link)}
                        className="text-gray-500 hover:text-blue-600 transition-colors duration-200"
                        title="Поделиться"
                      >
                        <span className="material-symbols-outlined text-2xl">share</span>
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <SharePollModal 
        pollLink={currentPollLink} 
        isOpen={isShareModalOpen} 
        onClose={closeShareModal} 
      />

      {isDeleteModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center p-4 z-50 transition-all duration-300">
          <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8 transform transition-all duration-300 scale-95 animate-modal-in">
            <div className="flex justify-between items-start mb-6">
              <h3 className="text-2xl font-semibold text-gray-800">Удаление опроса</h3>
              <button 
                onClick={closeDeleteModal}
                className="text-gray-400 hover:text-gray-600 transition-colors duration-200"
              >
                <span className="material-symbols-outlined text-2xl">close</span>
              </button>
            </div>
            <p className="text-gray-600 mb-8 text-base">Вы уверены, что хотите удалить этот опрос? Это действие нельзя отменить.</p>
            <div className="flex justify-end space-x-4">
              <button
                onClick={closeDeleteModal}
                className="px-5 py-2 border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors duration-200 text-gray-700 font-medium"
              >
                Отмена
              </button>
              <button
                onClick={handleDeletePoll}
                className="px-5 py-2 bg-gradient-to-r from-red-500 to-red-600 text-white rounded-xl hover:from-red-600 hover:to-red-700 transition-all duration-300 shadow-md"
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