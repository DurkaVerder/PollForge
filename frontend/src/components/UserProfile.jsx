export default function UserProfile() {
  const topics = [
    {
      title: "Technology",
      pollsCount: 8,
      description: "Programming languages, frameworks, and tech preferences",
      votes: "2.4k"
    },
    {
      title: "Work Environment",
      pollsCount: 5,
      description: "Remote work, office preferences, and productivity",
      votes: "1.2k"
    },
    {
      title: "UX/UI Design",
      pollsCount: 4,
      description: "Design trends, tools, and user experience patterns",
      votes: "950"
    },
    {
      title: "Software Development",
      pollsCount: 7,
      description: "Development practices, methodologies, and teams",
      votes: "1.8k"
    }
  ];

  return (
    <section className="bg-white rounded-lg shadow-md p-6 mb-8" id="profile">
      <div className="flex flex-col md:flex-row md:items-center mb-6 gap-6">
        <div className="flex-shrink-0">
          <div className="h-24 w-24 rounded-full overflow-hidden border-4 border-primary-100 shadow-md">
            <img
              src="https://images.unsplash.com/photo-1438761681033-6461ffad8d80?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3MzkyNDZ8MHwxfHNlYXJjaHw0fHx1c2VyfGVufDB8fHx8MTc0NjcxNTkzNXww&ixlib=rb-4.1.0&q=80&w=1080"
              alt="User profile"
              className="h-full w-full object-cover"
            />
          </div>
        </div>
        <div className="flex-1">
          <h2 className="text-2xl font-bold">John Doe</h2>
          <p className="text-gray-500">@johndoe</p>
          <p className="mt-2">
            Passionate about creating engaging polls and gathering insights.
          </p>
        </div>
        <div className="flex space-x-3 mt-4 md:mt-0">
          <button className="bg-primary-500 hover:bg-primary-600 text-white px-4 py-2 rounded-lg transition-colors duration-300">
            Edit Profile
          </button>
          <button className="border border-gray-300 hover:border-gray-400 px-4 py-2 rounded-lg transition-colors duration-300">
            Settings
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">24</div>
          <div className="text-gray-500">Polls Created</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">4.8k</div>
          <div className="text-gray-500">Total Votes</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-4 text-center transform hover:scale-105 transition-transform duration-300">
          <div className="text-3xl font-bold text-primary-600">256</div>
          <div className="text-gray-500">Poll Comments</div>
        </div>
      </div>

      <div className="text-xl font-bold mb-4">Poll Topics</div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {topics.map((topic, index) => (
          <button key={index} className="border rounded-lg p-4 hover:shadow-md transition-shadow duration-300 hover:bg-gray-50 cursor-pointer text-left flex flex-col">
            <div className="flex justify-between mb-2">
              <h4 className="font-semibold">{topic.title}</h4>
              <span className="text-sm text-gray-500">{topic.pollsCount} polls</span>
            </div>
            <p className="text-gray-500 text-sm mb-3">
              {topic.description}
            </p>
            <div className="mt-auto flex justify-between items-center">
              <span className="text-sm text-primary-600">{topic.votes} total votes</span>
              <span className="material-symbols-outlined text-primary-600">
                chevron_right
              </span>
            </div>
          </button>
        ))}
      </div>

      <div className="mt-8">
        <h3 className="text-xl font-bold mb-4">Topic Analytics</h3>
        <div className="bg-white rounded-lg shadow p-4 mb-6">
          <div className="h-[300px] w-full bg-gray-50 rounded-lg flex items-center justify-center">
            <svg className="w-full h-full" viewBox="0 0 600 300">
              <g transform="translate(50,20)">
                <rect x="0" y="30" width="80" height="200" fill="#818cf8" rx="4"></rect>
                <rect
                  x="100"
                  y="80"
                  width="80"
                  height="150"
                  fill="#a78bfa"
                  rx="4"
                ></rect>
                <rect
                  x="200"
                  y="130"
                  width="80"
                  height="100"
                  fill="#c084fc"
                  rx="4"
                ></rect>
                <rect
                  x="300"
                  y="180"
                  width="80"
                  height="50"
                  fill="#e879f9"
                  rx="4"
                ></rect>
                <rect
                  x="400"
                  y="110"
                  width="80"
                  height="120"
                  fill="#f472b6"
                  rx="4"
                ></rect>

                <text x="40" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Technology
                </text>
                <text x="140" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Work
                </text>
                <text x="240" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  UX/UI
                </text>
                <text x="340" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Dev
                </text>
                <text x="440" y="250" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Other
                </text>

                <text x="40" y="20" textAnchor="middle" fill="#4b5563">
                  8
                </text>
                <text x="140" y="70" textAnchor="middle" fill="#4b5563">
                  5
                </text>
                <text x="240" y="120" textAnchor="middle" fill="#4b5563">
                  4
                </text>
                <text x="340" y="170" textAnchor="middle" fill="#4b5563">
                  2
                </text>
                <text x="440" y="100" textAnchor="middle" fill="#4b5563">
                  5
                </text>

                <text
                  x="-40"
                  y="130"
                  textAnchor="middle"
                  transform="rotate(-90 -40 130)"
                  fill="#4b5563"
                  fontSize="14"
                >
                  Number of Polls
                </text>
                <text x="240" y="280" textAnchor="middle" fill="#4b5563" fontSize="14">
                  Topics
                </text>
                <text
                  x="240"
                  y="-10"
                  textAnchor="middle"
                  fill="#4b5563"
                  fontSize="16"
                  fontWeight="bold"
                >
                  Poll Distribution by Topic
                </text>
              </g>
            </svg>
          </div>
          <div className="flex justify-center mt-4">
            <button className="text-primary-600 hover:text-primary-700 font-medium hover:underline transition-all duration-300 flex items-center">
              View Detailed Analytics
              <span className="material-symbols-outlined ml-1">arrow_forward</span>
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}