INSERT INTO books (id, title, author) VALUES (1, 'Herman Melville', 'Moby Dick');
INSERT INTO books (id, title, author) VALUES (2, 'Leo Tolstoy', 'War and Peace');
INSERT INTO books (id, title, author) VALUES (3, 'William Shakespeare', 'Hamlet');
INSERT INTO books (id, title, author) VALUES (4, 'Homer', 'The Odyssey');
INSERT INTO books (id, title, author) VALUES (5, 'Mark Twain', 'The Adventures of Huckleberry Finn');

SELECT setval('books_id_seq', (SELECT MAX(id) FROM books));