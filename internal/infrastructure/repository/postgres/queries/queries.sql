--SQL:Create
INSERT INTO user (id, name) VALUES ($1, $2);
--end

--SQL:Get
SELECT * FROM user WHERE id = $1;
--end