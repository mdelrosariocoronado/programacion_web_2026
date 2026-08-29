CREATE TABLE IF NOT EXISTS emprendimientos (
    id_emprendimiento SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    url VARCHAR(255),
    rubro VARCHAR(100) NOT NULL,
    descripcion TEXT,
    logo VARCHAR(255),
    contacto VARCHAR(255),
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS usuarios (
    id_usuario SERIAL PRIMARY KEY,
    nombre_completo VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    clave VARCHAR(255) NOT NULL, 
    rol VARCHAR(50) NOT NULL CHECK (rol IN ('emprendedor', 'cliente', 'administrador')),
    id_emprendimiento INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_usuario_emprendimiento 
        FOREIGN KEY (id_emprendimiento) 
        REFERENCES emprendimientos(id_emprendimiento) 
        ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS publicaciones (
    id_publicacion SERIAL PRIMARY KEY,
    id_emprendimiento INT NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    contenido TEXT,
    imagen_url VARCHAR(255), -- Almacena la ruta física o link de la foto
    tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('producto', 'promocion', 'servicio', 'otro')),
    precio NUMERIC(10, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_publicaciones_emprendimientos
        FOREIGN KEY (id_emprendimiento) 
        REFERENCES emprendimientos(id_emprendimiento) 
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS suscripciones (
    id_suscripcion SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_emprendimiento INT NOT NULL,
    fecha_sub TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Relación con la tabla usuario
    CONSTRAINT fk_suscripciones_usuario
        FOREIGN KEY (id_usuario) 
        REFERENCES usuarios(id_usuario) 
        ON DELETE CASCADE,
        
    -- Relación con la tabla de emprendimiento
    CONSTRAINT fk_suscripciones_emprendimiento
        FOREIGN KEY (id_emprendimiento) 
        REFERENCES emprendimientos(id_emprendimiento) 
        ON DELETE CASCADE,
        
    -- Evita que un usuario se suscriba varias veces al mismo emprendimiento
    CONSTRAINT uq_usuario_emprendimiento 
        UNIQUE (id_usuario, id_emprendimiento)
);