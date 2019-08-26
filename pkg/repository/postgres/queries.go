package repository

const CreateChatMessages = `
CREATE TABLE IF NOT EXISTS chatmessage (
    id SERIAL,
    chatroom_id INT,
    sender_id INT,
    message VARCHAR,
    PRIMARY KEY (id),
    FOREIGN KEY (chatroom_id) REFERENCES chatroom (id),
    FOREIGN KEY (sender_id) REFERENCES users (id)
)`

const CreateChatRooms = `
CREATE TABLE IF NOT EXISTS chatroom (
    id SERIAL,
    name VARCHAR UNIQUE,
    owner_id INT,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES users (id)
)`

const CreateChatRoomUsers = `
CREATE TABLE IF NOT EXISTS chatroomusers (
    chatroom_id INT,
    user_id INTEGER,
    PRIMARY KEY (chatroom_id, user_id),
    FOREIGN KEY (chatroom_id) REFERENCES chatroom (id),
	FOREIGN KEY (user_id) REFERENCES users (id)
)`

const CreateUsers = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    name VARCHAR,
    PRIMARY KEY (id)
)`
