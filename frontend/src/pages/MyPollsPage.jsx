import Sidebar from '../components/Sidebar';

export default function MyPollsPage() {
  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      <div className="flex-1">
        <h2 className="text-2xl font-bold mb-4">Мои опросы</h2>
        {/* Здесь добавьте содержимое страницы с вашими опросами */}
      </div>
    </main>
  );
}
