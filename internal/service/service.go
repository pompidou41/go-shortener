package service

import (
	"context"
	"errors"
	"pompidou41/go-shortener/internal/codegen"
	"pompidou41/go-shortener/internal/storage"
)

type Service interface {
	Shorten(ctx context.Context, longUrl string) (string, error)
	Resolve(ctx context.Context, code string) (string, error)
}

type service struct {
	store   storage.Storage
	codeGen codegen.CodeGenerator
}

func New(store storage.Storage) *service {
	return &service{
		store: store,
	}
}

func (s *service) Shorten(ctx context.Context, longUrl string) (string, error) {
	for range 3 {
		code, err := s.codeGen.GenerateRandomString(8)

		if err != nil {
			return "", err
		}

		err = s.store.Save(ctx, code, longUrl)

		if err == nil {
			return code, err
		}

		if errors.Is(err, storage.ErrCodeExists) {
			continue
		}

		return "", err
	}

	return "", errors.New("Could not generate unique code after 3 attemps")
}

func (s *service) Resolve(ctx context.Context, code string) (string, error) {
	longUrl, err := s.store.Get(ctx, code)

	if err != nil {
		return "", err
	}

	return longUrl, nil
}
