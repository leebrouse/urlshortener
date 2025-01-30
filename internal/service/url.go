package service

import (
	"context"
	"fmt"
	"time"

	"githum.com/leebrouse/urlshortener/internal/model"
	"githum.com/leebrouse/urlshortener/internal/repo"
)

type ShortCodeGenerator interface {
	GengerateShortCode() (string, error)
}

type Cacher interface {
	SetURL(ctx context.Context, url repo.Url) error
	GetURL(ctx context.Context, shortCode string) (*repo.Url, error)
}

type URLService struct {
	querier             repo.Querier
	shortCodeGenerateor ShortCodeGenerator
	defaultDuration     time.Duration
	cache               Cacher
	baseURL             string
}

func (s *URLService) CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error) {
	var shortcode string
	var isCustom bool
	var expires_at time.Time

	if req.CustomCode != "" {
		//检查是否已经存在
		isavailable, err := s.querier.IsShortCodeAvailable(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}
		if !isavailable {
			return nil, fmt.Errorf("custom code %s is existed", req.CustomCode)
		}

		shortcode = req.CustomCode
		isavailable = true
	} else {
		code, err := s.getShortCode(ctx, 0)
		if err != nil {
			return nil, err
		}
		shortcode = code
	}

	if req.Duration != nil {
		expires_at = time.Now().Add(s.defaultDuration)
	} else {
		expires_at = time.Now().Add(time.Hour * time.Duration(*req.Duration))
	}

	//插入数据库
	url, err := s.querier.CreateURL(ctx, repo.CreateURLParams{
		OriginalUrl: req.OriginalURL,
		ShortUrl:    shortcode,
		IsCustom:    isCustom,
		ExpiresAt:   expires_at,
	})
	if err != nil {
		return nil, err
	}
	//存入缓存
	if err := s.cache.SetURL(ctx, url); err != nil {
		return nil, err
	}

	return &model.CreateURLResponse{
		ShortUrl:  s.baseURL + "/" + url.ShortUrl,
		ExpiresAt: url.ExpiresAt,
	}, nil
}

func (s *URLService) GetURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	//访问redis,是否shortCode已经存在
	url, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url != nil {
		return url.OriginalUrl, nil
	}

	//访问数据库
	url2,err:=s.querier.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	//存入缓存
	if err:=s.cache.SetURL(ctx, url2);err!=nil{
		return "", err
	}
	return url2.OriginalUrl,nil
}

func (s *URLService) getShortCode(ctx context.Context, n int) (string, error) {
	if n > 5 {
		return "", fmt.Errorf("try too many times")
	}
	shortCode, err := s.shortCodeGenerateor.GengerateShortCode()
	if err != nil {
		return "", err
	}
	isavailable, err := s.querier.IsShortCodeAvailable(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if isavailable {
		return shortCode, nil
	}

	return s.getShortCode(ctx, n+1)
}
