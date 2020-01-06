INSERT INTO musics (id, artist) VALUES (1, 'Maria Carey');
INSERT INTO musics (id, artist) VALUES (2, 'Agnes Monica');

SELECT setval('musics_id_seq', (SELECT MAX(id) FROM musics));