package config

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "github.com/lucky-cookie-waf/agent-cookie/models"
)

func LoadFromFile(configPath string) (*models.Config, error) {
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    
    cfg := GetDefault()
    err = yaml.Unmarshal(data, cfg)
    return cfg, err
}

func GetDefault() *models.Config {
    return &models.Config{
        ListenAddr:    "127.0.0.1:8080",
        CentralServer: "http://localhost:3000", // ❗️ 웹서버 도메인 생기면 바꾸기 ❗️
        LogPath:       "/var/log/modsecurity",
        Interval:      "30s",
        Debug:         false,
    }
}

func Load() (*models.Config, error) {
    return LoadFromFile("config.yaml")
}