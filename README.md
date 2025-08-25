# SICA

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue)](https://golang.org/)  
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-green)](https://www.postgresql.org/)  
![License](https://img.shields.io/badge/license-MIT-lightgrey.svg)

**SICA** is a lightweight backend for managing products and categories. It is designed to be a simple, extensible foundation for projects that require a product catalog backend.

---

## ✨ Features

- JWT-based authentication (single admin user)
    
- Product and category management
    
- Cloudinary image storage
    
- Passwords secured with bcrypt
    
- Modular structure in Go
    

---

## 📦 Requirements

- **Go** 1.23+ (tested version)
    
- **PostgreSQL**
    
- **Cloudinary account**
    

---

## 🚀 Installation

```bash
git clone https://github.com/dario61k/sica.git
cd sica
go mod tidy
go run ./cmd/api/main.go
```

---

## ⚙️ Environment Variables

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

## 🔌 API Endpoints

### Public

- `GET /api/get-all` → Fetch all visible categories and products
    

### Authentication

- `POST /api/login` → Admin login (password only)
    
- `POST /api/refresh-token` → Refresh access token
    
- `GET /api/auth` → Validate token
    

### Management (requires authentication)

- **Categories**
    
    - `GET /api/category` → List all categories
        
    - `POST /api/category` → Create a new category
        
    - `PUT /api/category/:id` → Update a category
        
    - `DELETE /api/category/:id` → Delete a category
        
- **Products**
    
    - `GET /api/product` → List all products
        
    - `POST /api/product` → Create a new product (with image)
        
    - `PUT /api/product/:id` → Update a product
        
    - `DELETE /api/product/:id` → Delete a product
        

---

## 🗄️ Database Setup (IMPORTANT)

The application will create the required tables automatically.  
**After the first run, execute the following SQL in your database** to enable necessary constraints and triggers for the `categories` table:

```sql
-- Add auxiliary column
ALTER TABLE categories 
ADD COLUMN IF NOT EXISTS internal_update BOOLEAN DEFAULT FALSE;

-- Unique constraint for order
ALTER TABLE categories 
ADD CONSTRAINT unique_order UNIQUE ("order") DEFERRABLE INITIALLY DEFERRED;

-- Update function (avoids recursion)
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

-- Clear flag function
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

-- Create triggers
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

-- Auto-order on insert
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

-- Update order on delete
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

---

## 🛠️ Contributing

Contributions, issues, and feature requests are welcome!  
Feel free to fork the repository and submit a PR.

---

## 📜 License

This project is licensed under the MIT License.
