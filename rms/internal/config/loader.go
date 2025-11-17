package config

import (
    "os"
    "time"

    "gopkg.in/yaml.v3"
)

type Core struct {
    AdminEmail        string 	`yaml:"admin_email"`
    AdminPassword     string 	`yaml:"admin_password"`
    AdminPasswordHash string 	`yaml:"admin_password_hash"`
    DBPath            string 	`yaml:"db_path"`
    HTTPS             bool   	`yaml:"https"`
    Port              int    	`yaml:"port"`
    IP                string 	`yaml:"ip"`
	GinMode           string 	`yaml:"gin_mode"`
	TrustedProxies    []string 	`yaml:"trusted_proxies"`
}

type Auth struct {
    JWTKey            string 	`yaml:"jwt_key"`
    JWTTTL            time.Duration 	`yaml:"jwt_ttl"`
    RefreshTokenTTL   time.Duration 	`yaml:"refresh_token_ttl"`
    CORSAllowOrigins  []string 	`yaml:"cors_allow_origins"`
}

type Workshop struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    ImagePath   string `yaml:"image_path"`
    Owner       string `yaml:"owner"`
}

type App struct {
    Config      Core     `yaml:"config"`
    Auth        Auth     `yaml:"auth"`
    CarWorkshop Workshop `yaml:"car_workshop"`
}

func Load(path string) (App, error) {
    var cfg App
    f, err := os.Open(path)
    if err != nil {
        return cfg, err
    }
    defer f.Close()
    if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
        return cfg, err
    }
    if cfg.Config.Port == 0 {
        cfg.Config.Port = 10800
    }
    if cfg.Config.IP == "" {
        cfg.Config.IP = "0.0.0.0"
    }
    if cfg.Config.DBPath == "" {
        cfg.Config.DBPath = "data/database.sqlite3"
    }

    if cfg.Config.GinMode == "" {
        cfg.Config.GinMode = "debug"
    }

    return cfg, nil
}