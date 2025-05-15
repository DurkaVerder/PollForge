import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import './styles/style.css';
import Header from './components/Header';
import Footer from './components/Footer';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import HomePage from './pages/HomePage';
import ProfilePage from './pages/ProfilePage';
import CreatePollPage from './pages/CreatePollPage';
import MyPollsPage from './pages/MyPollsPage';
// import ExplorePage from './pages/ExplorePage';
// import TrendingPage from './pages/TrendingPage';
import NotFoundPage from './pages/NotFoundPage';


const MainLayout = () => (
  <div className="w-[1200px] p-4 font-sans bg-gray-50 min-h-screen">
    <Header />
    <Outlet />
    <Footer />
  </div>
);


const AuthLayout = () => (
  <div className="w-[1200px] p-4 font-sans bg-gray-50 min-h-screen">
    <Outlet />
  </div>
);

export default function App() {
  return (
    <Router>
      <div id="webcrumbs">
        <Routes>
          {/* Маршруты с Header и Footer */}
          <Route element={<MainLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/create-poll" element={<CreatePollPage />} />
            <Route path="/my-polls" element={<MyPollsPage />} />
            {/* <Route path="/explore" element={<ExplorePage />} />
            <Route path="/trending" element={<TrendingPage />} /> */}
            <Route path="*" element={<NotFoundPage />} />
          </Route>
          {/* Маршруты без Header и Footer */}
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>
        </Routes>
      </div>
    </Router>
  );
}