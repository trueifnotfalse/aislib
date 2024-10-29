// Copyright (c) 2015, Marios Andreopoulos.
//
// This file is part of aislib.
//
//  Aislib is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//  Aislib is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
// along with aislib.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
    "bufio"
    "fmt"
    "github.com/trueifnotfalse/aislib"
    "os"
)

func main() {
    in := bufio.NewScanner(os.Stdin)
    in.Split(bufio.ScanLines)

    send := make(chan string, 1024*8)
    receive := make(chan aislib.Message, 1024*8)
    failed := make(chan error, 1024*8)
    stop := make(chan bool)

    done := make(chan bool)

    go aislib.Router(send, receive, stop, failed)

    go func() {
        var message aislib.Message
        var err error
        for {
            select {
            case message = <-receive:
                switch message.Type {
                case 1, 2, 3:
                    t, _ := aislib.DecodeClassAPositionReport(message.Payload)
                    fmt.Println(t)
                case 4:
                    t, _ := aislib.DecodeBaseStationReport(message.Payload)
                    fmt.Println(t)
                case 5:
                    t, _ := aislib.DecodeStaticVoyageData(message.Payload)
                    fmt.Println(t)
                case 8:
                    t, _ := aislib.DecodeBinaryBroadcast(message.Payload)
                    fmt.Println(t)
                case 18:
                    t, _ := aislib.DecodeClassBPositionReport(message.Payload)
                    fmt.Println(t)
                case 255:
                    done <- true
                default:
                    fmt.Printf("=== Message Type %2d ===\n", message.Type)
                    fmt.Printf(" Unsupported type \n\n")
                }
            case err = <-failed:
                fmt.Println(err.Error())
            }
        }
    }()

    for in.Scan() {
        send <- in.Text()
    }
    close(send)
    <-done
}
