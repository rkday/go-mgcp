/* mgcp_serv.go - MGCP server in Go

Copyright (C) 2013 Robert Day <rkd@rkd.me.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

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
