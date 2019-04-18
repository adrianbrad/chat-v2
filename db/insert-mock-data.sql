INSERT INTO "rooms" (room_id)
    VALUES ('room_a');

INSERT INTO "users" (user_id, nickname)
    VALUES ('user_a', 'someone'), ('user_b', 'random');


INSERT INTO "permissions" (permission_id)
    VALUES 
        ('send_message'),
        ('run'),
        ('dance'),
        ('drink');

INSERT INTO "users_permissions" (permission_id, user_id)
    VALUES 
        ('send_message', 'user_a'),
        ('run', 'user_a');