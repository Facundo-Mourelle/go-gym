CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE exercise_equipment_suggestions (
    exercise_id VARCHAR(255) REFERENCES exercises(id) ON DELETE CASCADE,
    equipment_type VARCHAR(100) NOT NULL,
    PRIMARY KEY (exercise_id, equipment_type)
);

CREATE TABLE exercise_pattern_contributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id VARCHAR(255) REFERENCES exercises(id) ON DELETE CASCADE,
    pattern VARCHAR(100) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    contribution FLOAT NOT NULL,
    range_of_motion VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE equipment (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    manufacturer VARCHAR(255) DEFAULT '',
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
    actual_weight DOUBLE PRECISION DEFAULT 0,
    cable_pulley_type VARCHAR(100) DEFAULT '',
    cable_stack_weights DOUBLE PRECISION[] DEFAULT '{}',
    cable_weight_increment DOUBLE PRECISION DEFAULT 0,
    resistance_profile_id VARCHAR(255),
    resistance_profile_name VARCHAR(255) DEFAULT '',
    movement_pattern VARCHAR(100) DEFAULT ''
);

CREATE TABLE workout_plans (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workout_exercises (
    id VARCHAR(255) PRIMARY KEY,
    workout_plan_id VARCHAR(255) REFERENCES workout_plans(id) ON DELETE CASCADE,
    exercise_id VARCHAR(255) REFERENCES exercises(id) ON DELETE SET NULL,
    "order" INTEGER NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER NOT NULL,
    reps_in_reserve INTEGER,
    notes TEXT
);

CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    workout_plan_id VARCHAR(255) REFERENCES workout_plans(id) ON DELETE SET NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    performed_sets JSONB,
    notes TEXT
);

CREATE TABLE routines (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    movement_patterns TEXT[] NOT NULL DEFAULT '{}',
    is_preset BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_exercises_user_id ON exercises(user_id);
CREATE INDEX idx_exercise_patterns_exercise_id ON exercise_pattern_contributions(exercise_id);
CREATE INDEX idx_exercise_patterns_pattern ON exercise_pattern_contributions(pattern);
CREATE INDEX idx_equipment_user_id ON equipment(user_id);
CREATE INDEX idx_equipment_type ON equipment(type);
CREATE INDEX idx_workout_plans_user_id ON workout_plans(user_id);
CREATE INDEX idx_workout_exercises_workout_plan_id ON workout_exercises(workout_plan_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_started_at ON sessions(started_at DESC);
CREATE INDEX idx_routines_user_id ON routines(user_id);
