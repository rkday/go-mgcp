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

func (msg MGCPMessage) String () (string) {
    result := fmt.Sprintf("%s %s %s %s\n", msg.Verb, msg.TransID, msg.EndpointName, msg.MGCPVer)
    for key, value := range msg.Parameters {
        result += fmt.Sprintf("%s: %s\n", key, value)
    }
    if msg.SDP != "" {
        result += "\n"
        result += msg.SDP
    }
    return result
}

func main() {
    netbuf := make([]byte, 4096)
    conn, _ := net.ListenPacket("udp", "127.0.0.1:2427")

    _, addr, _ := conn.ReadFrom(netbuf)

    fmt.Println(addr)

    fmt.Printf("%s\n\n\n", netbuf)

    response := parseMGCP(string(netbuf))

    fmt.Printf("%s\n", []byte(response))
    if _, err := conn.WriteTo([]byte(response), addr); err != nil {
        fmt.Println(err)
    }

}

func parseMGCP (packet string) (string) {

    var msg MGCPMessage
    msg.Parameters = make(map[string]string, 20)

    sdp := 0
    sdpString := ""

    lines := strings.Split(packet, "\n")
    mgcpCommand := strings.Split(lines[0], " ")
    
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
            msg.Parameters[mgcpParam[0]] = mgcpParam[1]
        } else {
            sdpString += line
            sdpString += "\n"
        }
    }

    if sdp == 1 {
        msg.SDP = sdpString
    }
    return fmt.Sprintf("%s", msg)
}
