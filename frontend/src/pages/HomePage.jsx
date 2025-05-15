import { NavLink } from 'react-router-dom';

import Sidebar from '../components/Sidebar';
import PollFeed from '../components/PollFeed';

export default function HomePage() {
  return (
    <main className="flex flex-col lg:flex-row gap-6">
      <Sidebar />
      
      <div className="flex-1">
        {/* Mobile Tabs */}
        <div className="lg:hidden flex bg-white rounded-lg shadow-md mb-6 overflow-hidden">
          <NavLink 
            to="/" 
            className={({isActive}) => 
              `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'} font-medium`
            }
          >
            Feed
          </NavLink>
          <NavLink
            to="/profile"
            className={({isActive}) => 
              `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
            }
          >
            Profile
          </NavLink>
          <NavLink
            to="/my-polls"
            className={({isActive}) => 
              `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
            }
          >
            My Polls
          </NavLink>
        </div>
        
        <PollFeed />
      </div>
    </main>
  );
}