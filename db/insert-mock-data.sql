INSERT INTO "rooms" (room_id)
    VALUES ('room_a');

INSERT INTO "users" (user_id, nickname)
    VALUES ('user_a', 'someone');


INSERT INTO "permissions" (permission_id)
    VALUES 
        ('talk'),
        ('run'),
        ('dance'),
        ('drink');

INSERT INTO "users_permissions" (permission_id, user_id)
    VALUES 
        ('talk', 'user_a'),
        ('run', 'user_a');