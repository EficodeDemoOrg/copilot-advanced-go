package repository

import (
	"sort"
	"sync"
	"time"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"

	"github.com/google/uuid"
)

// LocationRepository defines storage operations for locations.
type LocationRepository interface {
	Add(input models.LocationCreate) (models.Location, error)
	Get(id string) (models.Location, error)
	ListAll() []models.Location
	Update(id string, input models.LocationUpdate) (models.Location, error)
	Delete(id string) error
}

type locationRepo struct {
	mu    sync.RWMutex
	store map[string]models.Location
}

// NewLocationRepository returns an in-memory LocationRepository.
func NewLocationRepository() LocationRepository {
	return &locationRepo{
		store: make(map[string]models.Location),
	}
}

func (r *locationRepo) Add(input models.LocationCreate) (models.Location, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	loc := models.Location{
		ID:   uuid.New().String(),
		Name: input.Name,
		Coordinates: models.Coordinates{
			Lat: input.Lat,
			Lon: input.Lon,
		},
		CreatedAt: time.Now().UTC(),
	}
	r.store[loc.ID] = loc
	return loc, nil
}

func (r *locationRepo) Get(id string) (models.Location, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loc, ok := r.store[id]
	if !ok {
		return models.Location{}, &apperrors.LocationNotFoundError{ID: id}
	}
	return loc, nil
}

func (r *locationRepo) ListAll() []models.Location {
	r.mu.RLock()
	defer r.mu.RUnlock()

	locs := make([]models.Location, 0, len(r.store))
	for _, l := range r.store {
		locs = append(locs, l)
	}
	sort.Slice(locs, func(i, j int) bool {
		return locs[i].CreatedAt.Before(locs[j].CreatedAt)
	})
	return locs
}

func (r *locationRepo) Update(id string, input models.LocationUpdate) (models.Location, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	loc, ok := r.store[id]
	if !ok {
		return models.Location{}, &apperrors.LocationNotFoundError{ID: id}
	}

	if input.Name != nil {
		loc.Name = *input.Name
	}
	if input.Lat != nil {
		loc.Coordinates.Lat = *input.Lat
	}
	if input.Lon != nil {
		loc.Coordinates.Lon = *input.Lon
	}

	r.store[id] = loc
	return loc, nil
}

func (r *locationRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return &apperrors.LocationNotFoundError{ID: id}
	}
	delete(r.store, id)
	return nil
}
