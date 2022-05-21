package collect

type EventCollector interface {
    CachedCount() int
}

type collector struct {
    events []string
}

func StartEventCollector(eventChannel chan string) (EventCollector, error) {
    ec := &collector{}
    go ec.collectEvents(eventChannel)
    return ec, nil
}

func (ec *collector) collectEvents(eventChannel chan string) {
    for {
        event := <-eventChannel
        ec.events = append(ec.events, event)
    }
}

func (ec *collector) CachedCount() int {
    return len(ec.events)
}
