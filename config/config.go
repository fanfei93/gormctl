package config

type Config struct {
	Host        string `validate:"required"`
	Port        int
	User        string `validate:"required"`
	Password    string `validate:"required"`
	Database    string `validate:"required"`
	Table       string `validate:"required"`
	OutDir      string `validate:"required"`
	PackageName string
}
