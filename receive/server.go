package receive

import (
    "fmt"
    "net"
)

func StartEventServer(udpPort int, eventChannel chan string) error {
    addr := net.UDPAddr{
        IP: net.IPv4(0, 0, 0, 0),
        Port: udpPort,
    }
    conn, err := net.ListenUDP("udp4", &addr)
    if err != nil {
        return err
    }
    fmt.Println("Listening on UDP port", udpPort)
    go listenForEvents(conn, eventChannel)
    return nil
}

func listenForEvents(conn *net.UDPConn, eventChannel chan string) {
    defer conn.Close()
    buffer := make([]byte, 8192)
    for {
        size, _, err := conn.ReadFrom(buffer)
        if err != nil {
            println("event reception error", err.Error())
            continue
        }
        eventChannel <- string(buffer[0:size])
    }
}
