import { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { Bar, Line, Pie } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
} from 'chart.js';
import 'chartjs-adapter-date-fns';
import { FiDownload } from 'react-icons/fi';
import ChartDataLabels from 'chartjs-plugin-datalabels';

// Register Chart.js components and plugins
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  LineElement,
  PointElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  ChartDataLabels
);

// Plugin for white background on chart download
const whiteBackgroundPlugin = {
  id: 'whiteBackground',
  beforeDraw: (chart) => {
    const ctx = chart.ctx;
    ctx.save();
    ctx.globalCompositeOperation = 'destination-over';
    ctx.fillStyle = 'white';
    ctx.fillRect(0, 0, chart.width, chart.height);
    ctx.restore();
  },
};

// Chart data functions
const getCountChartData = (pollStat) => ({
  labels: pollStat.answers.map((answer) => answer.title),
  datasets: [
    {
      label: 'Количество ответов',
      data: pollStat.answers.map((answer) => answer.count),
      backgroundColor: 'rgba(54, 162, 235, 0.5)',
      borderColor: 'rgba(54, 162, 235, 1)',
      borderWidth: 1,
    },
  ],
});

const getTimeChartData = (pollStat, responseHistory, answerColors) => {
  const cumulativeData = {};
  responseHistory.forEach((entry) => {
    const stat = entry.stats.find((s) => s.question_title === pollStat.question_title);
    if (stat) {
      stat.answers.forEach((answer) => {
        const key = answer.title;
        if (!cumulativeData[key]) {
          cumulativeData[key] = { timestamps: [], counts: [] };
        }
        cumulativeData[key].timestamps.push(entry.timestamp);
        cumulativeData[key].counts.push(answer.count);
      });
    }
  });

  const datasets = Object.keys(cumulativeData).map((key) => {
    if (!answerColors[key]) {
      answerColors[key] = `#${Math.floor(Math.random() * 16777215).toString(16)}`;
    }
    return {
      label: key,
      data: cumulativeData[key].timestamps.map((timestamp, index) => ({
        x: timestamp,
        y: cumulativeData[key].counts[index],
      })),
      borderColor: answerColors[key],
      backgroundColor: 'rgba(0, 0, 0, 0)',
      borderWidth: 2,
      tension: 0.1,
    };
  });

  return { datasets };
};

const getPieChartData = (pollStat) => ({
  labels: pollStat.answers.map((answer) => answer.title),
  datasets: [
    {
      data: pollStat.answers.map((answer) => answer.count),
      backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'].slice(
        0,
        pollStat.answers.length
      ),
    },
  ],
});

// Chart options
const countChartOptions = {
  responsive: true,
  plugins: {
    legend: { position: 'top' },
    title: { display: true, text: 'Статистика ответов' },
    datalabels: { display: false },
  },
  scales: {
    y: { beginAtZero: true, ticks: { precision: 0 } },
  },
};

const timeChartOptions = {
  responsive: true,
  plugins: {
    legend: { position: 'top' },
    title: { display: true, text: 'Временная активность ответов' },
    datalabels: { display: false },
  },
  scales: {
    x: {
      type: 'time',
      time: { unit: 'minute', displayFormats: { minute: 'dd.MM.yyyy HH:mm' } },
      title: { display: true, text: 'Время' },
    },
    y: {
      beginAtZero: true,
      title: { display: true, text: 'Количество ответов' },
      ticks: { precision: 0 },
    },
  },
};

const getPieChartOptions = (showPercentages) => ({
  responsive: true,
  plugins: {
    legend: { position: 'top' },
    title: { display: true, text: 'Распределение ответов' },
    datalabels: {
      display: showPercentages,
      formatter: (value, context) => {
        const total = context.chart.data.datasets[0].data.reduce((a, b) => a + b, 0);
        const percentage = (value / total * 100).toFixed(1);
        return percentage > 0 ? `${percentage}%` : '';
      },
      color: 'white',
      font: { weight: 'bold' },
    },
    tooltip: {
      callbacks: {
        label: (context) => {
          if (!showPercentages) return context.label;
          const dataset = context.dataset;
          const total = dataset.data.reduce((a, b) => a + b, 0);
          const currentValue = dataset.data[context.dataIndex];
          const percentage = (currentValue / total * 100).toFixed(1);
          return percentage > 0 ? `${context.label}: ${percentage}%` : '';
        },
      },
    },
  },
});

