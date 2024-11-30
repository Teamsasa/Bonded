package db

import (
	"context"
	"testing"

	"bonded/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestCalendarRepository(t *testing.T) {
	repo := NewCalendarRepository()
	if repo == nil {
		t.Fatal("Failed to initialize CalendarRepository")
	}

	ctx := context.TODO()

	// Save
	calendar := &models.Calendar{
		ID:     "1",
		UserID: "user1",
		Name:   "Test Calendar",
	}
	err := repo.Save(ctx, calendar)
	assert.NoError(t, err, "Failed to save calendar")

	// FindByID
	found, err := repo.FindByID(ctx, "1")
	assert.NoError(t, err, "Failed to find calendar by ID")
	assert.Equal(t, calendar.ID, found.ID, "Calendar ID does not match")
	assert.Equal(t, calendar.UserID, found.UserID, "Calendar UserID does not match")
	assert.Equal(t, calendar.Name, found.Name, "Calendar Name does not match")

	// Update
	calendar.Name = "Updated Calendar"
	err = repo.Update(ctx, calendar)
	assert.NoError(t, err, "Failed to update calendar")

	// FindByUserID
	calendars, err := repo.FindByUserID(ctx, "user1")
	assert.NoError(t, err, "Failed to find calendars by UserID")
	assert.Len(t, calendars, 1, "Expected one calendar")
	assert.Equal(t, "Updated Calendar", calendars[0].Name, "Calendar name was not updated")

	// Delete
	err = repo.Delete(ctx, "1")
	assert.NoError(t, err, "Failed to delete calendar")

	// Verify Deletion
	found, err = repo.FindByID(ctx, "1")
	assert.Error(t, err, "Expected error when finding deleted calendar")
	assert.Nil(t, found, "Expected no calendar to be found after deletion")
}
