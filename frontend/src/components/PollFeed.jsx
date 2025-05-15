import PollCard from './PollCard';

export default function PollFeed() {
  const polls = [
    {
      id: 1,
      categories: ['Технологии', 'Программирование'],
      user: {
        name: 'Алиса Джонсон',
        avatar: 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwyfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080',
        time: '3 часа назад'
      },
      title: "Какой ваш любимый язык программирования?",
      questions: [
        {
          title: "Какой ваш любимый язык программирования?",
          options: [
            { label: 'JavaScript', percentage: 45 },
            { label: 'Python', percentage: 30 },
            { label: 'Java', percentage: 15 },
            { label: 'C++', percentage: 10 }
          ]
        },
        {
          title: "Какой у вас уровень опыта?",
          options: [
            { label: 'Начинающий (0-2 года)', percentage: 25 },
            { label: 'Средний (3-5 лет)', percentage: 40 },
            { label: 'Продвинутый (6+ лет)', percentage: 35 }
          ]
        }
      ],
      votes: 256,
      endsIn: 'через 2 дня',
      comments: 42
    },
    {
      id: 2,
      categories: ['Рабочая среда', 'Рабочее место'],
      user: {
        name: 'Марк Уилсон',
        avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHwzfHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080',
        time: 'Вчера'
      },
      title: 'Как вы предпочитаете работать?',
      questions: [
        {
          title: 'Как вы предпочитаете работать?',
          options: [
            { label: 'Удалённо', percentage: 60 },
            { label: 'В офисе', percentage: 25 },
            { label: 'Гибридный формат', percentage: 15 }
          ]
        }
      ],
      votes: 512,
      endsIn: 'через 5 дней',
      comments: 78
    }
  ];

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
