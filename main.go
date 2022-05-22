package main

import (
    "os"
    "strconv"
    
    "github.com/rodiongork/udptest/pkg/collect"
    "github.com/rodiongork/udptest/pkg/process"
    "github.com/rodiongork/udptest/pkg/receive"
    "github.com/rodiongork/udptest/pkg/web"
)

func main() {
    eventChannel := make(chan string, 5)
    eventCollector := collect.New()

    process.StartEventProcessor(eventChannel, eventCollector)

    portToListen := intFromEnv("port", 1961)

    if err := receive.StartEventServer(portToListen, eventChannel); err != nil {
        println("Failed to start event server", err.Error())
        os.Exit(2)
    }

    portForUi := intFromEnv("uiport", 4004)

    weberr := web.RunWebServer(portForUi, eventCollector)
    println("Web-server ends with", weberr.Error())
}

func intFromEnv(key string, defaultValue int) int {
    if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
        return v
    }
    return defaultValue
}
