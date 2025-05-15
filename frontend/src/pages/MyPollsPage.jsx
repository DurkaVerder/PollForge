import Sidebar from '../components/Sidebar';

export default function MyPollsPage() {
  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      <div className="flex-1">
        <h2 className="text-2xl font-bold mb-4">My Polls</h2>
        {/* Добавьте контент */}
      </div>
    </main>
  );
}