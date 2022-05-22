package collect

import (
    "context"
    "database/sql"
    "sync"
    "time"

    "github.com/rodiongork/udptest/pkg/entity"
    _ "github.com/mattn/go-sqlite3"
)

const maxCachedCount = 10000
const dbName = "./events.db"

type EventCollector interface {
    Store(event *entity.Event)
    Latest(count int, typ, sev string) []*entity.Event
    CachedCount() int
    TotalCount() int
    Retrieve(uuid string) (*entity.Event, error)
}

type collector struct {
    events []*entity.Event
    totalCount int
    lock sync.Mutex
    toPersist chan *entity.Event
}

func New() EventCollector {
    ec := &collector{ toPersist: make(chan *entity.Event, 100) }
    go ec.persistentStorage()
    return ec
}

func (ec *collector) Store(event *entity.Event) {
    ec.lock.Lock()
    defer ec.lock.Unlock()
    ec.events = append(ec.events, event)
    ec.totalCount++
    if len(ec.events) >= maxCachedCount {
        ec.events = ec.events[len(ec.events) - maxCachedCount * 2 / 3 :]
    }
    ec.toPersist <- event
}

func (ec *collector) CachedCount() int {
    return len(ec.events)
}

func (ec *collector) TotalCount() int {
    return ec.totalCount
}

func (ec *collector) Latest(count int, typ, sev string) []*entity.Event {
    events := ec.events
    result := make([]*entity.Event, 0, count)

    for i := len(events) - 1; i >= 0; i-- {
        event := events[i]
        if (typ != "" && event.Type != typ) || (sev != "" && event.Severity != sev) {
            continue
        }
        result = append(result, event)
        if len(result) >= count {
            break
        }
    }

    return result
}

func (ec *collector) persistentStorage() {
    for {
        events := make([]*entity.Event, 0)
        done:
        for len(events) < 50 {
            select {
                case event := <-ec.toPersist:
                    events = append(events, event)
                default:
                    break done
            }
        }
        if len(events) > 0 {
            persist(events)
        } else {
            time.Sleep(30 * time.Millisecond)
        }
    }
}

func persist(events []*entity.Event) {
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        println("database opening error", err.Error())
        return
    }
    defer db.Close()
    
    ctx := context.Background()
    
    tx, _ := db.BeginTx(ctx, nil)

    stmt, err := tx.Prepare(`insert into events (uuid, type, severity, ts, event, eventid, descr)
        values (?,?,?,?,?,?,?)`)
    
    for _, event := range events {
        stmt.Exec(event.UUID, event.Type, event.Severity, event.UnixTimestamp,
            event.Event, event.EventID, event.Description)
    }
    
    tx.Commit()
    
}

func (ec *collector) Retrieve(uuid string) (*entity.Event, error) {
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    event := entity.Event{UUID: uuid}

    err = db.QueryRow(`select type, severity, ts, event, eventid, descr from events where uuid = ?`, uuid).
        Scan(&event.Type, &event.Severity, &event.UnixTimestamp,
            &event.Event, &event.EventID, &event.Description)
    if err != nil {
        return nil, err
    }

    return &event, nil
}
