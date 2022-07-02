CREATE TABLE IF NOT EXISTS orders (
    uid TEXT UNIQUE NOT NULL PRIMARY KEY,
    json_order TEXT NOT NULL
)

    CREATE PROCEDURE getOrder(orderUID text, order text)
    LANGUAGE SQL
    AS $$
        INSERT INTO orders (uid, json_order) VALUES ($1, $2);
    $$;    
