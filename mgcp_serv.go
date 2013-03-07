package main

import (
"net"
"fmt"
)

func main() {
    netbuf := make([]byte, 4096)
    conn, _ := net.ListenPacket("udp", "127.0.0.1:2427")

    for {
        _, addr, _ := conn.ReadFrom(netbuf)

        go parseMGCP(string(netbuf), conn, addr)
    }
}

func MGCPCommandRespond(msg MGCPCommand, conn net.PacketConn, addr net.Addr) {
    MGCPCommandRespondEcho(msg, conn, addr)
}

func MGCPCommandRespondEcho(msg MGCPCommand, conn net.PacketConn, addr net.Addr) {
    str_reply := fmt.Sprintf("%s", msg)

    if _, err := conn.WriteTo([]byte(str_reply), addr); err != nil {
        fmt.Println(err)
    }

}

func MGCPCommandRespondErr(msg MGCPCommand, conn net.PacketConn, addr net.Addr) {
    var reply MGCPResponse
    reply.ResponseCode = "504"
    reply.TransID = msg.TransID
    reply.ResponseStr = "ERR"

    str_reply := fmt.Sprintf("%s", reply)

    if _, err := conn.WriteTo([]byte(str_reply), addr); err != nil {
        fmt.Println(err)
    }

}
