import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './styles/style.css';
import Header from './components/Header';
import Footer from './components/Footer';
import HomePage from './pages/HomePage';
import ProfilePage from './pages/ProfilePage';
import CreatePollPage from './pages/CreatePollPage';
import MyPollsPage from './pages/MyPollsPage';
// import ExplorePage from './pages/ExplorePage';
// import TrendingPage from './pages/TrendingPage';
import NotFoundPage from './pages/NotFoundPage';

export default function App() {
  return (
    <Router>
      <div id="webcrumbs">
        <div className="w-[1200px] p-4 font-sans bg-gray-50 min-h-screen">
          <Header />
          
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/create-poll" element={<CreatePollPage />} />
            <Route path="/my-polls" element={<MyPollsPage />} />
            {/* <Route path="/explore" element={<ExplorePage />} />
            <Route path="/trending" element={<TrendingPage />} /> */}
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
          
          <Footer />
        </div>
      </div>
    </Router>
  );
}