package session

import (
	"github.com/gabriel-panz/gojam/errors"
	"github.com/gabriel-panz/gojam/utils"
	"github.com/patrickmn/go-cache"
)

type SessionService struct {
	Cache *cache.Cache
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

	ses := Session{
		ID:         sId,
		UserTokens: []string{creatorToken},
		DjToken:    creatorToken,
	}

	s.Cache.Add(sId, ses, cache.DefaultExpiration)
	return &ses, nil
}

func (s SessionService) JoinSession(sId string, token string) error {
	v, found := s.Cache.Get(sId)
	if !found {
		return errors.ErrNotFound
	}
	ses := v.(*Session)
	if len(ses.UserTokens) >= 10 {
		return errors.ErrSessionFull
	}

	ses.UserTokens = append(ses.UserTokens, token)
	return nil
}

func (s SessionService) LeaveSession(sId string, token string) error {
	v, found := s.Cache.Get(sId)
	if !found {
		return errors.ErrNotFound
	}
	ses := v.(*Session)
	var rIndex int
	for i, t := range ses.UserTokens {
		if t == token {
			rIndex = i
		}
	}

	if ses.DjToken == token {
		ses.DjToken = ses.UserTokens[0]
	}

	if len(ses.UserTokens) == 1 {
		return s.DeleteSession(sId, token)
	}

	ses.UserTokens = append(ses.UserTokens[:rIndex], ses.UserTokens[rIndex+1:]...)

	s.Cache.Replace(sId, ses, cache.DefaultExpiration)

	return nil
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
