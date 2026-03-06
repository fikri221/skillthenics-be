-- name: GetOrderWithItems :many
-- Demonstrating JOIN to fetch Order, Items, and Product details
SELECT 
    o.id as order_id, 
    o.customer_name, 
    o.total_amount, 
    o.created_at as order_date,
    oi.id as item_id,
    oi.product_id,
    oi.quantity,
    oi.price as item_price,
    p.name as product_name
FROM ms_orders o
LEFT JOIN ms_order_items oi ON o.id = oi.order_id
LEFT JOIN ms_products p ON oi.product_id = p.id
WHERE o.id = ? AND o.rec_status = 'A';

-- name: CreateOrder :execresult
INSERT INTO ms_orders (id, customer_name, total_amount, rec_status) 
VALUES (?, ?, ?, 'A');

-- name: CreateOrderItem :execresult
INSERT INTO ms_order_items (order_id, product_id, quantity, price) 
VALUES (?, ?, ?, ?);

-- name: ListOrders :many
SELECT * FROM ms_orders WHERE rec_status = 'A' ORDER BY created_at DESC;
