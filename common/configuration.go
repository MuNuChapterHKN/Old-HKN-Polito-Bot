package common

import (
    "encoding/json"
    "os"
    "log"
)

type Configuration struct {
    TelegramAPI string
    ApiAiToken string
    ApiAiSessionId string
    ApiAiQueryLang string
}

func LoadConfiguration() *Configuration {
    file, err := os.Open("config/config.json")
    if err != nil {
        log.Printf("error, invalid config file");
    }
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err = decoder.Decode(&configuration)
    if err != nil {
        log.Printf("error:", err)
    }
    // Get the variables also from the environment
    token := os.Getenv("TelegramAPI")
    if token != "" {
        configuration.TelegramAPI = token
    }
    token = os.Getenv("ApiAiToken")
    if token != "" {
        configuration.ApiAiToken = token
    }
    token = os.Getenv("ApiAiSessionId")
    if token != "" {
        configuration.ApiAiSessionId = token
    }
    token = os.Getenv("ApiAiQueryLang")
    if token != "" {
        configuration.ApiAiQueryLang = token
    }
    return &configuration
}
