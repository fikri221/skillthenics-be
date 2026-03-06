-- +goose Up
-- ms_orders: Order header table
CREATE TABLE IF NOT EXISTS ms_orders (
    id VARCHAR(27) PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    total_amount DECIMAL(15, 2) NOT NULL,
    rec_status VARCHAR(1) NOT NULL DEFAULT 'A',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ms_order_items: Order detail table linked to products
CREATE TABLE IF NOT EXISTS ms_order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id VARCHAR(27) NOT NULL,
    product_id VARCHAR(27) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL, -- Historical price at the time of order
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES ms_orders(id),
    FOREIGN KEY (product_id) REFERENCES ms_products(id)
);

-- +goose Down
DROP TABLE IF EXISTS ms_order_items;
DROP TABLE IF EXISTS ms_orders;
