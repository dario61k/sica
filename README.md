# SICA

SICA es un backend para la gestión de productos y categorías. Pensado para servir como base en proyectos que requieran un backend sencillo para catálogo de productos. Simple y facil de extender.

**Indíce**

- [[#Características|Características]]
- [[#Requisitos|Requisitos]]
- [[#Instalación|Instalación]]
- [[#Variables de entorno|Variables de entorno]]
- [[#API Endpoints|API Endpoints]]
	- [[#API Endpoints#Públicos|Públicos]]
	- [[#API Endpoints#Autenticación|Autenticación]]
	- [[#API Endpoints#Gestión (requiere autenticación)|Gestión (requiere autenticación)]]
- [[#Base de datos (IMPORTANTE)|Base de datos (IMPORTANTE)]]

---
## Características

- Autenticación basada en JWT (un solo usuario administrador)
- Gestión de productos y categorías
- Almacenamiento de imágenes en Cloudinary
- Contraseñas protegidas con bcrypt
- Estructura modular en Go

---
## Requisitos

- Go 1.23+ (versión probada)
- PostgreSQL
- Cuenta de Cloudinary

---
## Instalación

```bash
git clone https://github.com/dario61k/sica.git
cd sica
go mod tidy
go run ./cmd/api/main.go
```

---
## Variables de entorno

```env
PORT=3000
DB_URL=postgresql://user:password@localhost/dbname
JWT_SECRET=your-secret-key
ADMIN_PASSWORD=your-admin-password
CLIENT_DOMAIN=http://localhost:3000
C_CLOUD_NAME=your-cloudinary-name
C_API_KEY=your-cloudinary-key
C_API_SECRET=your-cloudinary-secret
```

---
## API Endpoints

### Públicos
- `GET /api/get-all` - Obtiene todas las categorías y productos visibles

### Autenticación
- `POST /api/login` - Login del administrador (solo password)
- `POST /api/refresh-token` - Refresca el access token
- `GET /api/auth` - Verifica si el token es válido

### Gestión (requiere autenticación)
- `GET /api/category` - Lista todas las categorías
- `POST /api/category` - Crea una nueva categoría
- `PUT /api/category/:id` - Actualiza una categoría
- `DELETE /api/category/:id` - Elimina una categoría

- `GET /api/product` - Lista todos los productos
- `POST /api/product` - Crea un nuevo producto (con imagen)
- `PUT /api/product/:id` - Actualiza un producto
- `DELETE /api/product/:id` - Elimina un producto

---
## Base de datos (IMPORTANTE)

La aplicación creará automáticamente las tablas necesarias. **Después de la primera ejecución, ejecuta este SQL en tu base de datos** para habilitar configuraciones necesarias en la tabla  `categories` :

```sql
-- 1. Agregar columna auxiliar
ALTER TABLE categories 
ADD COLUMN IF NOT EXISTS internal_update BOOLEAN DEFAULT FALSE;

-- 2. Constraint para unicidad del orden
ALTER TABLE categories 
ADD CONSTRAINT unique_order UNIQUE ("order") DEFERRABLE INITIALLY DEFERRED;

-- 3. Función principal de actualización (evita recursión)
CREATE OR REPLACE FUNCTION categories_update() 
RETURNS TRIGGER AS $$
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
CREATE OR REPLACE FUNCTION categories_clear_flag() 
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.internal_update THEN
        UPDATE categories 
        SET internal_update = FALSE 
        WHERE id = NEW.id;
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
CREATE OR REPLACE FUNCTION categories_create() 
RETURNS TRIGGER AS $$
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
CREATE OR REPLACE FUNCTION categories_delete() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE categories 
    SET "order" = "order" - 1, internal_update = TRUE 
    WHERE "order" > OLD."order";

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER categories_delete
AFTER DELETE ON categories
FOR EACH ROW
EXECUTE FUNCTION categories_delete();

```
