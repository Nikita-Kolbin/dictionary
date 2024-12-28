package service

type tgClient interface {
}

type Service struct {
	jwtSecret string
	repo      repository
	storage   objectStorage
	cache     cache
}

func New(repo repository, storage objectStorage, cache cache, jwtSecret string) *Service {
	return &Service{
		jwtSecret: jwtSecret,
		repo:      repo,
		storage:   storage,
		cache:     cache,
	}
}
