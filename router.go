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

package aislib

import (
    "fmt"
    "strconv"
    "strings"
)

// Router accepts AIS radio sentences and process them. It checks their checksum,
// and AIS identifiers. If they are valid it tries to assemble the payload if it spans
// on multiple sentences. Upon success it returns the AIS Message at the out channel.
// Failed sentences go to the err channel.
// If the in channel is closed, then it sends a message with type 255 at the out channel.
// Your function can check for this message to know when it is safe to exit the program.
func Router(in chan string, out chan Message, stop chan bool, failed chan error) {
    count, ccount, padding := 0, 0, 0
    size, id := "0", "0"
    payload := ""
    var (
        err      error
        cache    [5]string
        sentence string
    )

    aisIdentifiers := map[string]bool{
        "ABVD": true, "ADVD": true, "AIVD": true, "ANVD": true, "ARVD": true,
        "ASVD": true, "ATVD": true, "AXVD": true, "BSVD": true, "SAVD": true,
    }
    for {
        select {
        case <-stop:
            return
        case sentence = <-in:
            if len(sentence) == 0 { // Do not process empty lines
                failed <- fmt.Errorf("empty line id '%s'", sentence)
                break
            }
            tokens := strings.Split(sentence, ",") // I think this takes the major portion of time for this function (after benchmarking)

            if !Nmea183ChecksumCheck(sentence) { // Checksum check
                failed <- fmt.Errorf("checksum failed in '%s'", sentence)
                break
            }

            if !aisIdentifiers[tokens[0][1:5]] { // Check for valid AIS identifier
                failed <- fmt.Errorf("sentence isn't AIVDM/AIVDO in '%s'", sentence)
                break
            }

            if tokens[1] == "1" { // One sentence message, process it immediately
                padding, _ = strconv.Atoi(tokens[6][:1])
                out <- Message{Type: MessageType(tokens[5]), Payload: tokens[5], Padding: uint8(padding)}
                if count > 1 { // Invalidate cache
                    for i := 0; i < count; i++ {
                        failed <- fmt.Errorf("incomplete/out of order span sentence in '%s'", cache[i])
                    }
                    count = 0
                    payload = ""
                }
            } else { // Message spans across sentences.
                ccount, err = strconv.Atoi(tokens[2])
                if err != nil {
                    failed <- fmt.Errorf("HERE '%s' in '%s'", tokens[2], sentence)
                    break
                }
                if ccount != count+1 || // If there are sentences with wrong seq.number in cache send them as failed
                    tokens[3] != id && count != 0 || // If there are sentences with different sequence id in cache , send old parts as failed
                    tokens[1] != size && count != 0 { // If there messages with wrong size in cache, send them as failed
                    for i := 0; i < count; i++ {
                        failed <- fmt.Errorf("incomplete/out of order span sentence in '%s'", cache[i])
                    }
                    if ccount != 1 { // The current one is invalid too
                        failed <- fmt.Errorf("incomplete/out of order span sentence in '%s'", sentence)
                        count = 0
                        payload = ""
                        break
                    }
                    count = 0
                    payload = ""
                }
                payload += tokens[5]
                cache[ccount-1] = sentence
                count++
                if ccount == 1 { // First message in sequence, get size and id
                    size = tokens[1]
                    id = tokens[3]
                } else if size == tokens[2] && count == ccount { // Last message in sequence, send it and clean up.
                    padding, _ = strconv.Atoi(tokens[6][:1])
                    out <- Message{Type: MessageType(payload), Payload: payload, Padding: uint8(padding)}
                    count = 0
                    payload = ""
                }
            }
            out <- Message{Type: 255, Payload: "", Padding: 0}
        }
    }
}
