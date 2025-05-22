import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import defaultAvatar from '../static/img/default-avatar.png';

// Базовый URL API
const API_BASE_URL = 'http://localhost:80/api';

export default function PollCard({ poll }) {
  const navigate = useNavigate();
  const [localPoll, setLocalPoll] = useState(poll);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [creator, setCreator] = useState(null);
  const [isLoadingCreator, setIsLoadingCreator] = useState(true);

  useEffect(() => {
    const fetchCreator = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/profile/user/${localPoll.creator_id}`, {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch creator');
        }

        const data = await response.json();
        setCreator(data);
      } catch (error) {
        console.error('Error fetching creator:', error);
      } finally {
        setIsLoadingCreator(false);
      }
    };

    fetchCreator();
  }, [localPoll.creator_id]);

  const handleAvatarClick = () => {
    if (localPoll.creator_id) {
      navigate(`/profile/${localPoll.creator_id}`);
    }
  };

  const handleVote = async (questionId, answerId, isSelected) => {
    setIsSubmitting(true);
    
    try {
      // Находим предыдущий выбранный ответ в этом вопросе
      const question = localPoll.questions.find(q => q.id === questionId);
      const prevSelectedAnswer = question?.answers.find(a => a.is_selected && a.id !== answerId);

      // Если был выбран другой вариант, сначала отменяем предыдущий выбор
      if (prevSelectedAnswer) {
        await fetch(`${API_BASE_URL}/vote/input`, {
          method: 'POST',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            id: prevSelectedAnswer.id,
            is_up_vote: false
          })
        });
      }

      // Отправляем новый выбор
      const response = await fetch(`${API_BASE_URL}/vote/input`, {
        method: 'POST',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          id: answerId,
          is_up_vote: !isSelected
        })
      });

      if (!response.ok) {
        throw new Error('Ошибка при отправке голоса');
      }

      // Обновляем локальное состояние
      setLocalPoll(prevPoll => {
        const updatedQuestions = prevPoll.questions.map(question => {
          if (question.id === questionId) {
            const updatedAnswers = question.answers.map(answer => {
              // Сбрасываем выбор для всех вариантов
              const newAnswer = {
                ...answer,
                is_selected: false,
                count_votes: answer.is_selected ? answer.count_votes - 1 : answer.count_votes
              };

              // Устанавливаем выбор для текущего варианта
              if (answer.id === answerId) {
                return {
                  ...newAnswer,
                  is_selected: !isSelected,
                  count_votes: isSelected ? newAnswer.count_votes : newAnswer.count_votes + 1
                };
              }

              return newAnswer;
            });

            // Пересчитываем общее количество голосов
            const newTotalVotes = updatedAnswers.reduce((sum, a) => sum + a.count_votes, 0);
            
            // Пересчитываем проценты
            const answersWithPercent = updatedAnswers.map(answer => ({
              ...answer,
              percent: newTotalVotes > 0 ? Math.round((answer.count_votes / newTotalVotes) * 100) : 0
            }));

            return {
              ...question,
              answers: answersWithPercent,
              total_count_votes: newTotalVotes
            };
          }
          return question;
        });

        return {
          ...prevPoll,
          questions: updatedQuestions
        };
      });
    } catch (error) {
      console.error('Ошибка при голосовании:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 transform hover:shadow-lg transition-all duration-300">
      <div>
        <div className="flex items-center mb-3">
          <span className="mr-2 bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
            Программирование
          </span>
          <span className="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
            Технологии
          </span>
        </div>
      </div>

      <div className="flex justify-between items-start mb-4">
        <div className="flex items-center">
          {isLoadingCreator ? (
            <div className="h-10 w-10 rounded-full bg-gray-200 animate-pulse mr-3"></div>
          ) : (
            <button 
              onClick={handleAvatarClick}
              className="mr-3 focus:outline-none"
            >
              <img
                src={creator?.avatar_url || defaultAvatar}
                alt="Аватар пользователя"
                className="h-10 w-10 rounded-full object-cover hover:ring-2 hover:ring-primary-500 transition-all duration-200"
                onError={(e) => {
                  e.target.src = 'https://images.unsplash.com/photo-1633332755192-727a05c4013d?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwxfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080';
                }}
              />
            </button>
          )}
          <div>
            <button 
              onClick={handleAvatarClick}
              className="font-semibold hover:text-primary-600 focus:outline-none text-left"
            >
              {isLoadingCreator ? 'Загрузка...' : creator?.name || 'Анонимный пользователь'}
            </button>
            <p className="text-sm text-gray-500">
              {new Date(localPoll.created_at).toLocaleDateString('ru-RU')}
            </p>
          </div>
        </div>
        <button className="text-gray-400 hover:text-gray-600">
          <span className="material-symbols-outlined">more_vert</span>
        </button>
      </div>

      <h3 className="text-xl font-semibold mb-2">{localPoll.title}</h3>
      <p className="text-gray-600 mb-4">{localPoll.description}</p>

      {localPoll.questions.map((question, qIndex) => (
        <div key={qIndex} className={qIndex < localPoll.questions.length - 1 ? "border-b pb-4 mb-4" : "mb-4"}>
          <h4 className="text-lg font-semibold mb-3">{question.title}</h4>
          <div className="space-y-3 mb-6">
            {question.answers.map((answer, aIndex) => (
              <div key={aIndex} className="flex items-center">
                <input
                  type="radio"
                  id={`poll${localPoll.id}_q${qIndex}_answer${aIndex}`}
                  name={`poll${localPoll.id}_q${qIndex}`}
                  className="h-4 w-4 text-primary-600"
                  checked={answer.is_selected}
                  onChange={() => handleVote(question.id, answer.id, answer.is_selected)}
                  disabled={isSubmitting}
                />
                <label htmlFor={`poll${localPoll.id}_q${qIndex}_answer${aIndex}`} className="ml-2 block w-full">
                  <div className="flex justify-between">
                    <span>{answer.title}</span>
                    <span className="text-sm text-gray-500">
                      {answer.count_votes} ({answer.percent}%)
                    </span>
                  </div>
                  <div className="mt-1 h-2 w-full bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-primary-500 rounded-full"
                      style={{width: `${answer.percent}%`}}
                    ></div>
                  </div>
                </label>
              </div>
            ))}
          </div>
        </div>
      ))}

      <div className="flex justify-between text-sm text-gray-500">
        <span>{localPoll.questions.reduce((sum, q) => sum + q.total_count_votes, 0)} голосов</span>
        <span>Заканчивается {new Date(localPoll.expires_at).toLocaleDateString('ru-RU')}</span>
      </div>
      <div className="mt-4 flex justify-between">
        <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors">
          <span className="material-symbols-outlined mr-1">comment</span>
          {localPoll.count_comments} комментариев
        </button>
        <div className="flex items-center">
          <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors mr-4">
            <span className="material-symbols-outlined mr-1">share</span>
            Поделиться
          </button>
          <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors">
            <span className="material-symbols-outlined mr-1">
              {localPoll.likes.is_liked ? 'favorite' : 'favorite_border'}
            </span>
            {localPoll.likes.count}
          </button>
        </div>
      </div>
    </div>
  );
}