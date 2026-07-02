package actions

type Actions struct {
	repo    Repository
	mailer  Mailer
	jwt     JwtManager[TokenPayload]
	sitemap Sitemap
}

func New(repo Repository, mailer Mailer, jwt JwtManager[TokenPayload], sitemap Sitemap) *Actions {
	return &Actions{
		repo:    repo,
		mailer:  mailer,
		jwt:     jwt,
		sitemap: sitemap,
	}
}
