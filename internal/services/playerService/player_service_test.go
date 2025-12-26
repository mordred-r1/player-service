package playerservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/mordred-r1/player-service/internal/models"
	ps "github.com/mordred-r1/player-service/internal/services/playerService"
)

// Mocks
type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Create(ctx context.Context, player *models.PlayerState) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}
func (m *mockStorage) Get(ctx context.Context, id string) (*models.PlayerState, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlayerState), args.Error(1)
}
func (m *mockStorage) Update(ctx context.Context, player *models.PlayerState) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}
func (m *mockStorage) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockProducer struct {
	mock.Mock
}

func (m *mockProducer) Produce(ctx context.Context, event *models.PlayerEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) GetPlayer(ctx context.Context, id string) (*models.PlayerState, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlayerState), args.Error(1)
}
func (m *mockCache) SetPlayer(ctx context.Context, p *models.PlayerState, ttl time.Duration) error {
	args := m.Called(ctx, p, ttl)
	return args.Error(0)
}
func (m *mockCache) DeletePlayer(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test Suite
type PlayerServiceSuite struct {
	suite.Suite
	storage  *mockStorage
	producer *mockProducer
	cache    *mockCache
	svc      *ps.PlayerService
}

func (s *PlayerServiceSuite) SetupTest() {
	s.storage = &mockStorage{}
	s.producer = &mockProducer{}
	s.cache = &mockCache{}
	s.svc = ps.NewPlayerService(context.Background(), s.storage, s.producer, s.cache)
}

func (s *PlayerServiceSuite) TestCreate_Success() {
	p := &models.PlayerState{ID: "p1", State: "created"}
	s.storage.On("Create", mock.Anything, p).Return(nil)
	s.producer.On("Produce", mock.Anything, mock.MatchedBy(func(e *models.PlayerEvent) bool { return e.ID == p.ID && e.State == p.State })).Return(nil)
	s.cache.On("SetPlayer", mock.Anything, p, mock.Anything).Return(nil)

	err := s.svc.Create(context.Background(), p)
	s.Require().NoError(err)
	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *PlayerServiceSuite) TestCreate_StorageError() {
	p := &models.PlayerState{ID: "p1", State: "created"}
	s.storage.On("Create", mock.Anything, p).Return(errors.New("db"))

	err := s.svc.Create(context.Background(), p)
	s.Require().Error(err)
	// producer should not be called
	s.producer.AssertNotCalled(s.T(), "Produce", mock.Anything, mock.Anything)
	s.cache.AssertNotCalled(s.T(), "SetPlayer", mock.Anything, mock.Anything, mock.Anything)
}

func (s *PlayerServiceSuite) TestGet_CacheHit() {
	p := &models.PlayerState{ID: "p2", State: "playing"}
	s.cache.On("GetPlayer", mock.Anything, "p2").Return(p, nil)

	res, err := s.svc.Get(context.Background(), "p2")
	s.Require().NoError(err)
	s.Equal(p, res)
	// storage should not be called
	s.storage.AssertNotCalled(s.T(), "Get", mock.Anything, mock.Anything)
}

func (s *PlayerServiceSuite) TestGet_CacheMiss_StorageHit() {
	p := &models.PlayerState{ID: "p3", State: "stopped"}
	s.cache.On("GetPlayer", mock.Anything, "p3").Return(nil, errors.New("miss"))
	s.storage.On("Get", mock.Anything, "p3").Return(p, nil)
	done := make(chan struct{}, 1)
	s.cache.On("SetPlayer", mock.Anything, p, mock.Anything).Run(func(args mock.Arguments) {
		select {
		default:
			// signal once
			done <- struct{}{}
		}
	}).Return(nil)

	res, err := s.svc.Get(context.Background(), "p3")
	s.Require().NoError(err)
	s.Equal(p, res)

	// wait for the async cache population (timeout to avoid flakiness)
	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		s.T().Fatal("timeout waiting for SetPlayer to be called")
	}

	s.storage.AssertExpectations(s.T())
	s.cache.AssertExpectations(s.T())
}

func (s *PlayerServiceSuite) TestUpdate_Success() {
	p := &models.PlayerState{ID: "p4", State: "paused"}
	s.storage.On("Update", mock.Anything, p).Return(nil)
	s.producer.On("Produce", mock.Anything, mock.MatchedBy(func(e *models.PlayerEvent) bool { return e.ID == p.ID && e.State == p.State })).Return(nil)
	s.cache.On("SetPlayer", mock.Anything, p, mock.Anything).Return(nil)

	err := s.svc.Update(context.Background(), p)
	s.Require().NoError(err)
	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func (s *PlayerServiceSuite) TestDelete_Success() {
	s.storage.On("Delete", mock.Anything, "p5").Return(nil)
	s.producer.On("Produce", mock.Anything, mock.MatchedBy(func(e *models.PlayerEvent) bool { return e.ID == "p5" && e.State == "deleted" })).Return(nil)
	s.cache.On("DeletePlayer", mock.Anything, "p5").Return(nil)

	err := s.svc.Delete(context.Background(), "p5")
	s.Require().NoError(err)
	s.storage.AssertExpectations(s.T())
	s.producer.AssertExpectations(s.T())
}

func TestPlayerServiceSuite(t *testing.T) {
	suite.Run(t, new(PlayerServiceSuite))
}
