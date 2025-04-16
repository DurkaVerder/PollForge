

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE forms(
    id SERIAL PRIMARY KEY,
    creator_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    private_key BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (creator_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE questions(
    id SERIAL PRIMARY KEY,
    form_id INT,
    number_order INT,
    title VARCHAR(255) NOT NULL,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE
);

CREATE TABLE answers(
    id SERIAL PRIMARY KEY,
    question_id INT,
    title VARCHAR(255) NOT NULL,
    number_order INT,
    count INT DEFAULT 0,
    FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE
);

CREATE TABLE comments(
    id SERIAL PRIMARY KEY,
    form_id INT,
    title VARCHAR(255) NOT NULL,
    count INT DEFAULT 0,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE
);

CREATE TABLE likes(
    id SERIAL PRIMARY KEY,
    form_id INT,
    user_id INT NOT NULL,
    count INT DEFAULT 0,
    FOREIGN KEY (form_id) REFERENCES forms (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);