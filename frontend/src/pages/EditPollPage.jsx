import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import Select from 'react-select';

export default function EditPollPage() {
  const MAX_TITLE_LENGTH = 100;
  const MAX_DESCRIPTION_LENGTH = 250;

  const { id } = useParams();
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    private_key: false,
    expires_at: '',
    theme_id: 1
  });
  const [questions, setQuestions] = useState([]);
  const [themes, setThemes] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingThemes, setIsLoadingThemes] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const customStyles = {
  control: (provided, state) => ({
    ...provided,
    borderColor: state.isFocused ? '#3b82f6' : '#d1d5db', // blue-500 при фокусе, gray-300 по умолчанию
    borderWidth: '1px',
    borderRadius: '0.5rem', // rounded-lg
    padding: '0.5rem', // py-2 px-4
    backgroundColor: '#fff', // bg-white
    boxShadow: state.isFocused ? '0 0 0 2px rgba(59, 130, 246, 0.5)' : 'none', // focus:ring-2 focus:ring-blue-500
    '&:hover': {
      borderColor: '#3b82f6', // hover:border-blue-500
    },
    transition: 'all 0.3s ease', // transition-all
    fontSize: '0.875rem', // text-sm
    fontWeight: 500, // font-medium
    minHeight: '2.5rem', // h-10
  }),
  menu: (provided) => ({
    ...provided,
    borderRadius: '0.5rem', // rounded-lg
    backgroundColor: '#fff', // bg-white
    boxShadow: '0 10px 15px rgba(0, 0, 0, 0.1)', // shadow-lg
    marginTop: '0.25rem', // mt-1
    zIndex: 10,
  }),
  option: (provided, state) => ({
    ...provided,
    backgroundColor: state.isFocused ? '#eff6ff' : '#fff', // bg-blue-50 при наведении, иначе bg-white
    color: state.isSelected ? '#8b5cf6' : '#1f2937', // purple-500 для выбранной, text-gray-800 для остальных
    padding: '0.5rem 1rem', // py-2 px-4
    fontSize: '0.875rem', // text-sm
    cursor: 'pointer',
    transition: 'background-color 0.2s ease',
    whiteSpace: 'normal', // для переноса длинных описаний
    '&:hover': {
      backgroundColor: '#dbeafe', // hover:bg-blue-100
    },
  }),
  singleValue: (provided) => ({
    ...provided,
    color: '#8b5cf6', // purple-500 для выбранной темы
    fontSize: '0.875rem', // text-sm
    whiteSpace: 'normal', // для переноса текста
    maxWidth: '100%', // чтобы текст не обрезался
  }),
  placeholder: (provided) => ({
    ...provided,
    color: '#9ca3af', // text-gray-400
    fontSize: '0.875rem', // text-sm
  }),
  dropdownIndicator: (provided, state) => ({
    ...provided,
    color: state.isFocused ? '#3b82f6' : '#6b7280', // blue-500 или gray-500
    '&:hover': {
      color: '#3b82f6', // hover:text-blue-500
    },
  }),
  clearIndicator: (provided) => ({
    ...provided,
    color: '#6b7280', // gray-500
    '&:hover': {
      color: '#ef4444', // hover:text-red-500
    },
  }),
  input: (provided) => ({
    ...provided,
    color: '#1f2937', // text-gray-800
    fontSize: '0.875rem', // text-sm
  }),
};

  useEffect(() => {
    const fetchThemes = async () => {
      try {
        const response = await fetch('http://localhost:80/api/forms/themes', {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch themes');
        }

        const data = await response.json();
        setThemes(data);
        setIsLoadingThemes(false);
      } catch (err) {
        setError(err.message);
        console.error('Fetch themes error:', err);
        setIsLoadingThemes(false);
      }
    };

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

        const expiresAt = data.expires_at 
          ? new Date(data.expires_at).toISOString().slice(0, 16)
          : '';

        setFormData({
          title: data.title,
          description: data.description,
          private_key: data.private_key,
          confidential: data.confidential,
          expires_at: expiresAt,
          theme_id: data.theme_id || 1
        });

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

    fetchThemes();
    fetchFormData();
  }, [id]);

  const handleFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleThemeChange = (e) => {
    const themeId = parseInt(e.target.value);
    setFormData(prev => ({
      ...prev,
      theme_id: themeId
    }));
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
          confidential: formData.confidential,
          private_key: formData.private_key,
          expires_at: utcISOString,
          theme_id: formData.theme_id,
        })
      });

      if (!formResponse.ok) {
        const errorData = await formResponse.json();
        throw new Error(errorData.message || 'Ошибка при обновлении формы');
      }

      navigate('/my-polls');
    } catch (err) {
      setError(err.message);
      console.error('Update poll error:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading || isLoadingThemes) {
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
              <li>Название опроса: до {MAX_TITLE_LENGTH} символов</li>
              <li>Описание: до {MAX_DESCRIPTION_LENGTH} символов</li>
            </ul>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-md">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit}>
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

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Тема опроса
                  </label>
                  <Select
                    name="theme_id"
                    value={
                      themes.length > 0
                        ? themes
                            .map(theme => ({
                              value: theme.id,
                              label: `${theme.name} - ${theme.description.slice(0, 50)}${theme.description.length > 50 ? '...' : ''}`,
                            }))
                            .find(option => option.value === formData.theme_id) || null
                        : null
                    }
                    onChange={(selectedOption) =>
                      handleThemeChange({ target: { value: selectedOption ? selectedOption.value : 1 } })
                    }
                    options={themes.map(theme => ({
                      value: theme.id,
                      label: `${theme.name} - ${theme.description.slice(0, 50)}${theme.description.length > 50 ? '...' : ''}`,
                    }))}
                    className="w-full"
                    classNamePrefix="select"
                    styles={customStyles}
                    placeholder="Выберите тему"
                    isSearchable={true}
                    isDisabled={themes.length === 0}
                    getOptionLabel={(option) => option.label}
                    getOptionValue={(option) => option.value}
                    formatOptionLabel={(option) =>
                      option.label ? (
                        <div className="flex flex-col">
                          <span className="font-medium">{option.label.split(' - ')[0]}</span>
                          <span className="text-xs text-gray-500">{option.label.split(' - ')[1] || 'Нет описания'}</span>
                        </div>
                      ) : (
                        <div className="flex flex-col">
                          <span className="font-medium">Без названия</span>
                          <span className="text-xs text-gray-500">Нет описания</span>
                        </div>
                      )
                    }
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

                <div className="flex items-center">
                    <input
                      type="checkbox"
                      id="confidential"
                      name="confidential"
                      checked={formData.confidential}
                      onChange={handleFormChange}
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    />
                    <label htmlFor="confidential" className="ml-2 block text-sm text-gray-700">
                      Скрыть результаты от других пользователей
                    </label>
                    <span className="ml-2 text-gray-500" title="Результаты будут видны только вам">ℹ️</span>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Дата окончания 
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

            <div className="mb-8">
              <h3 className="text-lg font-semibold mb-4">Вопросы</h3>
              {questions.map((question, qIndex) => (
                <div key={question.id} className="mb-6 p-4 border border-gray-200 rounded-lg">
                  <h4 className="font-medium mb-3">Вопрос {qIndex + 1}</h4>
                  <p className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100">{question.title}</p>

                  <div className="mt-4">
                    <label className="block text-sm font-medium text-gray-700 mb-3">
                      Варианты ответов
                    </label>
                    {question.answers.map((answer, aIndex) => (
                      <div key={answer.id} className="mb-2">
                        <p className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100">{answer.title}</p>
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