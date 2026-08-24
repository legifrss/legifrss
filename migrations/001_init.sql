CREATE TABLE jorf (id TEXT PRIMARY KEY, date TIMESTAMP NOT NULL, content JSONB, jorf_content JSONB);
CREATE TABLE jorf_text (id TEXT PRIMARY KEY, date TIMESTAMP NOT NULL, content JSONB, nature TEXT, author TEXT);
