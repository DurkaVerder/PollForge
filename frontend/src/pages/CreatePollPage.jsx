import Sidebar from '../components/Sidebar';

export default function CreatePollPage() {
  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      
      <div className="flex-1">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold mb-6">Создать новый опрос</h2>
          {/* Здесь добавьте форму для создания опроса */}
        </div>
      </div>
    </main>
  );
}
