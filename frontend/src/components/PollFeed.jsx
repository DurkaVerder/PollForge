import PollCard from './PollCard';

export default function PollFeed({ polls }) {
  return (
    <section className="mb-8" id="feed">
      <h2 className="text-2xl font-bold mb-4">Лента опросов</h2>
      <div className="space-y-6">
        {polls.map(poll => (
          <PollCard key={poll.id} poll={poll} />
        ))}
      </div>
    </section>
  );
}
