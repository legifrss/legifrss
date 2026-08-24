INSERT INTO jorf_text (id, date, content, nature, author) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT(id) DO UPDATE
SET
  date = $2,
  content = $3,
  nature = $4,
  author = $5
