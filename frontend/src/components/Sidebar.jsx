import { NavLink } from 'react-router-dom';

export default function Sidebar() {
  return (
    <aside className="hidden lg:block w-64 bg-white rounded-lg shadow-md p-4 h-fit sticky top-24">
      <nav>
        <ul className="space-y-2">
          <li>
            <NavLink
              to="/stream-line"
              end
              className={({ isActive }) => 
                `flex items-center p-3 rounded-lg transition-colors duration-200 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              <span className="material-symbols-outlined mr-3">dynamic_feed</span>
              Лента
            </NavLink>
          </li>
          <li>
            <NavLink
              to="/profile"
              className={({ isActive }) => 
                `flex items-center p-3 rounded-lg transition-colors duration-200 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              <span className="material-symbols-outlined mr-3">person</span>
              Мой профиль
            </NavLink>
          </li>
          <li>
            <NavLink
              to="/my-polls"
              className={({ isActive }) => 
                `flex items-center p-3 rounded-lg transition-colors duration-200 ${isActive ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-100'}`
              }
            >
              <span className="material-symbols-outlined mr-3">poll</span>
              Мои опросы
            </NavLink>
          </li>
          
        </ul>
      </nav>
    </aside>
  );
}