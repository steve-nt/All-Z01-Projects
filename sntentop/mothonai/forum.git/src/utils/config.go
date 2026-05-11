package utils

type Configuration struct {
	DbPath string
	Ip string
	Port string
}

func DefaultConfiguration() Configuration {
	return Configuration{
		DbPath: "./db.db",
		Ip: "127.0.0.1",
		Port: "8080",
	}
}
