import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
import './styles/style.css';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
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
import ProtectedRoute from './components/ProtectedRoute'; 
import EditProfilePage from './pages/EditProfilePage';
import PollDetailPage from './pages/PollDetailPage';
import AnotherPollPage from './pages/AnotherPollPage';
import EditPollPage from './pages/EditPollPage';
import AnotherProfilePage from './pages/AnotherProfilePage';

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
          {/* Маршруты без Header и Footer */}
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>

          {/* Защищённые маршруты с Header и Footer */}
          <Route element={<MainLayout />}>
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <HomePage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <ProfilePage />
                </ProtectedRoute>
              }
            />
             <Route
              path="/profile/:id"
              element={
                <ProtectedRoute>
                  <AnotherProfilePage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/create-poll"
              element={
                <ProtectedRoute>
                  <CreatePollPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/my-polls"
              element={
                <ProtectedRoute>
                  <MyPollsPage />
                </ProtectedRoute>
              }
            />
            <Route 
              path="/profile/:id/edit"
              element={
                <ProtectedRoute>
                  <EditProfilePage />
                </ProtectedRoute>
              }
            />

            <Route 
              path="/poll/:link"
              element={
                <ProtectedRoute>
                  <PollDetailPage />
                </ProtectedRoute>
              }
            />

            <Route 
              path="/poll/vote/:link"
              element={
                <ProtectedRoute>
                  <AnotherPollPage />
                </ProtectedRoute>
              }
            />

            <Route 
              path="/poll/edit/:id"
              element={
                <ProtectedRoute>
                  <EditPollPage />
                </ProtectedRoute>
              }
            />
            <Route path="/notfound" element={<NotFoundPage />} />
            {/* <Route path="/explore" element={<ProtectedRoute><ExplorePage /></ProtectedRoute>} /> */}
            {/* <Route path="/trending" element={<ProtectedRoute><TrendingPage /></ProtectedRoute>} /> */}
            <Route path="*" element={<NotFoundPage />} />
            
          </Route>
        </Routes>
         <ToastContainer position="top-right" autoClose={2000} />
      </div>
    </Router>
  );
}
