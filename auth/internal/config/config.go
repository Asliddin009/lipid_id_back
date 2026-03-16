package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                   string
	Host                   string
	Timeout                int
	DBDSN                  string
	DBSSL                  string
	DBTimeout              int
	JWTSecretKey           string
	AccessTokenExpiration  int
	RefreshTokenExpiration int
}

func NewConfig() (*Config, error) {
	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}
	host, err := getEnv("HOST")
	if err != nil {
		return nil, err
	}
	timeout, err := getEnvAsInt("SERVER_TIMEOUT")
	if err != nil {
		return nil, err
	}
	dbTimeout, err := getEnvAsInt("DB_TIMEOUT")
	if err != nil {
		return nil, err
	}
	jwtSecret, err := getEnv("JWT_SECRET_KEY")
	if err != nil {
		return nil, err
	}
	accessTokenExpiration, err := getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRATION")
	if err != nil {
		return nil, err
	}
	refreshTokenExpiration, err := getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRATION")
	if err != nil {
		return nil, err
	}
	dbSSL, err := getEnv("POSTGRES_USE_SSL")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USE_SSL из переменной окружения")
	}

	dbDSN, err := getDBDSN()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:                   port,
		Host:                   host,
		Timeout:                timeout,
		DBDSN:                  dbDSN,
		DBSSL:                  dbSSL,
		DBTimeout:              dbTimeout,
		JWTSecretKey:           jwtSecret,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
	}, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}

func getEnvAsInt(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not set", key)
	}
	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid integer", key)
	}
	return intValue, nil
}

func getDBDSN() (string, error) {
	dbHost, err := getEnv("POSTGRES_HOST")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_HOST из переменной окружения")
	}
	dbPort, err := getEnv("POSTGRES_PORT")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_PORT из переменной окружения")
	}
	dbUser, err := getEnv("POSTGRES_USER")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USER из переменной окружения")
	}
	dbPassword, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_PASSWORD из переменной окружения")
	}
	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_DB из переменной окружения")
	}
	dbSSL, err := getEnv("POSTGRES_USE_SSL")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USE_SSL из переменной окружения")
	}
	// Формирование строки подключения к базе данных
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSL)
	return dbDSN, nil

}
