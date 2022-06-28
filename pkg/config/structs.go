package config

type Config struct {
	App     Application
	DB      Database
	Auth    Auth
	Log     Logger
	Tinkoff TinkoffAPI
}

type Application struct {
	Host string
	Port string
}

type Logger struct {
	Level string
}

type Database struct {
	Name    string
	Pass    string
	Host    string
	User    string
	Type    string
	Port    string
	SSLMode string
}

type Auth struct {
	JWTSecret string
}

type TinkoffAPI struct {
	Token string
	URL   string
}
