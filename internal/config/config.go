package config

type TgBotConfig interface {
	Token() string
}

type TikTokAPIConfig interface {
	Token() string
}

type OpenAIAPIConfig interface {
	Token() string
}

type PgConfig interface {
	Dsn() string
}
