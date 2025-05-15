export default function PollCard({ poll }) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 transform hover:shadow-lg transition-all duration-300">
      <div>
        <div className="flex items-center mb-3">
          {poll.categories.map((category, index) => (
            <span 
              key={index}
              className={`${index === 0 ? 'mr-2' : ''} ${
                index % 2 === 0 ? 'bg-blue-100 text-blue-800' : 'bg-green-100 text-green-800'
              } text-xs font-medium px-2.5 py-0.5 rounded-full`}
            >
              {category}
            </span>
          ))}
        </div>
      </div>

      <div className="flex justify-between items-start mb-4">
        <div className="flex items-center">
          <img
            src={poll.user.avatar}
            alt="Аватар пользователя"
            className="h-10 w-10 rounded-full mr-3"
          />
          <div>
            <h3 className="font-semibold">{poll.user.name}</h3>
            <p className="text-sm text-gray-500">{poll.user.time}</p>
          </div>
        </div>
        <button className="text-gray-400 hover:text-gray-600">
          <span className="material-symbols-outlined">more_vert</span>
        </button>
      </div>

      {poll.questions.map((question, qIndex) => (
        <div key={qIndex} className={qIndex < poll.questions.length - 1 ? "border-b pb-4 mb-4" : "mb-4"}>
          <h4 className="text-xl font-semibold mb-3">{question.title}</h4>
          <div className="space-y-3 mb-6">
            {question.options.map((option, oIndex) => (
              <div key={oIndex} className="flex items-center">
                <input
                  type="radio"
                  id={`poll${poll.id}_q${qIndex}_option${oIndex}`}
                  name={`poll${poll.id}_q${qIndex}`}
                  className="h-4 w-4 text-primary-600"
                />
                <label htmlFor={`poll${poll.id}_q${qIndex}_option${oIndex}`} className="ml-2 block w-full">
                  <div className="flex justify-between">
                    <span>{option.label}</span>
                    <span className="text-sm text-gray-500">{option.percentage}%</span>
                  </div>
                  <div className="mt-1 h-2 w-full bg-gray-200 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-primary-500 rounded-full"
                      style={{width: `${option.percentage}%`}}
                    ></div>
                  </div>
                </label>
              </div>
            ))}
          </div>
        </div>
      ))}

      <div className="flex justify-between text-sm text-gray-500">
        <span>{poll.votes} голосов</span>
        <span>Заканчивается через {poll.endsIn}</span>
      </div>
      <div className="mt-4 flex justify-between">
        <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors">
          <span className="material-symbols-outlined mr-1">comment</span>
          {poll.comments} комментариев
        </button>
        <button className="flex items-center text-primary-600 hover:text-primary-700 transition-colors">
          <span className="material-symbols-outlined mr-1">share</span>
          Поделиться
        </button>
      </div>
    </div>
  );
}
