import PollCard from './PollCard';

export default function PollFeed({ polls, hasMore, onLoadMore, isLoadingMore }) {
  return (
    <>
      <style>
        {`
          @keyframes scaleIn {
            0% {
              opacity: 0;
              transform: scale(0.8);
            }
            100% {
              opacity: 1;
              transform: scale(1);
            }
          }
          .animate-scale-in {
            animation: scaleIn 0.4s ease-out forwards;
          }
        `}
      </style>
      <section className="mb-8" id="feed">
        <h2 style={{ fontSize: '1.5rem', fontWeight: 'bold', marginBottom: '1rem' }}>
          Лента опросов
        </h2>
        {polls.length === 0 ? (
          <p className="text-gray-500">
            Опросов нет. Создайте свой первый опрос или подождите, пока другие пользователи опубликуют опросы.
          </p>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
            {polls.map((poll, index) => (
              <div
                key={poll.id}
                className="animate-scale-in"
                style={{ animationDelay: `${index * 150}ms` }}
              >
                <PollCard poll={poll} />
              </div>
            ))}
            {hasMore && (
              <button
                onClick={onLoadMore}
                disabled={isLoadingMore}
                className="mt-4 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition-colors"
              >
                {isLoadingMore ? 'Загрузка...' : 'Загрузить еще'}
              </button>
            )}
          </div>
        )}
      </section>
    </>
  );
}