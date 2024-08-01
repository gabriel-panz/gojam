package session

import (
	"github.com/gabriel-panz/gojam/errors"
	"github.com/gabriel-panz/gojam/utils"
	"github.com/patrickmn/go-cache"
)

type SessionService struct {
	Cache *cache.Cache
}

func (s SessionService) GetSession(sId string) (*Session, error) {
	v, found := s.Cache.Get(sId)

	if !found {
		return nil, errors.ErrNotFound
	}

	return v.(*Session), nil
}

func (s SessionService) CreateSession(creatorToken string) (*Session, error) {
	sId, err := utils.GenerateRandomId()
	if err != nil {
		return nil, err
	}
	_, found := s.Cache.Get(sId)

	if found {
		return s.CreateSession(creatorToken)
	}

	ses := NewSession(sId, creatorToken)

	s.Cache.Add(sId, ses, cache.DefaultExpiration)
	return ses, nil
}

func (s SessionService) LeaveSession(sId string, token string) error {
	v, found := s.Cache.Get(sId)
	if !found {
		return errors.ErrNotFound
	}
	ses := v.(*Session)

	delete(ses.Conns, token)

	return s.Cache.Replace(sId, ses, cache.DefaultExpiration)
}

func (s SessionService) DeleteSession(sId string, token string) error {
	v, found := s.Cache.Get(sId)

	if found {
		return nil
	}

	ses := v.(*Session)
	if ses.DjToken != token {
		return errors.ErrUnauthorized
	}

	s.Cache.Delete(sId)
	return nil
}
