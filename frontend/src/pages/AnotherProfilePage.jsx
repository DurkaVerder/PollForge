
import { useParams} from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import AnotherUserProfile from '../components/AnotherUserProfile';

export default function AnotherProfilePage() {
  const { id } = useParams();

  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      <div className="flex-1">
        <AnotherUserProfile id = {id} />
      </div>
    </main>
  );
}
