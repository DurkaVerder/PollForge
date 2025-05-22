import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { v4 as uuidv4 } from 'uuid';

export default function CreatePollPage() {
  const MAX_QUESTIONS = 10;
  const MAX_ANSWERS = 12;
  const MAX_TITLE_LENGTH = 100;
  const MAX_DESCRIPTION_LENGTH = 250;
  const MAX_QUESTION_LENGTH = 200;
  const MAX_ANSWER_LENGTH = 100;

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    private_key: false,
    expires_at: ''
  });
  const [questions, setQuestions] = useState([{ id: uuidv4(), title: '', answers: [{ id: uuidv4(), value: '' }, { id: uuidv4(), value: '' }] }]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [tooltip, setTooltip] = useState({ text: '', x: 0, y: 0 });
  const navigate = useNavigate();

  const handleFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleQuestionChange = (id, value) => {
    if (value.length <= MAX_QUESTION_LENGTH) {
      const newQuestions = questions.map(question => {
        if (question.id === id) {
          return { ...question, title: value };
        }
        return question;
      });
      setQuestions(newQuestions);
    }
  };

  const handleAnswerChange = (qId, aId, value) => {
    if (value.length <= MAX_ANSWER_LENGTH) {
      const newQuestions = questions.map(question => {
        if (question.id === qId) {
          const newAnswers = question.answers.map(answer => {
            if (answer.id === aId) {
              return { ...answer, value };
            }
            return answer;
          });
          return { ...question, answers: newAnswers };
        }
        return question;
      });
      setQuestions(newQuestions);
    }
  };

  const addQuestion = () => {
    if (questions.length < MAX_QUESTIONS) {
      setQuestions([...questions, { id: uuidv4(), title: '', answers: [{ id: uuidv4(), value: '' }, { id: uuidv4(), value: '' }] }]);
    } else {
      setError(`Максимальное количество вопросов: ${MAX_QUESTIONS}`);
    }
  };

  const removeQuestion = (id) => {
    if (questions.length > 1) {
      const questionElement = document.getElementById(`question-${id}`);
      if (questionElement) {
        questionElement.style.opacity = '0';
        setTimeout(() => {
          const newQuestions = questions.filter(question => question.id !== id);
          setQuestions(newQuestions);
        }, 500);
      }
    }
  };

  const addAnswer = (qId) => {
    const newQuestions = questions.map(question => {
      if (question.id === qId) {
        if (question.answers.length < MAX_ANSWERS) {
          return { ...question, answers: [...question.answers, { id: uuidv4(), value: '' }] };
        } else {
          setError(`Максимальное количество ответов: ${MAX_ANSWERS}`);
          return question;
        }
      }
      return question;
    });
    setQuestions(newQuestions);
  };

  const removeAnswer = (qId, aId) => {
    const newQuestions = questions.map(question => {
      if (question.id === qId) {
        if (question.answers.length > 2) {
          const answerElement = document.getElementById(`answer-${aId}`);
          if (answerElement) {
            answerElement.style.opacity = '0';
            setTimeout(() => {
              const newAnswers = question.answers.filter(answer => answer.id !== aId);
              setQuestions(questions.map(q => q.id === qId ? { ...q, answers: newAnswers } : q));
            }, 500);
          }
          return question;
        } else {
          setError('У вопроса должно быть как минимум два варианта ответа');
          return question;
        }
      }
      return question;
    });
    setQuestions(newQuestions);
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
      const date = new Date(formData.expires_at);
      const utcISOString = date.toISOString();

      const formResponse = await fetch('http://localhost:80/api/forms/', {
        method: 'POST',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: formData.title,
          theme_id: 1,
          description: formData.description,
          private_key: formData.private_key,
          expires_at: utcISOString || null
        })
      });

      const formDataResponse = await formResponse.json();

      if (!formResponse.ok) {
        throw new Error(formDataResponse.message || 'Ошибка при создании формы');
      }

      const formId = formDataResponse.form_id;

      for (let qIndex = 0; qIndex < questions.length; qIndex++) {
        const question = questions[qIndex];

        if (question.title.length > MAX_QUESTION_LENGTH) {
          throw new Error(`Вопрос ${qIndex + 1} превышает максимальную длину (${MAX_QUESTION_LENGTH} символов)`);
        }

        const questionResponse = await fetch(`http://localhost:80/api/forms/${formId}/questions`, {
          method: 'POST',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            title: question.title,
            number_order: qIndex + 1
          })
        });

        const questionData = await questionResponse.json();

        if (!questionResponse.ok) {
          throw new Error(questionData.message || 'Ошибка при создании вопроса');
        }

        const questionId = questionData.question_id;

        for (let aIndex = 0; aIndex < question.answers.length; aIndex++) {
          const answer = question.answers[aIndex].value;

          if (answer.length > MAX_ANSWER_LENGTH) {
            throw new Error(`Ответ ${aIndex + 1} в вопросе ${qIndex + 1} превышает максимальную длину (${MAX_ANSWER_LENGTH} символов)`);
          }

          const answerResponse = await fetch(
            `http://localhost:80/api/forms/${formId}/questions/${questionId}/answers`,
            {
              method: 'POST',
              headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                title: answer,
                number_order: aIndex + 1
              })
            }
          );

          const answerData = await answerResponse.json();

          if (!answerResponse.ok) {
            throw new Error(answerData.message || 'Ошибка при создании ответа');
          }
        }
      }

      toast.success('Опрос успешно создан!', {
        position: "top-right",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined,
      });

      navigate('/my-polls');
    } catch (err) {
      setError(err.message);
      console.error('Create poll error:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className="flex flex-col lg:flex-row gap-6 min-h-screen">
      <Sidebar />
      <div className="flex-1">
        <div className="bg-white rounded-xl shadow-lg p-8 max-w-4xl mx-auto">
          <div className="flex items-center mb-8">
            <span className="material-symbols-outlined text-4xl text-blue-500 mr-3">add_circle</span>
            <h2 className="text-3xl font-bold text-gray-800">Создать новый опрос</h2>
          </div>

          <div className="mb-6 p-4 bg-blue-50 text-blue-700 rounded-lg text-sm">
            <p className="font-semibold">Ограничения:</p>
            <ul className="list-disc pl-5 mt-2 space-y-1">
              <li>Максимум {MAX_QUESTIONS} вопросов</li>
              <li>Максимум {MAX_ANSWERS} ответов на каждый вопрос</li>
              <li>Название опроса: до {MAX_TITLE_LENGTH} символов</li>
              <li>Описание: до {MAX_DESCRIPTION_LENGTH} символов</li>
              <li>Вопрос: до {MAX_QUESTION_LENGTH} символов</li>
              <li>Ответ: до {MAX_ANSWER_LENGTH} символов</li>
            </ul>
          </div>

          {error && (
            <div className="mb-6 p-4 bg-red-100 text-red-700 rounded-lg flex items-center justify-between fade-in">
              <div className="flex items-center">
                <span className="material-symbols-outlined mr-2">error</span>
                <span>{error}</span>
              </div>
              <button
                onClick={() => setError('')}
                className="text-red-500 hover:text-red-700 transition-colors"
              >
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="mb-10">
              <h3 className="text-xl font-semibold text-gray-800 mb-4">Основная информация</h3>
              <div className="space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Название опроса
                  </label>
                  <input
                    type="text"
                    name="title"
                    value={formData.title}
                    onChange={handleFormChange}
                    required
                    maxLength={MAX_TITLE_LENGTH}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    onFocus={(e) => {
                      const rect = e.target.getBoundingClientRect();
                      setTooltip({
                        text: `Максимальная длина: ${MAX_TITLE_LENGTH} символов`,
                        x: rect.left,
                        y: rect.bottom + 5
                      });
                    }}
                    onBlur={() => setTooltip({ text: '', x: 0, y: 0 })}
                  />
                  <div className="mt-2 flex justify-between items-center">
                    <span className="text-xs text-gray-500">
                      {formData.title.length}/{MAX_TITLE_LENGTH}
                    </span>
                    <div className="w-1/2 h-1 bg-gray-200 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-blue-500 transition-width duration-300"
                        style={{ width: `${(formData.title.length / MAX_TITLE_LENGTH) * 100}%` }}
                      ></div>
                    </div>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Описание
                  </label>
                  <textarea
                    name="description"
                    value={formData.description}
                    onChange={handleFormChange}
                    maxLength={MAX_DESCRIPTION_LENGTH}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    rows="3"
                    onFocus={(e) => {
                      const rect = e.target.getBoundingClientRect();
                      setTooltip({
                        text: `Максимальная длина: ${MAX_DESCRIPTION_LENGTH} символов`,
                        x: rect.left,
                        y: rect.bottom + 5
                      });
                    }}
                    onBlur={() => setTooltip({ text: '', x: 0, y: 0 })}
                  />
                  <div className="mt-2 flex justify-between items-center">
                    <span className="text-xs text-gray-500">
                      {formData.description.length}/{MAX_DESCRIPTION_LENGTH}
                    </span>
                    <div className="w-1/2 h-1 bg-gray-200 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-blue-500 transition-width duration-300"
                        style={{ width: `${(formData.description.length / MAX_DESCRIPTION_LENGTH) * 100}%` }}
                      ></div>
                    </div>
                  </div>
                </div>

                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="private_key"
                    name="private_key"
                    checked={formData.private_key}
                    onChange={handleFormChange}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded transition-transform"
                  />
                  <label htmlFor="private_key" className="ml-2 text-sm text-gray-700">
                    Приватный опрос
                  </label>
                  <span className="ml-2 text-gray-500" title="Приватный опрос доступен только по ссылке">ℹ️</span>
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
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    onFocus={(e) => {
                      const rect = e.target.getBoundingClientRect();
                      setTooltip({
                        text: "Опрос будет автоматически закрыт после этой даты",
                        x: rect.left,
                        y: rect.bottom + 5
                      });
                    }}
                    onBlur={() => setTooltip({ text: '', x: 0, y: 0 })}
                  />
                </div>
              </div>
            </div>

            <div className="mb-10">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-xl font-semibold text-gray-800">
                  Вопросы ({questions.length}/{MAX_QUESTIONS})
                </h3>
                <button
                  type="button"
                  onClick={addQuestion}
                  disabled={questions.length >= MAX_QUESTIONS}
                  className={`flex items-center px-3 py-2 rounded-lg text-sm transition-all ${
                    questions.length >= MAX_QUESTIONS
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                      : 'bg-blue-100 text-blue-700 hover:bg-blue-200 hover:scale-105'
                  }`}
                >
                  <span className="material-symbols-outlined mr-1">add</span>
                  Добавить вопрос
                </button>
              </div>

              {questions.map((question, qIndex) => (
                <div
                  key={question.id}
                  id={`question-${question.id}`}
                  className="mb-6 p-4 border border-gray-200 rounded-lg bg-gray-50"
                  style={{ opacity: 1, transition: 'opacity 0.5s ease' }}
                >
                  <div className="flex justify-between items-center mb-3">
                    <h4 className="text-lg font-medium text-gray-800">Вопрос {qIndex + 1}</h4>
                    {questions.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeQuestion(question.id)}
                        className="text-red-500 hover:text-red-700 transition-colors"
                      >
                        <span className="material-symbols-outlined">delete</span>
                      </button>
                    )}
                  </div>

                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Текст вопроса
                    </label>
                    <input
                      type="text"
                      value={question.title}
                      onChange={(e) => handleQuestionChange(question.id, e.target.value)}
                      required
                      maxLength={MAX_QUESTION_LENGTH}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                      onFocus={(e) => {
                        const rect = e.target.getBoundingClientRect();
                        setTooltip({
                          text: `Максимальная длина: ${MAX_QUESTION_LENGTH} символов`,
                          x: rect.left,
                          y: rect.bottom + 5
                        });
                      }}
                      onBlur={() => setTooltip({ text: '', x: 0, y: 0 })}
                    />
                    <div className="mt-2 flex justify-between items-center">
                      <span className="text-xs text-gray-500">
                        {question.title.length}/{MAX_QUESTION_LENGTH}
                      </span>
                      <div className="w-1/2 h-1 bg-gray-200 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-blue-500 transition-width duration-300"
                          style={{ width: `${(question.title.length / MAX_QUESTION_LENGTH) * 100}%` }}
                        ></div>
                      </div>
                    </div>
                  </div>

                  <div>
                    <div className="flex justify-between items-center mb-3">
                      <label className="block text-sm font-medium text-gray-700">
                        Варианты ответов ({question.answers.length}/{MAX_ANSWERS})
                      </label>
                      <button
                        type="button"
                        onClick={() => addAnswer(question.id)}
                        disabled={question.answers.length >= MAX_ANSWERS}
                        className={`flex items-center px-3 py-2 rounded-lg text-sm transition-all ${
                          question.answers.length >= MAX_ANSWERS
                            ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                            : 'bg-blue-100 text-blue-700 hover:bg-blue-200 hover:scale-105'
                        }`}
                      >
                        <span className="material-symbols-outlined mr-1">add</span>
                        Добавить ответ
                      </button>
                    </div>

                    {question.answers.map((answer, aIndex) => (
                      <div
                        key={answer.id}
                        id={`answer-${answer.id}`}
                        className="flex items-center mb-3"
                        style={{ opacity: 1, transition: 'opacity 0.5s ease' }}
                      >
                        <div className="flex-1">
                          <input
                            type="text"
                            value={answer.value}
                            onChange={(e) => handleAnswerChange(question.id, answer.id, e.target.value)}
                            required
                            maxLength={MAX_ANSWER_LENGTH}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                            onFocus={(e) => {
                              const rect = e.target.getBoundingClientRect();
                              setTooltip({
                                text: `Максимальная длина: ${MAX_ANSWER_LENGTH} символов`,
                                x: rect.left,
                                y: rect.bottom + 5
                              });
                            }}
                            onBlur={() => setTooltip({ text: '', x: 0, y: 0 })}
                          />
                          <div className="mt-2 flex justify-between items-center">
                            <span className="text-xs text-gray-500">
                              {answer.value.length}/{MAX_ANSWER_LENGTH}
                            </span>
                            <div className="w-1/2 h-1 bg-gray-200 rounded-full overflow-hidden">
                              <div
                                className="h-full bg-blue-500 transition-width duration-300"
                                style={{ width: `${(answer.value.length / MAX_ANSWER_LENGTH) * 100}%` }}
                              ></div>
                            </div>
                          </div>
                        </div>
                        {question.answers.length > 2 && (
                          <button
                            type="button"
                            onClick={() => removeAnswer(question.id, answer.id)}
                            className="ml-2 text-red-500 hover:text-red-700 transition-colors"
                          >
                            <span className="material-symbols-outlined">close</span>
                          </button>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={isSubmitting}
                className={`flex items-center justify-center h-12 w-32 px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all ${
                  isSubmitting ? 'opacity-70 cursor-not-allowed' : ''
                }`}
              >
                {isSubmitting ? (
                  <span className="animate-spin material-symbols-outlined">refresh</span>
                ) : (
                  'Создать'
                )}
              </button>
            </div>
          </form>

          <ToastContainer />

          {tooltip.text && (
            <div
              className="tooltip"
              style={{ top: tooltip.y, left: tooltip.x }}
            >
              {tooltip.text}
            </div>
          )}
        </div>
      </div>
    </main>
  );
}