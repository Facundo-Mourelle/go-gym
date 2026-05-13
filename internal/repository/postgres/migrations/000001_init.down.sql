ALTER TABLE equipment DROP COLUMN IF EXISTS manufacturer;
ALTER TABLE equipment DROP COLUMN IF EXISTS user_id;
ALTER TABLE equipment DROP COLUMN IF EXISTS actual_weight;
ALTER TABLE equipment DROP COLUMN IF EXISTS cable_pulley_type;
ALTER TABLE equipment DROP COLUMN IF EXISTS cable_stack_weights;
ALTER TABLE equipment DROP COLUMN IF EXISTS resistance_profile_name;

DROP TABLE equipment;
DROP TABLE sessions;
DROP TABLE workout_exercises;
DROP TABLE workout_plans;
DROP TABLE exercise_pattern_contributions;
DROP TABLE exercise_equipment_suggestions;
DROP TABLE exercises;
DROP TABLE users;
