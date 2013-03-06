package main

import (
"net"
"fmt"
"strings"
)

type MGCPMessage struct {
    Verb string
    TransID string
    EndpointName string
    MGCPVer string
    Parameters map[string]string
    SDP string
}

func main() {
    netbuf := make([]byte, 65535)
    conn, _ := net.ListenPacket("udp", "127.0.0.1:2427")

    _, addr, _ := conn.ReadFrom(netbuf)

    fmt.Println(addr)
    fmt.Printf("%s", netbuf)

    parseMGCP(string(netbuf))

    response := "500 Internal Server Error"

    conn.WriteTo([]byte(response), addr)

}

func parseMGCP (packet string) {

    var msg MGCPMessage
    msg.Parameters = make(map[string]string, 500)

    sdp := 0
    sdpString := ""

    lines := strings.Split(packet, "\n")
    mgcpCommand := strings.Split(lines[0], " ")
    
    fmt.Printf("Verb is %s\nTransaction ID is %s\nEndpoint name is %s\nMGCP version is %s\n", mgcpCommand[0], mgcpCommand[1], mgcpCommand[2], mgcpCommand[3:])
    
    msg.Verb = mgcpCommand[0]
    msg.TransID = mgcpCommand[1]
    msg.EndpointName = mgcpCommand[2]
    msg.MGCPVer = strings.Join(mgcpCommand[3:], " ")

    for _, line := range lines[1:] {
        if line == "" {
            sdp = 1
            continue
        }
    
        if sdp == 0 {
            mgcpParam := strings.Split(line, ": ")
            fmt.Printf("Value of %s is %s\n", mgcpParam[0], mgcpParam[1])
            msg.Parameters[mgcpParam[0]] = mgcpParam[1]
        } else {
            sdpString += line
            sdpString += "\n"
        }
    }

    if sdp == 1 {
        fmt.Printf("SDP is %s\n", sdpString)
        msg.SDP = sdpString
    }
    fmt.Printf("%s", msg)
}
