# Go (Gin) + PostgreSQL CRUD API

Proyek RESTful API menggunakan Go (Gin Engine) dan PostgreSQL, dilengkapi implementasi Stored Procedure.

## Persyaratan
- Docker & Docker Compose

## Cara Menjalankan
```bash
docker compose up -d --build


API Base URL: http://localhost:8081

PostgreSQL Host Port: 5433 (User: postgres, Pass: postgres123, DB: inventory_db)

Endpoints
GET /api/items - List data item

GET /api/items/:id - Ambil detail item

POST /api/items - Tambah item

PUT /api/items/:id - Update data item

DELETE /api/items/:id - Hapus data item



Salin script Stored Procedure ke dalam file `procedure.sql` di dalam folder ini agar reviewer bisa langsung membaca kodenya[cite: 1]:
```sql
CREATE TABLE IF NOT EXISTS stock_logs (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    change_type VARCHAR(10) NOT NULL,
    quantity INT NOT NULL,
    previous_stock INT NOT NULL,
    current_stock INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE PROCEDURE sp_deduct_item_stock(
    p_item_id INT,
    p_quantity INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_current_stock INT;
BEGIN
    IF p_quantity <= 0 THEN
        RAISE EXCEPTION 'Jumlah kuantitas harus lebih dari 0';
    END IF;

    SELECT stock INTO v_current_stock
    FROM items
    WHERE id = p_item_id
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Item dengan ID % tidak ditemukan', p_item_id;
    END IF;

    IF v_current_stock < p_quantity THEN
        RAISE EXCEPTION 'Stok tidak mencukupi. Sisa stok: %, Diminta: %', v_current_stock, p_quantity;
    END IF;

    UPDATE items
    SET stock = stock - p_quantity,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_item_id;

    INSERT INTO stock_logs (item_id, change_type, quantity, previous_stock, current_stock)
    VALUES (p_item_id, 'DEDUCT', p_quantity, v_current_stock, v_current_stock - p_quantity);
END;
$$;