# Test file for ExerciseService and ExercisePostgresRepository

This file will contain integration tests for the exercise-related functionalities.

## Test Cases

### 1. Exercise Creation Tests

- Test successful creation of a custom exercise.
- Test creation from a template.
- Test creation with invalid data (e.g., missing required fields).

### 2. Exercise Listing Tests

- Test listing exercises for a specific user.
- Test filtering by pattern, equipment, and muscle group.
- Test that exercises are returned in alphabetical order.

### 3. Exercise Update Tests

- Test updating an existing user-created exercise.
- Test that system exercises cannot be updated.
- Test that unauthorized users cannot update exercises.

### 4. Exercise Deletion Tests

- Test deleting a user-created exercise.
- Test that system exercises cannot be deleted.
- Test that unauthorized users cannot delete exercises.
