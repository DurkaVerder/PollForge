import PollCard from './PollCard';

export default function PollFeed() {
  const polls = [
    {
      id: 1,
      categories: ['Technology', 'Programming'],
      user: {
        name: 'Alice Johnson',
        avatar: 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwyfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080',
        time: '3 hours ago'
      },
      title: "What's your favorite programming language?",
      questions: [
        {
          title: "What's your favorite programming language?",
          options: [
            { label: 'JavaScript', percentage: 45 },
            { label: 'Python', percentage: 30 },
            { label: 'Java', percentage: 15 },
            { label: 'C++', percentage: 10 }
          ]
        },
        {
          title: "What's your experience level?",
          options: [
            { label: 'Beginner (0-2 years)', percentage: 25 },
            { label: 'Intermediate (3-5 years)', percentage: 40 },
            { label: 'Advanced (6+ years)', percentage: 35 }
          ]
        }
      ],
      votes: 256,
      endsIn: '2 days',
      comments: 42
    },
    {
      id: 2,
      categories: ['Work Environment', 'Workplace'],
      user: {
        name: 'Mark Wilson',
        avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwzfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080',
        time: 'Yesterday'
      },
      title: 'How do you prefer to work?',
      questions: [
        {
          title: 'How do you prefer to work?',
          options: [
            { label: 'Remote', percentage: 60 },
            { label: 'Office', percentage: 25 },
            { label: 'Hybrid', percentage: 15 }
          ]
        }
      ],
      votes: 512,
      endsIn: '5 days',
      comments: 78
    }
  ];

  return (
    <section className="mb-8" id="feed">
      <h2 className="text-2xl font-bold mb-4">Polls Feed</h2>
      <div className="space-y-6">
        {polls.map(poll => (
          <PollCard key={poll.id} poll={poll} />
        ))}
      </div>
    </section>
  );
}