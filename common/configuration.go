package common

import (
    "encoding/json"
    "os"
    "log"
    "reflect"
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
    params := []string{
        "TelegramAPI",
        "ApiAiToken",
        "ApiAiSessionId",
        "ApiAiQueryLang",
    }
    r := reflect.ValueOf(configuration)
    for _, p := range params {
        token := os.Getenv(p)
        v := reflect.Indirect(r).FieldByName(p)
        if v.CanSet() == true {
            v.SetString(token)
        }
    }
    return &configuration
}
