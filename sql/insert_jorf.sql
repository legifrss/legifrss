INSERT INTO jorf (id, date, content) VALUES ($1, $2, $3)
ON CONFLICT(id) DO UPDATE
SET
  date = $2,
  content = $3
