/* mgcp_common.go - common parts of go-mgcp

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
"strings"
"strconv"
)

type MGCPResponse struct {
    ResponseCode string
    TransID string
    ResponseStr string
    Parameters map[string]string
    SDP string
}

type MGCPCommand struct {
    Verb string
    TransID string
    EndpointName string
    MGCPVer string
    Parameters map[string]string
    SDP string
}

func (msg MGCPCommand) String () (string) {
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

func (msg MGCPResponse) String () (string) {
    result := fmt.Sprintf("%s %s %s\n", msg.ResponseCode, msg.TransID, msg.ResponseStr)
    for key, value := range msg.Parameters {
        result += fmt.Sprintf("%s: %s\n", key, value)
    }
    if msg.SDP != "" {
        result += "\n"
        result += msg.SDP
    }
    return result
}

func parseMGCP (packet string, conn net.PacketConn, addr net.Addr) {

    lines := strings.Split(packet, "\n")
    mgcpCommand := strings.Split(lines[0], " ")

    if _, err := strconv.Atoi(mgcpCommand[0]); err == nil {
        parseMGCPResponse(lines, conn, addr)
    }
    parseMGCPCommand(lines, conn, addr)
}

func parseMGCPResponse(lines []string, conn net.PacketConn, addr net.Addr) {
    var msg MGCPResponse
    msg.Parameters = make(map[string]string, 20)

    sdp := 0
    sdpString := ""

    mgcpCommand := strings.Split(lines[0], " ")
    msg.ResponseCode = mgcpCommand[0]
    msg.TransID = mgcpCommand[1]
    msg.ResponseStr = strings.Join(mgcpCommand[2:], " ")

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
    //MGCPResponseRespond(msg, conn, addr)
}

func parseMGCPCommand(lines []string, conn net.PacketConn, addr net.Addr) {
    var msg MGCPCommand
    msg.Parameters = make(map[string]string, 20)

    sdp := 0
    sdpString := ""

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
    MGCPCommandRespond(msg, conn, addr)
}
