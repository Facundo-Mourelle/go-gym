CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE exercises (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id VARCHAR(255) REFERENCES users(id),
    source VARCHAR(50),
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE exercise_equipment_suggestions (
    exercise_id VARCHAR(255) REFERENCES exercises(id),
    equipment_type VARCHAR(100) NOT NULL,
    PRIMARY KEY (exercise_id, equipment_type)
);

CREATE TABLE exercise_pattern_contributions (
    exercise_id VARCHAR(255) REFERENCES exercises(id),
    pattern VARCHAR(100) NOT NULL,
    contribution FLOAT NOT NULL,
    range_of_motion VARCHAR(50),
    notes TEXT,
    PRIMARY KEY (exercise_id, pattern)
);

CREATE TABLE workout_plans (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workout_exercises (
    id VARCHAR(255) PRIMARY KEY,
    workout_plan_id VARCHAR(255),
    exercise_id VARCHAR(255) REFERENCES exercises(id),
    "order" INTEGER NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    reps_in_reserve INTEGER,
    notes TEXT
);

CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id),
    workout_plan_id VARCHAR(255),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    performed_sets JSONB,
    notes TEXT
);

CREATE TABLE equipment (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    resistance_profile_id VARCHAR(255)
);

CREATE TABLE routines (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    movement_patterns TEXT[] NOT NULL DEFAULT '{}',
    is_preset BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_routines_user_id ON routines(user_id);
