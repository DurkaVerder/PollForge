import { useState } from 'react';
import { Link } from 'react-router-dom';
import { FiBarChart2, FiEdit, FiShare2, FiLock, FiAward, FiStar } from 'react-icons/fi';
import { FaTelegram, FaVk, FaGithub } from 'react-icons/fa';
import { motion } from 'framer-motion';
import analyticsDemo from '../static/img/analytics-demo.png';
import creatorDemo from '../static/img/creator-demo.gif';
import hero from '../static/img/hero.png';
import user1 from '../static/img/user1.jpg';
import user2 from '../static/img/user2.jpg';
import user3 from '../static/img/user3.jpg';

export default function LandingPage() {
  const [activeFeature, setActiveFeature] = useState(0);

  const features = [
    {
      icon: <FiEdit className="text-3xl mb-4 text-indigo-600" />,
      title: "Простое создание",
      description: "Интуитивный конструктор для быстрых опросов."
    },
    {
      icon: <FiBarChart2 className="text-3xl mb-4 text-indigo-600" />,
      title: "Аналитика в реальном времени",
      description: "Мгновенные результаты с графиками."
    },
    {
      icon: <FiShare2 className="text-3xl mb-4 text-indigo-600" />,
      title: "Мгновенный доступ",
      description: "Делитесь опросами одной ссылкой."
    },
    {
      icon: <FiLock className="text-3xl mb-4 text-indigo-600" />,
      title: "Конфиденциальность",
      description: "Защита данных и анонимность."
    }
  ];

  const testimonials = [
    {
      name: "Алексей Шматов",
      role: "Проктолог",
      text: "Незаменимый инструмент для клиентов!",
      image: user1,
    },
    {
      name: "Никита Трушин",
      role: "Преподаватель",
      text: "Опросы стали увлекательными!",
      image: user2,
    },
    {
      name: "Семён Иванов",
      role: "IT-специалист",
      text: "Простой сбор обратной связи.",
      image: user3,
    }
  ];

  const currentYear = new Date().getFullYear();

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <header className="min-h-screen py-24 px-4 text-center relative overflow-hidden">
        <div className="absolute inset-0 bg-indigo-100 opacity-50"></div>
        <div className="relative z-10 max-w-4xl mx-auto">
          <motion.h1 
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-5xl md:text-6xl font-bold text-indigo-900 mb-6"
          >
            Создавайте опросы <span className="text-indigo-600">без усилий</span>
          </motion.h1>
          <motion.p 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            className="text-xl text-gray-700 mb-10"
          >
            PollForge — инструмент для анализа данных
          </motion.p>
          <motion.div 
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.4 }}
            className="flex justify-center space-x-4"
          >
            <Link 
              to="/register" 
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-3 rounded-full font-bold text-lg transition-colors shadow-lg"
            >
              Начать бесплатно
            </Link>
          </motion.div>
        </div>
        <motion.img 
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 1, delay: 0.6 }}
          src={hero}
          alt="PollForge в действии" 
          className="mt-10 mx-auto max-w-2xl rounded-3xl shadow-2xl"
        />
      </header>

      {/* Features Showcase */}
      <section className="py-16 px-4 bg-white">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-4xl font-bold text-center mb-16 text-indigo-900">Почему PollForge?</h2>
          
          <div className="grid md:grid-cols-2 gap-12">
            <div className="space-y-8">
              {features.map((feature, index) => (
                <motion.div 
                  key={index}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.5, delay: index * 0.1 }}
                  className={`p-6 rounded-3xl cursor-pointer transition-all ${activeFeature === index ? 'bg-indigo-50 border-l-4 border-indigo-600' : 'hover:bg-gray-100'}`}
                  onClick={() => setActiveFeature(index)}
                >
                  <div className="flex items-start">
                    <div className="mr-4">
                      {feature.icon}
                    </div>
                    <div>
                      <h3 className="text-xl font-semibold mb-2 text-indigo-800">{feature.title}</h3>
                      <p className="text-gray-700">{feature.description}</p>
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
            
            <div className="flex items-center justify-center">
              <motion.div 
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5 }}
                className="bg-gray-100 rounded-3xl p-8 w-full h-80 flex items-center justify-center shadow-lg"
              >
                {activeFeature === 0 && (
                  <img src={creatorDemo} alt="Создание опроса" className="rounded-2xl shadow-lg" />
                )}
                {activeFeature === 1 && (
                  <img src={analyticsDemo} alt="Аналитика" className="rounded-2xl shadow-lg" />
                )}
                {activeFeature === 2 && (
                  <img src="../static/img/share-demo.gif" alt="Поделиться опросом" className="rounded-2xl shadow-lg" />
                )}
                {activeFeature === 3 && (
                  <img src="../static/img/security-demo.png" alt="Безопасность" className="rounded-2xl shadow-lg" />
                )}
              </motion.div>
            </div>
          </div>
        </div>
      </section>

      {/* Testimonials */}
      <section className="py-20 px-4 bg-indigo-50">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-4xl font-bold text-center mb-16 text-indigo-900">Отзывы</h2>
          <div className="grid md:grid-cols-3 gap-8">
            {testimonials.map((testimonial, index) => (
              <motion.div 
                key={index}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: index * 0.1 }}
                className="bg-white p-6 rounded-3xl shadow-md hover:shadow-lg transition-shadow"
              >
                <div className="flex items-center mb-4">
                  <img src={testimonial.image} alt={testimonial.name} className="w-12 h-12 rounded-full mr-4" />
                  <div>
                    <h4 className="font-semibold text-indigo-800">{testimonial.name}</h4>
                    <p className="text-sm text-gray-600">{testimonial.role}</p>
                  </div>
                </div>
                <p className="text-gray-700 italic">"{testimonial.text}"</p>
                <div className="flex justify-end mt-4">
                  {[...Array(5)].map((_, i) => (
                    <FiStar key={i} className="text-indigo-400" />
                  ))}
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Use Cases */}
      <section className="py-20 px-4 bg-gray-50">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-4xl font-bold text-center mb-16 text-indigo-900">Для кого?</h2>
          
          <div className="grid md:grid-cols-3 gap-8">
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5 }}
              className="bg-white p-8 rounded-3xl shadow-sm hover:shadow-md transition-shadow"
            >
              <FiAward className="text-3xl mb-4 text-indigo-600" />
              <h3 className="text-xl font-semibold mb-3 text-indigo-800">Маркетологи</h3>
              <p className="text-gray-700">Исследуйте рынок и улучшайте продукты.</p>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
              className="bg-white p-8 rounded-3xl shadow-sm hover:shadow-md transition-shadow"
            >
              <FiAward className="text-3xl mb-4 text-indigo-600" />
              <h3 className="text-xl font-semibold mb-3 text-indigo-800">Преподаватели</h3>
              <p className="text-gray-700">Интерактивные опросы для студентов.</p>
            </motion.div>
            
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
              className="bg-white p-8 rounded-3xl shadow-sm hover:shadow-md transition-shadow"
            >
              <FiAward className="text-3xl mb-4 text-indigo-600" />
              <h3 className="text-xl font-semibold mb-3 text-indigo-800">HR-специалисты</h3>
              <p className="text-gray-700">Собирайте обратную связь сотрудников.</p>
            </motion.div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4 bg-indigo-100 text-indigo-900">
        <div className="max-w-4xl mx-auto text-center">
          <motion.h2 
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-4xl font-bold mb-6"
          >
            Готовы начать?
          </motion.h2>
          <motion.p 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            className="text-xl mb-10"
          >
            Начните собирать данные уже сегодня.
          </motion.p>
          <motion.div 
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.4 }}
          >
            <Link 
              to="/register" 
              className="inline-block bg-indigo-600 text-white hover:bg-indigo-700 px-10 py-4 rounded-full font-bold text-lg transition-colors shadow-lg"
            >
              Зарегистрироваться
            </Link>
          </motion.div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 rounded-t-3xl shadow-md p-6 mt-6 text-center text-white">
        <p>© {currentYear} PollForge. Все права защищены.</p>
        <div className="flex justify-center space-x-4 mt-4">
          <a href="#" className="text-gray-400 hover:text-white transition-colors duration-300">
            <FaTelegram className="text-lg" />
          </a>
          <a href="#" className="text-gray-400 hover:text-white transition-colors duration-300">
            <FaVk className="text-lg" />
          </a>
          <a href="https://github.com/DurkaVerder/PollForge" className="text-gray-400 hover:text-white transition-colors duration-300">
            <FaGithub className="text-lg" />
          </a>
        </div>
        <div className="mt-4">
          <a href="#" className="text-gray-400 hover:text-white mx-2">Блог</a>
          <a href="#" className="text-gray-400 hover:text-white mx-2">FAQ</a>
          <a href="#" className="text-gray-400 hover:text-white mx-2">Контакты</a>
        </div>
      </footer>
    </div>
  );
}