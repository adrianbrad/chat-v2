CREATE TABLE permissions
(
  permission_id       TEXT NOT NULL CONSTRAINT permissions_pkey PRIMARY KEY
);

CREATE TABLE users
(
  user_id             TEXT NOT NULL CONSTRAINT users_pkey PRIMARY KEY,
  nickname            TEXT NOT NULL,
  created_at          TIMESTAMP DEFAULT now() NOT NULL,
  updated_at          TIMESTAMP
);

CREATE TABLE users_permissions
(
  user_id             TEXT NOT NULL CONSTRAINT FK_Users_Permissions_Users REFERENCES users,
  permission_id       TEXT NOT NULL CONSTRAINT FK_Users_Permissions_Permissions REFERENCES permissions,
  CONSTRAINT users_permissions_pkey PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE rooms
(
  room_id             TEXT NOT NULL CONSTRAINT rooms_pkey PRIMARY KEY
);

CREATE TABLE messages
(
  message_id          BIGSERIAL NOT NULL CONSTRAINT messages_pkey PRIMARY KEY,
  content             TEXT NOT NULL,
  created_at          TIMESTAMP DEFAULT now() NOT NULL,
  room_id             TEXT NOT NULL CONSTRAINT FK_Messages_Rooms REFERENCES rooms,
  user_id             TEXT NOT NULL CONSTRAINT FK_Messages_Users REFERENCES users
);
