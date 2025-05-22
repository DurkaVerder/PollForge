import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import defaultAvatar from '../static/img/default-avatar.png';

const API_BASE_URL = 'http://localhost:80/api';

export default function PollCard({ poll }) {
  const navigate = useNavigate();
  const [localPoll, setLocalPoll] = useState(poll);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [creator, setCreator] = useState(null);
  const [isLoadingCreator, setIsLoadingCreator] = useState(true);
  const [comments, setComments] = useState([]);
  const [isCommentsOpen, setIsCommentsOpen] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [editingCommentId, setEditingCommentId] = useState(null);
  const [editCommentText, setEditCommentText] = useState('');
  const [usersInfo, setUsersInfo] = useState({});
  const [isLoadingComments, setIsLoadingComments] = useState(false);

  useEffect(() => {
    const fetchCreator = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/profile/user/${localPoll.creator_id}`, {
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('authToken'),
          },
        });
        if (!response.ok) throw new Error('Failed to fetch creator');
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

  useEffect(() => {
    if (isCommentsOpen) fetchComments();
  }, [isCommentsOpen]);

  const fetchComments = async () => {
    setIsLoadingComments(true);
    try {
      const response = await fetch(`${API_BASE_URL}/comments/forms/${localPoll.id}/comments`, {
        headers: { Authorization: 'Bearer ' + localStorage.getItem('authToken') },
      });
      if (!response.ok) throw new Error('Failed to fetch comments');
      const data = await response.json();
      const commentsData = data.comments || [];
      setComments(commentsData);

      const usersData = {};
      for (const comment of commentsData) {
        if (!usersData[comment.user_id]) {
          try {
            const userResponse = await fetch(`${API_BASE_URL}/profile/user/${comment.user_id}`, {
              headers: { Authorization: 'Bearer ' + localStorage.getItem('authToken') },
            });
            if (userResponse.ok) {
              const userData = await userResponse.json();
              usersData[comment.user_id] = userData;
            }
          } catch (error) {
            console.error(`Error fetching user ${comment.user_id}:`, error);
          }
        }
      }
      setUsersInfo(usersData);
    } catch (error) {
      console.error('Error fetching comments:', error);
      setComments([]);
    } finally {
      setIsLoadingComments(false);
    }
  };

  const handleAddComment = async () => {
    if (!newComment.trim()) return;

    try {
      const response = await fetch(`${API_BASE_URL}/comments/forms/${localPoll.id}/comments`, {
        method: 'POST',
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          description: newComment.trim(),
        }),
      });
      if (!response.ok) throw new Error('Failed to add comment');
      await fetchComments();
      setNewComment('');
      // Обновляем счетчик комментариев
      setLocalPoll(prev => ({
        ...prev,
        count_comments: prev.count_comments + 1
      }));
    } catch (error) {
      console.error('Error adding comment:', error);
    }
  };

  const handleUpdateComment = async (commentId) => {
    if (!editCommentText.trim()) return;

    try {
      const response = await fetch(`${API_BASE_URL}/comments/forms/${localPoll.id}/comments/${commentId}`, {
        method: 'PUT',
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          description: editCommentText.trim(),
        }),
      });
      if (!response.ok) throw new Error('Failed to update comment');
      await fetchComments();
      setEditingCommentId(null);
      setEditCommentText('');
    } catch (error) {
      console.error('Error updating comment:', error);
    }
  };

  const handleDeleteComment = async (commentId) => {
    try {
      const response = await fetch(`${API_BASE_URL}/comments/forms/${localPoll.id}/comments/${commentId}`, {
        method: 'DELETE',
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('authToken'),
        },
      });
      if (!response.ok) throw new Error('Failed to delete comment');
      await fetchComments();
      // Обновляем счетчик комментариев
      setLocalPoll(prev => ({
        ...prev,
        count_comments: prev.count_comments - 1
      }));
    } catch (error) {
      console.error('Error deleting comment:', error);
    }
  };

  const handleAvatarClick = () => {
    if (localPoll.creator_id) {
      if (localPoll.creator_id === parseInt(localStorage.getItem('userId'))) {
        navigate('/profile');
      } else {
        navigate(`/profile/${localPoll.creator_id}`);
      }
    }
  };

  const handleLike = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/like/forms/${localPoll.id}`, {
        method: 'POST',
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          is_liked: !localPoll.likes.is_liked
        }),
      });

      if (!response.ok) throw new Error('Ошибка при отправке лайка');

      setLocalPoll(prev => ({
        ...prev,
        likes: {
          count: prev.likes.is_liked ? prev.likes.count - 1 : prev.likes.count + 1,
          is_liked: !prev.likes.is_liked
        }
      }));
    } catch (error) {
      console.error('Ошибка при лайке:', error);
    }
  };

  const handleVote = async (questionId, answerId, isSelected) => {
    setIsSubmitting(true);
    try {
      const question = localPoll.questions.find((q) => q.id === questionId);
      const prevSelectedAnswer = question?.answers.find((a) => a.is_selected && a.id !== answerId);

      if (prevSelectedAnswer) {
        await fetch(`${API_BASE_URL}/vote/input`, {
          method: 'POST',
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('authToken'),
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: prevSelectedAnswer.id,
            is_up_vote: false,
          }),
        });
      }

      const response = await fetch(`${API_BASE_URL}/vote/input`, {
        method: 'POST',
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('authToken'),
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: answerId,
          is_up_vote: !isSelected,
        }),
      });

      if (!response.ok) throw new Error('Ошибка при отправке голоса');

      setLocalPoll((prevPoll) => {
        const updatedQuestions = prevPoll.questions.map((question) => {
          if (question.id === questionId) {
            const updatedAnswers = question.answers.map((answer) => {
              const newAnswer = {
                ...answer,
                is_selected: false,
                count_votes: answer.is_selected ? answer.count_votes - 1 : answer.count_votes,
              };
              if (answer.id === answerId) {
                return {
                  ...newAnswer,
                  is_selected: !isSelected,
                  count_votes: isSelected ? newAnswer.count_votes : newAnswer.count_votes + 1,
                };
              }
              return newAnswer;
            });

            const newTotalVotes = updatedAnswers.reduce((sum, a) => sum + a.count_votes, 0);
            const answersWithPercent = updatedAnswers.map((answer) => ({
              ...answer,
              percent: newTotalVotes > 0 ? Math.round((answer.count_votes / newTotalVotes) * 100) : 0,
            }));

            return {
              ...question,
              answers: answersWithPercent,
              total_count_votes: newTotalVotes,
            };
          }
          return question;
        });

        return {
          ...prevPoll,
          questions: updatedQuestions,
        };
      });
    } catch (error) {
      console.error('Ошибка при голосовании:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-lg p-6 transform hover:shadow-xl transition-all duration-300 w-full mx-auto">
      <div className="flex items-center mb-4 space-x-2">
        <span className="bg-blue-100 text-blue-800 text-xs font-medium px-3 py-1 rounded-full">
          {localPoll.theme}
        </span>
      </div>

      <div className="flex justify-between items-start mb-6">
        <div className="flex items-center space-x-4">
          {isLoadingCreator ? (
            <div className="h-12 w-12 rounded-full bg-gray-200 animate-pulse"></div>
          ) : (
            <button onClick={handleAvatarClick} className="focus:outline-none">
              <img
                src={creator?.avatar_url || defaultAvatar}
                alt="Аватар"
                className="h-12 w-12 rounded-full object-cover hover:ring-2 hover:ring-blue-500 transition-all duration-200"
                onError={(e) => (e.target.src = defaultAvatar)}
              />
            </button>
          )}
          <div>
            <button
              onClick={handleAvatarClick}
              className="font-semibold text-gray-800 hover:text-blue-600 transition-colors"
            >
              {isLoadingCreator ? 'Загрузка...' : creator?.name || 'Аноним'}
            </button>
            <p className="text-sm text-gray-500">
              {new Date(localPoll.created_at).toLocaleDateString('ru-RU')}
            </p>
          </div>
        </div>
        <button className="text-gray-400 hover:text-gray-600 transition-colors">
          <span className="material-symbols-outlined">more_vert</span>
        </button>
      </div>

      <h3 className="text-2xl font-bold text-gray-900 mb-3">{localPoll.title}</h3>
      <p className="text-gray-600 mb-6 leading-relaxed">{localPoll.description}</p>

      {localPoll.questions.map((question, qIndex) => (
        <div key={qIndex} className={qIndex < localPoll.questions.length - 1 ? 'border-b pb-6 mb-6' : 'mb-6'}>
          <h4 className="text-lg font-semibold text-gray-800 mb-4">{question.title}</h4>
          <div className="space-y-4">
            {question.answers.map((answer, aIndex) => (
              <div key={aIndex} className="flex items-center">
                <input
                  type="radio"
                  id={`poll${localPoll.id}_q${qIndex}_answer${aIndex}`}
                  name={`poll${localPoll.id}_q${qIndex}`}
                  className="h-5 w-5 text-blue-600 focus:ring-blue-500"
                  checked={answer.is_selected}
                  onChange={() => handleVote(question.id, answer.id, answer.is_selected)}
                  disabled={isSubmitting}
                />
                <label
                  htmlFor={`poll${localPoll.id}_q${qIndex}_answer${aIndex}`}
                  className="ml-3 block w-full cursor-pointer hover:bg-gray-50 p-2 rounded transition-colors duration-200"
                >
                  <div className="flex justify-between items-center">
                    <span className="text-gray-700">{answer.title}</span>
                    <span className="text-sm text-gray-500">
                      {answer.count_votes} ({answer.percent}%)
                    </span>
                  </div>
                  <div className="mt-2 h-2 w-full bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-primary-500 rounded-full transition-all duration-300"
                      style={{ width: `${answer.percent}%` }}
                    />
                  </div>
                </label>
              </div>
            ))}
          </div>
        </div>
      ))}

      <div className="flex justify-between text-sm text-gray-500 mb-6">
        <span>{localPoll.questions.reduce((sum, q) => sum + q.total_count_votes, 0)} голосов</span>
        <span>Заканчивается {new Date(localPoll.expires_at).toLocaleDateString('ru-RU')}</span>
      </div>

      <div className="flex justify-between items-center">
        <button
          className="flex items-center text-primary-600 hover:text-primary-800 transition-colors duration-200"
          onClick={() => setIsCommentsOpen(!isCommentsOpen)}
        >
          <span className="material-symbols-outlined mr-2">comment</span>
          {localPoll.count_comments} комментариев
        </button>
        <div className="flex items-center space-x-4">
          <button 
            className="flex items-center text-primary-600 hover:text-primary-800 transition-colors duration-200"
            onClick={handleLike}
          >
            <span className="material-symbols-outlined mr-2">
              {localPoll.likes.is_liked ? 'favorite' : 'favorite_border'}
            </span>
            {localPoll.likes.count}
          </button>
        </div>
      </div>

      {isCommentsOpen && (
        <div className="mt-6 border-t pt-6">
          <div className="mb-6">
            <textarea
              className="w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              rows="3"
              placeholder="Напишите ваш комментарий..."
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
            />
            <button
              className="mt-3 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition-colors duration-200"
              onClick={handleAddComment}
            >
              Отправить
            </button>
          </div>

          {isLoadingComments ? (
            <div className="flex justify-center py-4">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
            </div>
          ) : comments.length > 0 ? (
            <div className="space-y-4">
              {comments.map((comment) => {
                const userInfo = usersInfo[comment.user_id] || null;
                return (
                  <div key={comment.id} className="border-b pb-4 last:border-b-0">
                    <div className="flex items-start mb-2">
                      <img
                        src={userInfo?.avatar_url || defaultAvatar}
                        alt="Аватар"
                        className="h-8 w-8 rounded-full object-cover mr-3"
                        onError={(e) => (e.target.src = defaultAvatar)}
                      />
                      <div className="flex-1">
                        <div className="flex justify-between items-center">
                          <span className="font-semibold text-gray-800">{userInfo?.name || 'Аноним'}</span>
                          <span className="text-xs text-gray-500">
                            {new Date(comment.created_at).toLocaleString('ru-RU')}
                          </span>
                        </div>
                        {editingCommentId === comment.id ? (
                          <div className="mt-2">
                            <textarea
                              className="w-full p-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                              rows="2"
                              value={editCommentText}
                              onChange={(e) => setEditCommentText(e.target.value)}
                            />
                            <div className="flex space-x-2 mt-2">
                              <button
                                className="bg-primary-500 text-white px-3 py-1 rounded-lg hover:bg-primary-600 transition-colors"
                                onClick={() => handleUpdateComment(comment.id)}
                              >
                                Сохранить
                              </button>
                              <button
                                className="bg-gray-500 text-white px-3 py-1 rounded-lg hover:bg-gray-600 transition-colors"
                                onClick={() => {
                                  setEditingCommentId(null);
                                  setEditCommentText('');
                                }}
                              >
                                Отмена
                              </button>
                            </div>
                          </div>
                        ) : (
                          <>
                            <p className="text-gray-700 mt-1">{comment.description}</p>
                            {comment.user_id === parseInt(localStorage.getItem('userId')) && (
                              <div className="flex space-x-3 mt-2">
                                <button
                                  className="text-primary-500 hover:text-blue-600 transition-colors"
                                  onClick={() => {
                                    setEditingCommentId(comment.id);
                                    setEditCommentText(comment.description);
                                  }}
                                >
                                  <span className="material-symbols-outlined text-sm">edit</span>
                                </button>
                                <button
                                  className="text-red-500 hover:text-red-600 transition-colors"
                                  onClick={() => handleDeleteComment(comment.id)}
                                >
                                  <span className="material-symbols-outlined text-sm">delete</span>
                                </button>
                              </div>
                            )}
                          </>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <p className="text-gray-500">Комментариев пока нет</p>
          )}
        </div>
      )}
    </div>
  );
}