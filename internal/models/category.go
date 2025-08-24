package models

type Category struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string    `json:"name" gorm:"type:text;not null"`
	Order          uint      `json:"order" gorm:"type:INTEGER;not null"`
	InternalUpdate bool      `json:"-" gorm:"type:boolean;default:false"`

	Products []Product `json:"products" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// Create triggers for create, update, delete category with recursion protection

/*
-- 1. Agregar columna auxiliar
ALTER TABLE categories ADD COLUMN IF NOT EXISTS internal_update BOOLEAN DEFAULT FALSE;

-- 2. Constraint para unicidad del orden
ALTER TABLE categories ADD CONSTRAINT unique_order UNIQUE ("order") DEFERRABLE INITIALLY DEFERRED;

-- 3. Función principal de actualización (evita recursión)
CREATE OR REPLACE FUNCTION categories_update() RETURNS TRIGGER AS $$
DECLARE
    cant INTEGER;
BEGIN
    IF NEW.internal_update THEN
        RETURN NEW;
    END IF;

    SELECT COUNT(*) INTO cant FROM categories;
    IF NEW."order" < 1 OR NEW."order" > cant THEN
        RAISE EXCEPTION 'new order out of range';
    END IF;

    UPDATE categories
    SET "order" = OLD."order", internal_update = TRUE
    WHERE "order" = NEW."order" AND id != NEW.id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4. Función para limpiar el flag
CREATE OR REPLACE FUNCTION categories_clear_flag() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.internal_update THEN
        UPDATE categories SET internal_update = FALSE WHERE id = NEW.id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 5. Crear triggers de actualización con protección
DROP TRIGGER IF EXISTS categories_update_trigger ON categories;
DROP TRIGGER IF EXISTS categories_clear_flag_trigger ON categories;

CREATE TRIGGER categories_update_trigger
BEFORE UPDATE ON categories
FOR EACH ROW
WHEN (OLD."order" IS DISTINCT FROM NEW."order")
EXECUTE FUNCTION categories_update();

CREATE TRIGGER categories_clear_flag_trigger
AFTER UPDATE ON categories
FOR EACH ROW
WHEN (NEW.internal_update IS TRUE)
EXECUTE FUNCTION categories_clear_flag();

-- 6. Trigger para inserción automática de orden
CREATE OR REPLACE FUNCTION categories_create() RETURNS TRIGGER AS $$
DECLARE
    cant INTEGER;
BEGIN
    SELECT COUNT(*) INTO cant FROM categories;
    NEW."order" = cant + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER categories_create
BEFORE INSERT ON categories
FOR EACH ROW
EXECUTE FUNCTION categories_create();

-- 7. Trigger para actualización del orden al eliminar
CREATE OR REPLACE FUNCTION categories_delete() RETURNS TRIGGER AS $$
BEGIN
    UPDATE categories SET "order" = "order" - 1, internal_update = TRUE WHERE "order" > OLD."order";
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER categories_delete
AFTER DELETE ON categories
FOR EACH ROW
EXECUTE FUNCTION categories_delete();

*/
