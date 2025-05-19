export default function Footer() {

  const currentYear = new Date().getFullYear();


  return (
    <footer className="bg-white rounded-lg shadow-md p-6 mt-6 text-center">
      <p>&copy; {currentYear} PollForge. All rights reserved.</p>
      <div className="flex justify-center space-x-4 mt-4">
        <a href="#" className="text-gray-500 hover:text-primary-600 transition-colors duration-300">
          <i className="fa-brands fa-telegram text-lg"></i>
        </a>
        <a href="#" className="text-gray-500 hover:text-primary-600 transition-colors duration-300">
          <i className="fa-brands fa-facebook text-lg"></i>
        </a>
        <a href="#" className="text-gray-500 hover:text-primary-600 transition-colors duration-300">
          <i className="fa-brands fa-instagram text-lg"></i>
        </a>
        <a href="https://github.com/DurkaVerder/PollForge" className="text-gray-500 hover:text-primary-600 transition-colors duration-300">
          <i className="fa-brands fa-github text-lg"></i>
        </a>
      </div>
    </footer>
  );
}