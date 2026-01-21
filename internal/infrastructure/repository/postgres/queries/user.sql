--SQL:Create
INSERT INTO users (name) VALUES ($1)
RETURNING *;
--end

--SQL:Get
SELECT * FROM users WHERE id = $1;
--end

--SQL:GetAll
SELECT * FROM users;
--end

--SQL:Update
UPDATE users SET name = $2 WHERE id = $1
RETURNING *;
--end

--SQL:Delete
DELETE FROM users WHERE id = $1;
--end