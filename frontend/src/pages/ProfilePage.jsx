import UserProfile from '../components/UserProfile';
import Sidebar from '../components/Sidebar';

export default function ProfilePage() {
  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      
      <div className="flex-1">
        <UserProfile />
      </div>
    </main>
  );
}