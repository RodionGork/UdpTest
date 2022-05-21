package main

import (
    "fmt"
    "os"
    "time"
    
    "github.com/rodiongork/udptest/receive"
    "github.com/rodiongork/udptest/collect"
)

func main() {
    eventChannel := make(chan string, 5)
    
    collector, err := collect.StartEventCollector(eventChannel)
    if err != nil {
        println("Failed to start event collector", err.Error())
        os.Exit(1)
    }
    
    if err := receive.StartEventServer(1961, eventChannel); err != nil {
        println("Failed to start event server", err.Error())
        os.Exit(1)
    }
    
    for {
        fmt.Println("Events:", collector.CachedCount())
        time.Sleep(5 * time.Second)
    }
}

