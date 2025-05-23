

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    role      VARCHAR(20) DEFAULT 'user',
    is_banned BOOLEAN DEFAULT FALSE,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    bio VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE password_resets (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	token VARCHAR(255) NOT NULL UNIQUE,
	expires_at TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE forms(
    id SERIAL PRIMARY KEY,
    theme_id INT NOT NULL,  
    creator_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    link VARCHAR(255) NOT NULL,
    private_key BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE questions(
    id SERIAL PRIMARY KEY,
    form_id INT,
    creator_id INT NOT NULL,
    number_order INT,
    required BOOLEAN DEFAULT FALSE,
    title VARCHAR(255) NOT NULL,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE
);

CREATE TABLE answers(
    id SERIAL PRIMARY KEY,
    question_id INT,
    creator_id INT NOT NULL,
    form_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    number_order INT,
    count BIGINT DEFAULT 0,
    FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE
);

CREATE TABLE answers_chosen (
    answer_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (answer_id, user_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (answer_id) REFERENCES answers (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE comments(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    form_id INT NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP WITH TIME ZONE NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE
);

CREATE TABLE likes(
    id SERIAL PRIMARY KEY,
    form_id INT,
    count INT DEFAULT 0,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE
);

CREATE TABLE likes_forms (
    id SERIAL PRIMARY KEY,
    form_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE answered_polls (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    form_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE
);

CREATE TABLE themes(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (role, name, email, password)
VALUES (
  'admin', 
  'superadmin', 
  'admin@example.com', 
  '$2a$10$EixZaYVK1fsbw1ZfbX3OXe.Pu6q2/o6HqZV0K6L/BCGgDCbXJbad.'
)

CREATE INDEX idx_users_id_name ON users (id, name);

CREATE INDEX idx_forms_id_creator_id_link ON forms (id, creator_id, link);

CREATE INDEX idx_forms_link ON forms (link);

CREATE INDEX idx_questions_id_form_id ON questions (id, form_id);

CREATE INDEX idx_answers_id_question_id ON answers (id, question_id);

CREATE INDEX idx_comments_id_user_id_form_id ON comments (id, user_id, form_id);

CREATE INDEX idx_likes_id_form_id ON likes (id, form_id);

CREATE INDEX idx_answered_polls_id_user_id_form_id ON answered_polls (id, user_id, form_id);


INSERT INTO themes (name, description) VALUES
('Образ жизни', 'Хобби, путешествия, семья и повседневные привычки'),
('Здоровье и благополучие', 'Физическое и ментальное здоровье, спорт и питание'),
('Технологии и инновации', 'Гаджеты, интернет и их роль в жизни'),
('Окружающая среда', 'Экология и забота о природе'),
('Образование и карьера', 'Обучение, работа и профессиональный рост'),
('Культура и искусство', 'Литература, музыка, кино и творчество'),
('Политика и общество', 'Социальные вопросы и права человека'),
('Экономика и финансы', 'Личные деньги и экономические тренды'),
('Спорт и развлечения', 'Спортивные события и популярная культура'),
('Наука и исследования', 'Открытия, космос и новые технологии'),
('Путешествия и туризм', 'Приключения и культурный обмен'),
('Еда и кулинария', 'Рецепты, гастрономия и традиции'),
('Мода и стиль', 'Тенденции в одежде и самовыражении'),
('Семья и отношения', 'Семейная жизнь, дружба и воспитание'),
('Дом и интерьер', 'Идеи для уюта и декора дома'),
('Финансовая грамотность', 'Управление финансами и бюджет'),
('Саморазвитие', 'Личностный рост и мотивация'),
('Игры и досуг', 'Видеоигры, настольные игры и отдых'),
('История и наследие', 'Прошлое и его влияние на сегодня'),
('Будущее и инновации', 'Прогнозы о технологиях и обществе'),
('Социальные сети', 'Влияние соцсетей на общение и жизнь'),
('Этика и мораль', 'Вопросы морали и ценностей в обществе'),
('Глобализация', 'Взаимосвязь культур и экономик'),
('Медиа и новости', 'Роль СМИ в формировании мнений'),
('Волонтерство', 'Помощь обществу и благотворительность'),
('Безопасность', 'Личная и общественная безопасность'),
('Транспорт и мобильность', 'Транспортные системы и их развитие'),
('Искусственный интеллект', 'Применение ИИ в жизни и работе'),
('Энергия и ресурсы', 'Альтернативная энергия и устойчивое развитие'),
('Литература и чтение', 'Книги и их влияние на мышление'),
('Музыка и концерты', 'Музыкальные жанры и живые выступления'),
('Кино и сериалы', 'Тренды в кинематографе и стриминге'),
('Фотография', 'Искусство фотографии и визуальный контент'),
('Танцы и движение', 'Танцевальные стили и физическая активность'),
('Психология', 'Понимание поведения и эмоций'),
('Религия и духовность', 'Духовные практики и верования'),
('Архитектура', 'Дизайн зданий и городских пространств'),
('Садоводство', 'Выращивание растений и озеленение'),
('Животные и природа', 'Забота о животных и их роль в жизни'),
('Космос и астрономия', 'Исследование звезд и вселенной');