package jwt

type Config struct {
	SecretKey string `yaml:"secretKey"`
	Domain    string `yaml:"domain"`
}
