CREATE TYPE item_type_enum AS ENUM (
    'MATERIAL',
    'MANUFACTURED'
    );

CREATE TABLE IF NOT EXISTS items
(
    id                  UUID PRIMARY KEY,
    name                VARCHAR(255)             NOT NULL,
    description         TEXT,
    sku                 VARCHAR(100) UNIQUE,
    type                item_type_enum           NOT NULL,
    is_active           BOOLEAN                  NOT NULL DEFAULT true,
    can_be_sold         BOOLEAN                  NOT NULL DEFAULT false,
    price_sale_in_cents BIGINT                   NOT NULL DEFAULT 0,
    stock_quantity      NUMERIC(10, 3)           NOT NULL DEFAULT 0,
    unit_of_measure     VARCHAR(20)              NOT NULL,
    minimum_stock_level INT                      NOT NULL DEFAULT 0,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX IF NOT EXISTS idx_items_sku ON items (sku);
CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items (deleted_at);

CREATE TABLE IF NOT EXISTS recipes
(
    id         UUID PRIMARY KEY,
    product_id UUID NOT NULL UNIQUE,
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_lines
(
    id        UUID PRIMARY KEY,
    recipe_id UUID           NOT NULL,
    item_id   UUID           NOT NULL,
    quantity  NUMERIC(10, 3) NOT NULL,
    CONSTRAINT fk_recipe FOREIGN KEY (recipe_id) REFERENCES recipes (id) ON DELETE CASCADE,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE RESTRICT
);

CREATE TRIGGER set_timestamp_items
    BEFORE UPDATE
    ON items
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();