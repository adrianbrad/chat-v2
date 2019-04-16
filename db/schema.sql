DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS rooms;


CREATE TABLE users
(
  user_id             TEXT NOT NULL CONSTRAINT users_pkey PRIMARY KEY,
  nickname            TEXT NOT NULL,
  created_at          TIMESTAMP DEFAULT now() NOT NULL,
  updated_at          TIMESTAMP
);

CREATE TABLE rooms
(
  room_id           TEXT NOT NULL CONSTRAINT rooms_pkey PRIMARY KEY
);

CREATE TABLE messages
(
  message_id          BIGSERIAL NOT NULL CONSTRAINT messages_pkey PRIMARY KEY,
  content             TEXT NOT NULL,
  created_at          TIMESTAMP DEFAULT now() NOT NULL,
  room_id             TEXT NOT NULL CONSTRAINT FK_Messages_Rooms REFERENCES rooms,
  user_id             TEXT NOT NULL CONSTRAINT FK_Messages_Users REFERENCES users
);