export default function PollStatsPage() {
  const { formId } = useParams();
  const [stats, setStats] = useState([]);
  const [responseHistory, setResponseHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('count');
  const [showPercentages, setShowPercentages] = useState(false);
  const socketRef = useRef(null);
  const answerColors = useRef({});

  useEffect(() => {
    const userId = localStorage.getItem('userId');
    if (!userId) {
      setError('User ID not found in local storage');
      setLoading(false);
      return;
    }

    const socketUrl = `ws://localhost:8086/ws?user_id=${userId}&form_id=${formId}`;
    socketRef.current = new WebSocket(socketUrl);

    socketRef.current.onopen = () => {
      console.log('WebSocket connected');
      setLoading(false);
    };

    socketRef.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const newData = { timestamp: new Date(), stats: data.stats };
        setResponseHistory((prev) => [...prev, newData]);
        setStats(data.stats);
      } catch (err) {
        console.error('Error parsing WebSocket message:', err);
        setError('Failed to parse server response');
      }
    };

    socketRef.current.onerror = (error) => {
      console.error('WebSocket error:', error);
      setError('WebSocket connection error');
      setLoading(false);
    };

    socketRef.current.onclose = () => {
      console.log('WebSocket disconnected');
    };

    return () => {
      if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
        socketRef.current.close();
      }
    };
  }, [formId]);

  const lastUpdated = responseHistory.length > 0 ? responseHistory[responseHistory.length - 1].timestamp : null;

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center h-screen">
        <div className="text-red-500 text-lg">{error}</div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-center mb-8">Статистика опроса</h1>
      <div className="flex justify-center space-x-4 mb-6">
        {['count', 'time', 'distribution'].map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-6 py-2 rounded-full font-medium transition-colors ${
              activeTab === tab ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            {tab === 'count' ? 'Количество ответов' : tab === 'time' ? 'Временная активность' : 'Распределение ответов'}
          </button>
        ))}
        {activeTab === 'distribution' && (
          <button
            onClick={() => setShowPercentages((prev) => !prev)}
            className={`px-6 py-2 rounded-full font-medium transition-colors ${
              showPercentages ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            {showPercentages ? 'Скрыть проценты' : 'Показать проценты'}
          </button>
        )}
      </div>
      {lastUpdated && (
        <p className="text-sm text-gray-500 text-center mb-6">
          Последнее обновление: {new Date(lastUpdated).toLocaleString()}
        </p>
      )}
      {stats.length === 0 ? (
        <div className="text-center py-10">
          <p className="text-gray-500">Данные статистики пока недоступны</p>
        </div>
      ) : (
        <div className="space-y-8">
          {stats.map((pollStat, index) => (
            <PollStatChart
              key={index}
              pollStat={pollStat}
              activeTab={activeTab}
              responseHistory={responseHistory}
              answerColors={answerColors.current}
              showPercentages={showPercentages}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function PollStatChart({ pollStat, activeTab, responseHistory, answerColors, showPercentages }) {
  const barChartRef = useRef(null);
  const lineChartRef = useRef(null);
  const pieChartRef = useRef(null);

  const handleDownloadChart = () => {
    let chartRef;
    let chartType;
    switch (activeTab) {
      case 'count':
        chartRef = barChartRef;
        chartType = 'bar';
        break;
      case 'time':
        chartRef = lineChartRef;
        chartType = 'line';
        break;
      case 'distribution':
        chartRef = pieChartRef;
        chartType = 'pie';
        break;
      default:
        return;
    }

    const chart = chartRef.current;
    if (chart) {
      const link = document.createElement('a');
      link.download = `${pollStat.question_title}_${chartType}.png`;
      link.href = chart.toBase64Image();
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } else {
      console.error('Chart instance not found');
    }
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow-lg">
      <h2 className="text-2xl font-semibold mb-4 text-gray-800">{pollStat.question_title}</h2>
      <div className="relative h-96 mb-4">
        {activeTab === 'count' ? (
          <Bar
            ref={barChartRef}
            data={getCountChartData(pollStat)}
            options={countChartOptions}
            plugins={[whiteBackgroundPlugin]}
          />
        ) : activeTab === 'time' ? (
          <Line
            ref={lineChartRef}
            data={getTimeChartData(pollStat, responseHistory, answerColors)}
            options={timeChartOptions}
            plugins={[whiteBackgroundPlugin]}
          />
        ) : (
          <Pie
            ref={pieChartRef}
            data={getPieChartData(pollStat)}
            options={getPieChartOptions(showPercentages)}
            plugins={[whiteBackgroundPlugin]}
          />
        )}
      </div>
      <button
        onClick={handleDownloadChart}
        className="flex items-center space-x-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors mb-4"
      >
        <FiDownload />
        <span>Скачать диаграмму</span>
      </button>
      <div>
        <h3 className="text-lg font-medium mb-2 text-gray-700">Детали ответов:</h3>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Ответ
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Количество
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {pollStat.answers.map((answer, idx) => (
                <tr key={idx} className={idx % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">{answer.title}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">{answer.count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}