import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

export default function EditPollPage() {
  const MAX_QUESTIONS = 10;
  const MAX_ANSWERS = 12;
  const MAX_TITLE_LENGTH = 100;
  const MAX_DESCRIPTION_LENGTH = 300;
  const MAX_QUESTION_LENGTH = 200;
  const MAX_ANSWER_LENGTH = 100;

  const { id } = useParams();
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    private_key: false,
    expires_at: ''
  });
  const [questions, setQuestions] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchFormData = async () => {
      try {
        const response = await fetch(`http://localhost:80/api/forms/${id}`, {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch form data');
        }

        const data = await response.json();

        // Format the expires_at date for the datetime-local input
        const expiresAt = data.expires_at 
          ? new Date(data.expires_at).toISOString().slice(0, 16)
          : '';

        setFormData({
          title: data.title,
          description: data.description,
          private_key: data.private_key,
          expires_at: expiresAt
        });

        // Process questions and answers from the single response
        const processedQuestions = data.questions.map(question => ({
          id: question.id,
          title: question.title,
          answers: question.answers.map(answer => ({
            id: answer.id,
            title: answer.title
          }))
        }));

        setQuestions(processedQuestions);
        setIsLoading(false);
      } catch (err) {
        setError(err.message);
        console.error('Fetch form error:', err);
        setIsLoading(false);
      }
    };

    fetchFormData();
  }, [id]);

  const handleFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleQuestionChange = (index, value) => {
    if (value.length <= MAX_QUESTION_LENGTH) {
      const newQuestions = [...questions];
      newQuestions[index].title = value;
      setQuestions(newQuestions);
    }
  };

  const handleAnswerChange = (qIndex, aIndex, value) => {
    if (value.length <= MAX_ANSWER_LENGTH) {
      const newQuestions = [...questions];
      newQuestions[qIndex].answers[aIndex].title = value;
      setQuestions(newQuestions);
    }
  };

  const addQuestion = () => {
    if (questions.length < MAX_QUESTIONS) {
      setQuestions([...questions, { title: '', answers: [{ title: '' }] }]);
    } else {
      setError(`Максимальное количество вопросов: ${MAX_QUESTIONS}`);
    }
  };

  const removeQuestion = async (index) => {
    if (questions[index].id) {
      // Existing question - delete from server
      try {
        const response = await fetch(
          `http://localhost:80/api/forms/${id}/questions/${questions[index].id}`,
          {
            method: 'DELETE',
            headers: {
              'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            },
          }
        );

        if (!response.ok) {
          throw new Error('Failed to delete question');
        }

        const newQuestions = [...questions];
        newQuestions.splice(index, 1);
        setQuestions(newQuestions);
      } catch (err) {
        setError(err.message);
        console.error('Delete question error:', err);
      }
    } else {
      // New question - just remove from state
      if (questions.length > 1) {
        const newQuestions = [...questions];
        newQuestions.splice(index, 1);
        setQuestions(newQuestions);
      }
    }
  };

  const addAnswer = (qIndex) => {
    const newQuestions = [...questions];
    if (newQuestions[qIndex].answers.length < MAX_ANSWERS) {
      newQuestions[qIndex].answers.push({ title: '' });
      setQuestions(newQuestions);
    } else {
      setError(`Максимальное количество ответов: ${MAX_ANSWERS}`);
    }
  };

  const removeAnswer = async (qIndex, aIndex) => {
    const question = questions[qIndex];
    const answer = question.answers[aIndex];

    if (answer.id) {
      // Existing answer - delete from server
      try {
        const response = await fetch(
          `http://localhost:80/api/forms/${id}/questions/${question.id}/answers/${answer.id}`,
          {
            method: 'DELETE',
            headers: {
              'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            },
          }
        );

        if (!response.ok) {
          throw new Error('Failed to delete answer');
        }

        const newQuestions = [...questions];
        newQuestions[qIndex].answers.splice(aIndex, 1);
        setQuestions(newQuestions);
      } catch (err) {
        setError(err.message);
        console.error('Delete answer error:', err);
      }
    } else {
      // New answer - just remove from state
      if (question.answers.length > 1) {
        const newQuestions = [...questions];
        newQuestions[qIndex].answers.splice(aIndex, 1);
        setQuestions(newQuestions);
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');

    if (formData.title.length > MAX_TITLE_LENGTH) {
      setError(`Название опроса не должно превышать ${MAX_TITLE_LENGTH} символов`);
      setIsSubmitting(false);
      return;
    }

    if (formData.description.length > MAX_DESCRIPTION_LENGTH) {
      setError(`Описание не должно превышать ${MAX_DESCRIPTION_LENGTH} символов`);
      setIsSubmitting(false);
      return;
    }

    try {
      // Update form data
      const date = formData.expires_at ? new Date(formData.expires_at) : null;
      const utcISOString = date ? date.toISOString() : null;

      const formResponse = await fetch(`http://localhost:80/api/forms/${id}`, {
        method: 'PUT',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: formData.title,
          description: formData.description,
          private_key: formData.private_key,
          expires_at: utcISOString
        })
      });

      if (!formResponse.ok) {
        const errorData = await formResponse.json();
        throw new Error(errorData.message || 'Ошибка при обновлении формы');
      }

      // Process questions and answers
      for (let qIndex = 0; qIndex < questions.length; qIndex++) {
        const question = questions[qIndex];

        if (question.title.length > MAX_QUESTION_LENGTH) {
          throw new Error(`Вопрос ${qIndex + 1} превышает максимальную длину (${MAX_QUESTION_LENGTH} символов)`);
        }

        if (question.id) {
          // Update existing question
          const questionResponse = await fetch(
            `http://localhost:80/api/forms/${id}/questions/${question.id}`,
            {
              method: 'PUT',
              headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                title: question.title,
                number_order: qIndex + 1
              })
            }
          );

          if (!questionResponse.ok) {
            const errorData = await questionResponse.json();
            throw new Error(errorData.message || 'Ошибка при обновлении вопроса');
          }
        } else {
          // Create new question
          const questionResponse = await fetch(
            `http://localhost:80/api/forms/${id}/questions`,
            {
              method: 'POST',
              headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                title: question.title,
                number_order: qIndex + 1
              })
            }
          );

          if (!questionResponse.ok) {
            const errorData = await questionResponse.json();
            throw new Error(errorData.message || 'Ошибка при создании вопроса');
          }

          const questionData = await questionResponse.json();
          question.id = questionData.question_id;
        }

        // Process answers for this question
        for (let aIndex = 0; aIndex < question.answers.length; aIndex++) {
          const answer = question.answers[aIndex];

          if (answer.title.length > MAX_ANSWER_LENGTH) {
            throw new Error(`Ответ ${aIndex + 1} в вопросе ${qIndex + 1} превышает максимальную длину (${MAX_ANSWER_LENGTH} символов)`);
          }

          if (answer.id) {
            // Update existing answer
            const answerResponse = await fetch(
              `http://localhost:80/api/forms/${id}/questions/${question.id}/answers/${answer.id}`,
              {
                method: 'PUT',
                headers: {
                  'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                  title: answer.title,
                  number_order: aIndex + 1
                })
              }
            );

            if (!answerResponse.ok) {
              const errorData = await answerResponse.json();
              throw new Error(errorData.message || 'Ошибка при обновлении ответа');
            }
          } else {
            // Create new answer
            const answerResponse = await fetch(
              `http://localhost:80/api/forms/${id}/questions/${question.id}/answers`,
              {
                method: 'POST',
                headers: {
                  'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                  title: answer.title,
                  number_order: aIndex + 1
                })
              }
            );

            if (!answerResponse.ok) {
              const errorData = await answerResponse.json();
              throw new Error(errorData.message || 'Ошибка при создании ответа');
            }
          }
        }
      }

      navigate('/my-polls');
    } catch (err) {
      setError(err.message);
      console.error('Update poll error:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <main className="flex flex-col lg:flex-row gap-6">
        <Sidebar />
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
            <p className="mt-4 text-gray-600">Загрузка данных опроса...</p>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      <div className="flex-1">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center mb-6">
            <span className="material-symbols-outlined text-4xl text-primary-500 mr-2">edit</span>
            <h2 className="text-2xl font-bold">Редактировать опрос</h2>
          </div>

          <div className="mb-4 p-3 bg-blue-50 text-blue-700 rounded-md text-sm">
            <p>Ограничения:</p>
            <ul className="list-disc pl-5 mt-1 space-y-1">
              <li>Максимум {MAX_QUESTIONS} вопросов</li>
              <li>Максимум {MAX_ANSWERS} ответов на каждый вопрос</li>
              <li>Название опроса: до {MAX_TITLE_LENGTH} символов</li>
              <li>Описание: до {MAX_DESCRIPTION_LENGTH} символов</li>
              <li>Вопрос: до {MAX_QUESTION_LENGTH} символов</li>
              <li>Ответ: до {MAX_ANSWER_LENGTH} символов</li>
            </ul>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-md">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit}>
            {/* Основная информация об опросе */}
            <div className="mb-8">
              <h3 className="text-lg font-semibold mb-4">Основная информация</h3>
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between items-center mb-1">
                    <label className="block text-sm font-medium text-gray-700">
                      Название опроса
                    </label>
                    <span className="text-xs text-gray-500">
                      {formData.title.length}/{MAX_TITLE_LENGTH}
                    </span>
                  </div>
                  <input
                    type="text"
                    name="title"
                    value={formData.title}
                    onChange={handleFormChange}
                    required
                    maxLength={MAX_TITLE_LENGTH}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  />
                  <div className="mt-1 h-1 w-full bg-gray-200 rounded-full">
                    <div
                      className="h-1 bg-primary-500 rounded-full progress-bar"
                      style={{ width: `${(formData.title.length / MAX_TITLE_LENGTH) * 100}%` }}
                    ></div>
                  </div>
                </div>

                <div>
                  <div className="flex justify-between items-center mb-1">
                    <label className="block text-sm font-medium text-gray-700">
                      Описание
                    </label>
                    <span className="text-xs text-gray-500">
                      {formData.description.length}/{MAX_DESCRIPTION_LENGTH}
                    </span>
                  </div>
                  <textarea
                    name="description"
                    value={formData.description}
                    onChange={handleFormChange}
                    maxLength={MAX_DESCRIPTION_LENGTH}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    rows="3"
                  />
                  <div className="mt-1 h-1 w-full bg-gray-200 rounded-full">
                    <div
                      className="h-1 bg-primary-500 rounded-full"
                      style={{ width: `${(formData.description.length / MAX_DESCRIPTION_LENGTH) * 100}%` }}
                    ></div>
                  </div>
                </div>

                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="private_key"
                    name="private_key"
                    checked={formData.private_key}
                    onChange={handleFormChange}
                    className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                  />
                  <label htmlFor="private_key" className="ml-2 block text-sm text-gray-700">
                    Приватный опрос
                  </label>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Дата окончания (необязательно)
                  </label>
                  <input
                    type="datetime-local"
                    name="expires_at"
                    value={formData.expires_at}
                    onChange={handleFormChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  />
                </div>
              </div>
            </div>

            {/* Вопросы */}
            <div className="mb-8">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-semibold">Вопросы ({questions.length}/{MAX_QUESTIONS})</h3>
                <button
                  type="button"
                  onClick={addQuestion}
                  disabled={questions.length >= MAX_QUESTIONS}
                  className={` px-3 py-1 rounded-lg text-sm ${
                    questions.length >= MAX_QUESTIONS
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                      : 'bg-primary-50 text-primary-600 hover:bg-primary-100'
                  }`}
                >
                  + Добавить вопрос
                </button>
              </div>

              {questions.map((question, qIndex) => (
                <div key={question.id || `new-${qIndex}`} className="mb-6 p-4 border border-gray-200 rounded-lg">
                  <div className="flex justify-between items-center mb-3">
                    <h4 className="font-medium">Вопрос {qIndex + 1}</h4>
                    <button
                      type="button"
                      onClick={() => removeQuestion(qIndex)}
                      className="text-red-500 text-sm hover:text-red-700"
                    >
                      Удалить вопрос
                    </button>
                  </div>

                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-1">
                      <label className="block text-sm font-medium text-gray-700">
                        Текст вопроса
                      </label>
                      <span className="text-xs text-gray-500">
                        {question.title.length}/{MAX_QUESTION_LENGTH}
                      </span>
                    </div>
                    <input
                      type="text"
                      value={question.title}
                      onChange={(e) => handleQuestionChange(qIndex, e.target.value)}
                      required
                      maxLength={MAX_QUESTION_LENGTH}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    />
                    <div className="mt-1 h-1 w-full bg-gray-200 rounded-full">
                      <div
                        className="h-1 bg-primary-500 rounded-full"
                        style={{ width: `${(question.title.length / MAX_QUESTION_LENGTH) * 100}%` }}
                      ></div>
                    </div>
                  </div>

                  <div>
                    <div className="flex justify-between items-center mb-3">
                      <label className="block text-sm font-medium text-gray-700">
                        Варианты ответов ({question.answers.length}/{MAX_ANSWERS})
                      </label>
                      <button
                        type="button"
                        onClick={() => addAnswer(qIndex)}
                        disabled={question.answers.length >= MAX_ANSWERS}
                        className={`px-3 py-1 rounded-lg text-sm ${
                          question.answers.length >= MAX_ANSWERS
                            ? 'bg-gray-100 text-gray-600 cursor-not-allowed'
                            : 'bg-primary-50 text-primary-600 hover:bg-primary-100'
                        }`}
                      >
                        + Добавить ответ
                      </button>
                    </div>

                    {question.answers.map((answer, aIndex) => (
                      <div key={answer.id || `new-${aIndex}`} className="flex items-center mb-2">
                        <div className="flex-1">
                          <div className="flex justify-between items-center mb-1">
                            <span className="text-xs text-gray-500">
                              Ответ {aIndex + 1}
                            </span>
                            <span className="text-xs text-gray-500">
                              {answer.title.length}/{MAX_ANSWER_LENGTH}
                            </span>
                          </div>
                          <input
                            type="text"
                            value={answer.title}
                            onChange={(e) => handleAnswerChange(qIndex, aIndex, e.target.value)}
                            required
                            maxLength={MAX_ANSWER_LENGTH}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                          />
                          <div className="mt-1 h-1 w-full bg-gray-200 rounded-full">
                            <div
                              className="h-1 bg-primary-500 rounded-full"
                              style={{ width: `${(answer.title.length / MAX_ANSWER_LENGTH) * 100}%` }}
                            ></div>
                          </div>
                        </div>
                        <button
                          type="button"
                          onClick={() => removeAnswer(qIndex, aIndex)}
                          className="ml-2 text-red-500 hover:text-red-700"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>

            <div className="flex justify-end space-x-4">
              <button
                type="button"
                onClick={() => navigate('/my-polls')}
                className="h-10 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className={`h-10 px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                  isSubmitting ? 'opacity-70 cursor-not-allowed' : ''
                }`}
              >
                {isSubmitting ? 'Сохранение...' : 'Сохранить'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  );
}