-- +goose Up
-- Change id from INT to VARCHAR(27) to support KSUID
ALTER TABLE ms_products MODIFY id VARCHAR(27);

-- +goose Down
-- This is destructive if data exists, but necessary for rollback logic
ALTER TABLE ms_products MODIFY id INT AUTO_INCREMENT;
