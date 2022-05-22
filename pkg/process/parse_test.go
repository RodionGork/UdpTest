package process

import (
    "testing"
)

const sampleMessage =
`* * * FIRE WARN 20210320T221735Z003
HOUSE ON FIRE                 13-41
Fire started on 2nd floor by kids!!
END
`

func TestParse(t *testing.T) {
    event := parseEvent(sampleMessage)
    if event == nil {
        t.Errorf("Parsing failed")
    }
    t.Log(event.UUID, event.Type, event.Severity, event.UnixTimestamp)
    t.Log(event.Event, event.EventID)
    t.Log(event.Description)
}
