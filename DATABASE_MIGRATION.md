# Database Migration to GORM

This document describes the migration from raw SQL to GORM ORM for the SpiceRoute project.

## Changes Made

### 1. Added GORM Models (`pkg/models/models.go`)

Created comprehensive GORM models that correspond to the protobuf messages:

- **User**: Base user entity with UUID primary key
- **Preference**: User dietary preferences with foreign key to User
- **Recipe**: Recipe entity with all necessary fields and arrays
- **Feedback**: User feedback on dishes with foreign keys to User and Recipe

### 2. Database Configuration (`pkg/database/database.go`)

- Centralized database connection management
- Automatic migration functionality
- Proper connection cleanup

### 3. Updated Services

All Go services have been updated to use GORM instead of raw SQL:

#### Feedback Service (`services/feedback/main.go`)

- Replaced raw SQL with GORM transactions
- Added proper time parsing for `cooked_at` field
- Implemented upsert functionality using GORM's `FirstOrCreate`

#### Profile Service (`services/profile/main.go`)

- Replaced raw SQL with GORM operations
- Added `GetPreference` method implementation
- Used GORM's `FirstOrCreate` for upsert operations

#### Recipe Service (`services/recipes/main.go`)

- Replaced raw SQL with GORM operations
- Added filtering capabilities for cuisine and spicy preferences
- Proper conversion between GORM models and protobuf messages

### 4. Dependencies

Added to `go.mod`:

- `gorm.io/gorm` - Core GORM library
- `gorm.io/driver/postgres` - PostgreSQL driver for GORM

## Database Schema

The GORM models will create the following tables:

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Preferences Table

```sql
CREATE TABLE preferences (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    cuisines TEXT[],
    allergies TEXT[],
    budget_week DOUBLE PRECISION,
    spicy BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Recipes Table

```sql
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    cuisine VARCHAR,
    prep_minutes INTEGER,
    calories INTEGER,
    ingredients TEXT[],
    cost DOUBLE PRECISION,
    shelf_life_days INTEGER,
    tags TEXT[],
    nutrition TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Feedback Table

```sql
CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    dish_id UUID NOT NULL,
    rating INTEGER,
    skipped BOOLEAN,
    substituted_with VARCHAR,
    comment TEXT,
    cooked_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (dish_id) REFERENCES recipes(id)
);
```

## Running Migrations

### Option 1: Automatic Migration (Recommended)

Each service will automatically run migrations on startup. Just ensure the `DB_DSN` environment variable is set.

### Option 2: Manual Migration

Run the standalone migration script:

```bash
export DB_DSN="postgres://username:password@localhost:5432/spiceroute?sslmode=disable"
go run scripts/migrate.go
```

## Environment Variables

Required environment variable:

- `DB_DSN`: PostgreSQL connection string

Example:

```
DB_DSN=postgres://username:password@localhost:5432/spiceroute?sslmode=disable
```

## Benefits of GORM Migration

1. **Type Safety**: Compile-time checking of database operations
2. **Automatic Migrations**: Schema changes are handled automatically
3. **Relationships**: Easy handling of foreign key relationships
4. **Query Building**: Fluent API for building complex queries
5. **Hooks**: Lifecycle hooks for validation and business logic
6. **Soft Deletes**: Built-in support for soft deletes
7. **Transactions**: Simplified transaction management

## Backward Compatibility

The API interfaces remain unchanged - all protobuf messages and gRPC service definitions are preserved. The migration is internal to the services.

## Testing

To test the migration:

1. Set up a PostgreSQL database
2. Set the `DB_DSN` environment variable
3. Run the migration script or start any service
4. Verify tables are created correctly
5. Test API endpoints to ensure functionality is preserved
