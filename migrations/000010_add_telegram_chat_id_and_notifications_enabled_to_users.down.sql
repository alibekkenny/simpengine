ALTER TABLE users
DROP COLUMN IF EXISTS telegram_chat_id
DROP COLUMN IF EXISTS notifications_enabled;