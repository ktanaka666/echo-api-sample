use testdb;

INSERT INTO users (`name`, `email`) VALUES
  ('テストユーザ', 'test@example.com'),
  ('テストユーザ２', 'test2@example.com');

INSERT INTO tasks (`name`, `user_id`, `category_id`, `status_id`, `end_date`, `publish_flg`) VALUES
  ('タスク１', 1, 1, 1, '2019-12-31 00:00:00', 1),
  ('タスク２(期限なし)', 2, 2, 1, NULL, 1),
  ('タスク３(PoC)', 1, 3, 0, '2019-12-31 00:00:00', 0);

INSERT INTO categories (`name`) VALUES
  ('期限ありタスク'), ('期限なしタスク'), ('PoCタスク');