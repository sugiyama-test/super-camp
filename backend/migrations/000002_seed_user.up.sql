INSERT INTO users (id, name, email) VALUES (1, 'テストユーザー', 'test@example.com')
ON CONFLICT (id) DO NOTHING;
