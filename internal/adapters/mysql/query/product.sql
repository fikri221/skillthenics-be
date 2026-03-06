-- name: ListProducts :many
SELECT * FROM ms_products 
WHERE rec_status = 'A' 
AND (name LIKE CONCAT('%', sqlc.arg(search), '%') OR sqlc.arg(search) = '')
LIMIT ? OFFSET ?;

-- name: GetProductByID :one
SELECT * FROM ms_products WHERE id = ? AND rec_status = 'A';

-- name: CreateProduct :execresult
INSERT INTO ms_products (id, name, price, rec_status) VALUES (?, ?, ?, 'A');

-- name: UpdateProduct :execresult
UPDATE ms_products SET name = ?, price = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND rec_status = 'A';

-- name: DeleteProduct :execresult
UPDATE ms_products SET rec_status = 'D', updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: CountProducts :one
SELECT COUNT(*) FROM ms_products 
WHERE rec_status = 'A' 
AND (name LIKE CONCAT('%', sqlc.arg(search), '%') OR sqlc.arg(search) = '');

-- name: CheckProductNameExists :one
SELECT EXISTS(
    SELECT 1 FROM ms_products 
    WHERE name = ? AND rec_status = 'A'
);

-- name: CheckProductNameExistsForOther :one
SELECT EXISTS(
    SELECT 1 FROM ms_products 
    WHERE name = ? AND id != ? AND rec_status = 'A'
);