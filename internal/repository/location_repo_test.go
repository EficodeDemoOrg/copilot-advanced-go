package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"
	"copilot-go-advanced/internal/repository"
	"copilot-go-advanced/internal/testhelpers"
)

func newRepo() repository.LocationRepository {
	return repository.NewLocationRepository()
}

func TestAdd_ReturnsLocationWithID(t *testing.T) {
	repo := newRepo()
	input := testhelpers.MakeLocationCreate()

	loc, err := repo.Add(input)

	require.NoError(t, err)
	assert.NotEmpty(t, loc.ID)
	assert.Equal(t, input.Name, loc.Name)
	assert.Equal(t, input.Lat, loc.Coordinates.Lat)
	assert.Equal(t, input.Lon, loc.Coordinates.Lon)
	assert.WithinDuration(t, time.Now().UTC(), loc.CreatedAt, 5*time.Second)
}

func TestAdd_MultipleEntriesHaveUniqueIDs(t *testing.T) {
	repo := newRepo()
	input := testhelpers.MakeLocationCreate()

	loc1, _ := repo.Add(input)
	loc2, _ := repo.Add(input)

	assert.NotEqual(t, loc1.ID, loc2.ID)
}

func TestGet_ReturnsExistingLocation(t *testing.T) {
	repo := newRepo()
	loc, _ := repo.Add(testhelpers.MakeLocationCreate())

	got, err := repo.Get(loc.ID)

	require.NoError(t, err)
	assert.Equal(t, loc.ID, got.ID)
}

func TestGet_NonExistentReturnsLocationNotFoundError(t *testing.T) {
	repo := newRepo()

	_, err := repo.Get("00000000-0000-0000-0000-000000000000")

	var locErr *apperrors.LocationNotFoundError
	assert.ErrorAs(t, err, &locErr)
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", locErr.ID)
}

func TestListAll_EmptyReturnsEmptySlice(t *testing.T) {
	repo := newRepo()
	locs := repo.ListAll()
	assert.Empty(t, locs)
}

func TestListAll_SortedByCreatedAtAscending(t *testing.T) {
	repo := newRepo()

	// Add slightly offset entries (timestamps will differ due to uuid + time.Now)
	for _, name := range []string{"Alpha", "Beta", "Gamma"} {
		n := name
		repo.Add(testhelpers.MakeLocationCreate(func(lc *models.LocationCreate) { lc.Name = n })) //nolint:errcheck
	}

	locs := repo.ListAll()
	require.Len(t, locs, 3)
	for i := 1; i < len(locs); i++ {
		assert.False(t, locs[i].CreatedAt.Before(locs[i-1].CreatedAt),
			"entry %d should not be before entry %d", i, i-1)
	}
}

func TestDelete_RemovesExistingLocation(t *testing.T) {
	repo := newRepo()
	loc, _ := repo.Add(testhelpers.MakeLocationCreate())

	err := repo.Delete(loc.ID)

	require.NoError(t, err)
	_, getErr := repo.Get(loc.ID)
	assert.Error(t, getErr)
}

func TestDelete_NonExistentReturnsLocationNotFoundError(t *testing.T) {
	repo := newRepo()

	err := repo.Delete("missing-id")

	var locErr *apperrors.LocationNotFoundError
	assert.ErrorAs(t, err, &locErr)
}

func TestUpdate_PartialName(t *testing.T) {
	repo := newRepo()
	loc, _ := repo.Add(testhelpers.MakeLocationCreate())

	newName := "New York"
	updated, err := repo.Update(loc.ID, models.LocationUpdate{Name: &newName})

	require.NoError(t, err)
	assert.Equal(t, "New York", updated.Name)
	assert.Equal(t, loc.Coordinates.Lat, updated.Coordinates.Lat) // unchanged
}

func TestUpdate_PartialCoordinates(t *testing.T) {
	repo := newRepo()
	loc, _ := repo.Add(testhelpers.MakeLocationCreate())

	newLat := 40.71
	updated, err := repo.Update(loc.ID, models.LocationUpdate{Lat: &newLat})

	require.NoError(t, err)
	assert.Equal(t, 40.71, updated.Coordinates.Lat)
	assert.Equal(t, loc.Name, updated.Name) // unchanged
}

func TestUpdate_NonExistentReturnsLocationNotFoundError(t *testing.T) {
	repo := newRepo()

	name := "x"
	_, err := repo.Update("missing-id", models.LocationUpdate{Name: &name})

	var locErr *apperrors.LocationNotFoundError
	assert.ErrorAs(t, err, &locErr)
}
