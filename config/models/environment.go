package model

import "crypto/ecdsa"

type Environment struct {
	PrivateKey        *ecdsa.PrivateKey
	PublicKey         *ecdsa.PublicKey
	Redis             RedisEnvironment
	CloudFlare        CloudFlareEnvironment
	Cache             CacheEnvironment
	AvatarPlaceholder string `env:"AVATAR_PLACEHOLDER"`
	RabbitMQURL       string `env:"RABBITMQ_URL"`
	APIPort           string `env:"API_PORT"`
	ConnectionString  string `env:"CONNECTION_STRING"`
	FrontURL          string `env:"FRONT_URL"`
}

type RedisEnvironment struct {
	DB       int    `env:"REDIS_DB"`
	Address  string `env:"REDIS_ADDRESS"`
	Password string `env:"REDIS_PASSWORD"`
}

type CloudFlareEnvironment struct {
	CloudFlareAccountAPI string `env:"CLOUD_FLARE_ACCOUNT_API"`
	CloudFlareApiKey     string `env:"CLOUD_FLARE_API_KEY"`
}

type CacheEnvironment struct {
	SessionExp      int `env:"SESSION_EXP"`
	CacheExp        int `env:"CACHE_EXP"`
	Hash2FADuration int `env:"HASH_2FA_DURATION"`
	Code2FADuration int `env:"CODE_2FA_DURATION"`
}
