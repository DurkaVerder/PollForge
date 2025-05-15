import { Link } from 'react-router-dom';

export default function Sidebar() {
  return (
    <aside className="hidden lg:block w-64 bg-white rounded-lg shadow-md p-4 h-fit sticky top-24">
      <nav>
        <ul className="space-y-2">
          <li>
            <Link
              to="/"
              className="flex items-center p-3 rounded-lg bg-primary-50 text-primary-700"
            >
              <span className="material-symbols-outlined mr-3">dynamic_feed</span>
              Feed
            </Link>
          </li>
          <li>
            <Link
              to="/profile"
              className="flex items-center p-3 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            >
              <span className="material-symbols-outlined mr-3">person</span>
              My Profile
            </Link>
          </li>
          <li>
            <Link
              to="/my-polls"
              className="flex items-center p-3 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            >
              <span className="material-symbols-outlined mr-3">poll</span>
              My Polls
            </Link>
          </li>
          <li>
            <Link
              to="/explore"
              className="flex items-center p-3 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            >
              <span className="material-symbols-outlined mr-3">explore</span>
              Explore
            </Link>
          </li>
          <li>
            <Link
              to="/trending"
              className="flex items-center p-3 rounded-lg hover:bg-gray-100 transition-colors duration-200"
            >
              <span className="material-symbols-outlined mr-3">trending_up</span>
              Trending
            </Link>
          </li>
        </ul>
      </nav>
    </aside>
  );
}
