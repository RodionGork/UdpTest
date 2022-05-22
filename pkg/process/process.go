package process

import (
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"

    "github.com/rodiongork/udptest/pkg/collect"
    "github.com/rodiongork/udptest/pkg/entity"
)

var ParseErrorLines = errors.New("message contains too few lines")
var ParseErrorStart = errors.New("message doesn't start with stars")
var ParseErrorEnd = errors.New("message doesn't end with END")
var ParseErrorHeader = errors.New("message header has bad format")
var ParseErrorTimestamp = errors.New("message has unparsable timestamp")

func StartEventProcessor(eventChannel chan string, eventCollector collect.EventCollector) {
    go collectEvents(eventChannel, eventCollector)
}

func collectEvents(eventChannel chan string, eventCollector collect.EventCollector) {
    for {
        msg := <-eventChannel
        if event, err := parseEvent(msg); err == nil {
            eventCollector.Store(event)
        } else {
            println(err.Error())
        }
    }
}

func parseEvent(msg string) (*entity.Event, error) {
    lines := strings.Split(strings.TrimSpace(msg), "\n")
    
    if len(lines) < 3 {
        return nil, ParseErrorLines
    }
    
    if !strings.HasPrefix(lines[0], "* * * ") {
        return nil, ParseErrorStart
    }
    
    if lines[len(lines) - 1] != "END" {
        return nil, ParseErrorEnd
    }
    
    header := strings.Split(lines[0][6:], " ")
    if len(header) != 3 {
        return nil, ParseErrorHeader
    }
    
    timeStr := header[2]
    ts, err := time.Parse("20060102T150405", timeStr[: len(timeStr) - 4])
    if err != nil {
        return nil, ParseErrorTimestamp
    }

    event := &entity.Event {
        Type: header[0],
        Severity: header[1],
        UnixTimestamp: ts.Unix(),
        UUID: uuid.New().String(),
    }

    if lastSpace := strings.LastIndex(lines[1], " "); lastSpace >= 0 {
        event.Event = strings.TrimSpace(lines[1][:lastSpace])
        event.EventID = strings.TrimSpace(lines[1][lastSpace:])
    }

    event.Description = strings.Join(lines[2 : len(lines) - 1], "\n")

    return event, nil
}
