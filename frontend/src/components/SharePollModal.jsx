import { useState } from 'react';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';


const SharePollModal = ({ pollLink, isOpen, onClose }) => {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(pollLink);
      setIsCopied(true);
      toast.success('Ссылка скопирована!', { autoClose: 2000 });
      setTimeout(() => setIsCopied(false), 2000);
    } catch (err) {
      console.error('Ошибка копирования:', err);
      toast.error('Не удалось скопировать', { autoClose: 2000 });
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-lg font-bold mb-4">Поделиться опросом</h2>
        <div className="flex items-center space-x-2">
          <input
            type="text"
            value={pollLink}
            readOnly
            className="w-full p-2 border rounded bg-gray-100"
          />
          <button
            onClick={handleCopy}
            className="p-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            title={isCopied ? 'Скопировано!' : 'Копировать'}
          >
            {isCopied ? (
              '✓'
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                />
              </svg>
            )}
          </button>
        </div>
        <button
          onClick={onClose}
          className="mt-4 p-3 w-30 bg-primary-500 rounded-lg hover:bg-primary-600"
        >
          Закрыть
        </button>
      </div>
      <ToastContainer position="top-right" autoClose={2000} />
    </div>
  );
};

export default SharePollModal;