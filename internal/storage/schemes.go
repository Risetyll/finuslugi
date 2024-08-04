package storage

const (
	TableSchema = `CREATE TABLE IF NOT EXISTS materials (
    	uuid SERIAL PRIMARY KEY,
    	material_type VARCHAR(50) CHECK (material_type IN ('статья', 'видеоролик', 'презентация')),
    	publication_status VARCHAR(50) CHECK (publication_status IN ('архивный', 'активный')),
    	title VARCHAR(255) NOT NULL,
    	content TEXT,
    	creation_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    	modification_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`
	MaterialIndexSchema  = `CREATE INDEX IF NOT EXISTS idx_material_type ON materials(material_type);`
	DateIndexSchema      = `CREATE INDEX IF NOT EXISTS idx_creation_date ON materials(creation_date);`
	CompositeIndexSchema = `CREATE INDEX IF NOT EXISTS idx_material_type_creation_date ON materials(material_type, creation_date);`
)
