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
