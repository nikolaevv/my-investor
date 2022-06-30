package config

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
}
