

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

CREATE INDEX idx_users_id_name ON users (id, name);

CREATE INDEX idx_forms_id_creator_id_link ON forms (id, creator_id, link);

CREATE INDEX idx_forms_link ON forms (link);

CREATE INDEX idx_questions_id_form_id ON questions (id, form_id);

CREATE INDEX idx_answers_id_question_id ON answers (id, question_id);

CREATE INDEX idx_comments_id_user_id_form_id ON comments (id, user_id, form_id);

CREATE INDEX idx_likes_id_form_id ON likes (id, form_id);

CREATE INDEX idx_answered_polls_id_user_id_form_id ON answered_polls (id, user_id, form_id);

