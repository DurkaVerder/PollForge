import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Sidebar from '../components/Sidebar';

export default function CreatePollPage() {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    private_key: false,
    expires_at: ''
  });
  const [questions, setQuestions] = useState([
    { title: '', answers: [''] }
  ]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleQuestionChange = (index, value) => {
    const newQuestions = [...questions];
    newQuestions[index].title = value;
    setQuestions(newQuestions);
  };

  const handleAnswerChange = (qIndex, aIndex, value) => {
    const newQuestions = [...questions];
    newQuestions[qIndex].answers[aIndex] = value;
    setQuestions(newQuestions);
  };

  const addQuestion = () => {
    setQuestions([...questions, { title: '', answers: [''] }]);
  };

  const removeQuestion = (index) => {
    if (questions.length > 1) {
      const newQuestions = [...questions];
      newQuestions.splice(index, 1);
      setQuestions(newQuestions);
    }
  };

  const addAnswer = (qIndex) => {
    const newQuestions = [...questions];
    newQuestions[qIndex].answers.push('');
    setQuestions(newQuestions);
  };

  const removeAnswer = (qIndex, aIndex) => {
    const newQuestions = [...questions];
    if (newQuestions[qIndex].answers.length > 1) {
      newQuestions[qIndex].answers.splice(aIndex, 1);
      setQuestions(newQuestions);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');

    try {
      const date = new Date(formData.expires_at);
      const utcISOString = date.toISOString();
      // 1. Создаем форму
      const formResponse = await fetch('http://localhost:80/api/forms/', {
        method: 'POST',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        

        body: JSON.stringify({
          title: formData.title,
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

      // 2. Создаем вопросы
      for (let qIndex = 0; qIndex < questions.length; qIndex++) {
        const question = questions[qIndex];
        
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

        // 3. Создаем ответы для вопроса
        for (let aIndex = 0; aIndex < question.answers.length; aIndex++) {
          const answer = question.answers[aIndex];
          
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

      
      navigate('/my-polls');
    } catch (err) {
      setError(err.message);
      console.error('Create poll error:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      
      <div className="flex-1">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold mb-6">Создать новый опрос</h2>
          
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
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Название опроса
                  </label>
                  <input
                    type="text"
                    name="title"
                    value={formData.title}
                    onChange={handleFormChange}
                    required
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Описание
                  </label>
                  <textarea
                    name="description"
                    value={formData.description}
                    onChange={handleFormChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    rows="3"
                  />
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
                <h3 className="text-lg font-semibold">Вопросы</h3>
                <button
                  type="button"
                  onClick={addQuestion}
                  className="px-3 py-1 bg-primary-50 text-primary-600 rounded-lg text-sm hover:bg-primary-100"
                >
                  + Добавить вопрос
                </button>
              </div>

              {questions.map((question, qIndex) => (
                <div key={qIndex} className="mb-6 p-4 border border-gray-200 rounded-lg">
                  <div className="flex justify-between items-center mb-3">
                    <h4 className="font-medium">Вопрос {qIndex + 1}</h4>
                    {questions.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeQuestion(qIndex)}
                        className="text-red-500 text-sm hover:text-red-700"
                      >
                        Удалить вопрос
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
                      onChange={(e) => handleQuestionChange(qIndex, e.target.value)}
                      required
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    />
                  </div>

                  <div>
                    <div className="flex justify-between items-center mb-3">
                      <label className="block text-sm font-medium text-gray-700">
                        Варианты ответов
                      </label>
                      <button
                        type="button"
                        onClick={() => addAnswer(qIndex)}
                        className="px-3 py-1 bg-primary-50 text-primary-600 rounded-lg text-sm hover:bg-primary-100"
                      >
                        + Добавить ответ
                      </button>
                    </div>

                    {question.answers.map((answer, aIndex) => (
                      <div key={aIndex} className="flex items-center mb-2">
                        <input
                          type="text"
                          value={answer}
                          onChange={(e) => handleAnswerChange(qIndex, aIndex, e.target.value)}
                          required
                          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                        />
                        {question.answers.length > 1 && (
                          <button
                            type="button"
                            onClick={() => removeAnswer(qIndex, aIndex)}
                            className="ml-2 text-red-500 hover:text-red-700"
                          >
                            ×
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
                className={`px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
                  isSubmitting ? 'opacity-70 cursor-not-allowed' : ''
                }`}
              >
                {isSubmitting ? 'Создание...' : 'Создать опрос'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  );
}