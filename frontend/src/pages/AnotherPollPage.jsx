import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

const API_BASE_URL = 'http://localhost:80/api';

export default function AnotherPollPage() {
  const { link } = useParams();
  const navigate = useNavigate();
  const [poll, setPoll] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchPoll = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/forms/link/${link}`);
        
        if (!response.ok) {
          throw new Error('Ошибка загрузки опроса');
        }

        const data = await response.json();
        // Добавляем расчет процентов для ответов
        const pollWithPercents = {
          ...data,
          questions: data.questions.map(question => {
            const totalVotes = question.answers.reduce((sum, a) => sum + a.count, 0);
            return {
              ...question,
              answers: question.answers.map(answer => ({
                ...answer,
                percent: totalVotes > 0 ? Math.round((answer.count / totalVotes) * 100) : 0,
                is_selected: false // Добавляем поле для отслеживания выбора
              })),
              total_count_votes: totalVotes
            };
          })
        };
        setPoll(pollWithPercents);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPoll();
  }, [link]);

  const handleVote = async (questionId, answerId, isSelected) => {
    setIsSubmitting(true);
    
    try {
      // Находим предыдущий выбранный ответ в этом вопросе
      const question = poll.questions.find(q => q.id === questionId);
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
      setPoll(prevPoll => {
        const updatedQuestions = prevPoll.questions.map(question => {
          if (question.id === questionId) {
            const updatedAnswers = question.answers.map(answer => {
              // Сбрасываем выбор для всех вариантов
              const newAnswer = {
                ...answer,
                is_selected: false,
                count: answer.is_selected ? answer.count - 1 : answer.count
              };

              // Устанавливаем выбор для текущего варианта
              if (answer.id === answerId) {
                return {
                  ...newAnswer,
                  is_selected: !isSelected,
                  count: isSelected ? newAnswer.count : newAnswer.count + 1
                };
              }

              return newAnswer;
            });

            // Пересчитываем общее количество голосов
            const newTotalVotes = updatedAnswers.reduce((sum, a) => sum + a.count, 0);
            
            // Пересчитываем проценты
            const answersWithPercent = updatedAnswers.map(answer => ({
              ...answer,
              percent: newTotalVotes > 0 ? Math.round((answer.count / newTotalVotes) * 100) : 0
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

    } catch (err) {
      setError(err.message);
    } finally {
      setIsSubmitting(false);
    }
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

  if (error || !poll) {
    navigate("/404");
  }

  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      
      <div className="flex-1">
        <div className="bg-white rounded-lg shadow-md p-6 transform hover:shadow-lg transition-all duration-300">
          <div>
            <div className="flex items-center mb-3">
              <span className="mr-2 bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                Опрос
              </span>
            </div>
          </div>

          <div className="flex justify-between items-start mb-4">
            <div className="flex items-center">
              <img
                src="https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwyfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080"
                alt="Аватар пользователя"
                className="h-10 w-10 rounded-full mr-3"
              />
              <div>
                <h3 className="font-semibold">Анонимный пользователь</h3>
                <p className="text-sm text-gray-500">{new Date(poll.created_at).toLocaleDateString('ru-RU')}</p>
              </div>
            </div>
            <button className="text-gray-400 hover:text-gray-600">
              <span className="material-symbols-outlined">more_vert</span>
            </button>
          </div>

          <h3 className="text-xl font-semibold mb-2">{poll.title}</h3>
          <p className="text-gray-600 mb-4">{poll.description}</p>

          {poll.questions.map((question, qIndex) => (
            <div key={qIndex} className={qIndex < poll.questions.length - 1 ? "border-b pb-4 mb-4" : "mb-4"}>
              <h4 className="text-lg font-semibold mb-3">
                {question.title}
                {question.required && <span className="text-red-500 ml-1">*</span>}
              </h4>
              <div className="space-y-3 mb-6">
                {question.answers.map((answer, aIndex) => (
                  <div key={aIndex} className="flex items-center">
                    <input
                      type="radio"
                      id={`poll${poll.id}_q${qIndex}_answer${aIndex}`}
                      name={`poll${poll.id}_q${qIndex}`}
                      className="h-4 w-4 text-primary-600"
                      checked={answer.is_selected}
                      onChange={() => handleVote(question.id, answer.id, answer.is_selected)}
                      disabled={isSubmitting}
                    />
                    <label htmlFor={`poll${poll.id}_q${qIndex}_answer${aIndex}`} className="ml-2 block w-full">
                      <div className="flex justify-between">
                        <span>{answer.title}</span>
                        <span className="text-sm text-gray-500">
                          {answer.count} ({answer.percent}%)
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
            <span>{poll.questions.reduce((sum, q) => sum + q.total_count_votes, 0)} голосов</span>
            <span>
              {poll.expires_at ? `Заканчивается ${new Date(poll.expires_at).toLocaleDateString('ru-RU')}` : 'Без ограничений'}
            </span>
          </div>
          <div className="mt-4 flex justify-between">
            <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors">
              <span className="material-symbols-outlined mr-1">comment</span>
              0 комментариев
            </button>
            <div className="flex items-center">
              <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors mr-4">
                <span className="material-symbols-outlined mr-1">share</span>
                Поделиться
              </button>
              <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors">
                <span className="material-symbols-outlined mr-1">favorite_border</span>
                0
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}