package config

type Config struct {
	Api      Api      `json:"api"`
	Postgres Postgres `json:"postgres"`
}

type Api struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
}

type Postgres struct {
	User     string `env:"PG_USER,required" `
	Password string `env:"PG_PASS,required"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Db       string `json:"db"`
	Sslmode  string `json:"sslmode"`
	MaxConns int    `json:"max_conns"`
}
