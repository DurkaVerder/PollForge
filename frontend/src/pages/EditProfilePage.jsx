import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

export default function EditProfilePage() {
  const [formData, setFormData] = useState({
    name: '',
    bio: '',
    avatar: null
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const navigate = useNavigate();
  const { id } = useParams();


  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch(`http://localhost:80/api/profile/`, {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken')
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка загрузки профиля');
        }

        const data = await response.json();
        setFormData({
          name: data.name,
          bio: data.bio || '',
          avatar: null
        });
      } catch (err) {
        setError(err.message);
        console.error('Profile fetch error:', err);
      }
    };

    fetchProfile();
  }, [id]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleFileChange = (e) => {
    setFormData(prev => ({
      ...prev,
      avatar: e.target.files[0]
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');
    setSuccess('');

    try {
      // 1. Update name if changed
      if (formData.name !== '') {
        const nameResponse = await fetch('http://localhost:80/api/profile/name', {
          method: 'PUT',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: formData.name
          })
        });

        const nameData = await nameResponse.json();
        if (!nameResponse.ok) {
          throw new Error(nameData.message || 'Ошибка при обновлении имени');
        }
      }

      // 2. Update bio if changed
      if (formData.bio !== '') {
        const bioResponse = await fetch('http://localhost:80/api/profile/bio', {
          method: 'PUT',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            bio: formData.bio
          })
        });

        const bioData = await bioResponse.json();
        if (!bioResponse.ok) {
          throw new Error(bioData.message || 'Ошибка при обновлении описания');
        }
      }

      // 3. Update avatar if selected
      if (formData.avatar) {
        const formDataAvatar = new FormData();
        formDataAvatar.append('avatar', formData.avatar);

        const avatarResponse = await fetch('http://localhost:80/api/profile/avatar', {
          method: 'POST',
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
          },
          body: formDataAvatar
        });

        const avatarData = await avatarResponse.json();
        if (!avatarResponse.ok) {
          throw new Error(avatarData.message || 'Ошибка при загрузке аватарки');
        }
      }

      setSuccess('Профиль успешно обновлён');
      setTimeout(() => {
        navigate(`/profile`);
      }, 1500);
    } catch (err) {
      setError(err.message);
      console.error('Profile update error:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-8 max-w-md mx-auto mt-8">
      <h2 className="text-2xl font-bold mb-4">Редактирование профиля</h2>

      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-md">
          {error}
        </div>
      )}

      {success && (
        <div className="mb-4 p-3 bg-green-50 text-green-600 rounded-md">
          {success}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Имя
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Описание профиля
            </label>
            <textarea
              name="bio"
              value={formData.bio}
              onChange={handleInputChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              rows="3"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Аватар (PNG или JPG)
            </label>
            <input
              type="file"
              accept=".png,.jpg,.jpeg"
              onChange={handleFileChange}
              className="block w-full text-sm text-gray-500
                file:mr-4 file:py-2 file:px-4
                file:rounded-lg file:border-0
                file:text-sm file:font-semibold
                file:bg-primary-50 file:text-primary-700
                hover:file:bg-primary-100"
            />
          </div>
        </div>

        <div className="mt-6 flex justify-end space-x-3">
          <button
            type="button"
            onClick={() => navigate(`/profile`)}
            className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
            disabled={isSubmitting}
          >
            Отмена
          </button>
          <button
            type="submit"
            className="px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600 disabled:opacity-70"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Сохранение...' : 'Сохранить изменения'}
          </button>
        </div>
      </form>
    </div>
  );
}