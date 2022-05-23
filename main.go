package main

import (
    "os"
    "strconv"
    "strings"
    
    "github.com/rodiongork/udptest/pkg/collect"
    "github.com/rodiongork/udptest/pkg/process"
    "github.com/rodiongork/udptest/pkg/receive"
    "github.com/rodiongork/udptest/pkg/web"
)

func main() {
    eventChannel := make(chan string, 5)
    eventCollector := collect.New()

    process.StartEventProcessor(eventChannel, eventCollector)

    portsToListen := intListFromEnv("UDP_PORTS", []int{1961})

    for _, port := range portsToListen {
        if err := receive.StartEventServer(port, eventChannel); err != nil {
            println("Failed to start event server on port", port, err.Error())
            os.Exit(2)
        }
    }

    portForUi := intFromEnv("UI_PORT", 4004)

    weberr := web.RunWebServer(portForUi, eventCollector)
    println("Web-server ends with", weberr.Error())
}

func intFromEnv(key string, defaultValue int) int {
    if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
        return v
    }
    return defaultValue
}

func intListFromEnv(key string, defaultValue []int) []int {
    values := os.Getenv(key)
    if values == "" {
        return defaultValue
    }
    result := make([]int, 0)
    for _, value := range strings.Split(values, ",") {
        if port, err := strconv.Atoi(value); err == nil {
            result = append(result, port)
        }
    }
    return result
}
