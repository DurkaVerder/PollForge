import { NavLink } from 'react-router-dom';

import Sidebar from '../components/Sidebar';
import PollFeed from '../components/PollFeed';

export default function HomePage() {
  return (
    <main className="px-4 sm:px-6 lg:px-8">
      <div className="max-w-screen-xl mx-auto flex flex-col lg:flex-row gap-6">
        <Sidebar />
        
        <div className="flex-1">
          {/* Мобильные вкладки */}
          <div className="lg:hidden flex bg-white rounded-lg shadow-md mb-6 overflow-hidden">
            <NavLink 
              to="/" 
              className={({isActive}) => 
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'} font-medium`
              }
            >
              Лента
            </NavLink>
            <NavLink
              to="/profile"
              className={({isActive}) => 
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              Профиль
            </NavLink>
            <NavLink
              to="/my-polls"
              className={({isActive}) => 
                `flex-1 p-3 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              Мои опросы
            </NavLink>
          </div>

          <PollFeed />
        </div>
      </div>
    </main>
  );
}
