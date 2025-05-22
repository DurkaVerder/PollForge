import PollCard from './PollCard';

export default function PollFeed({ polls }) {
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
        </div>
      </section>
    </>
  );
}