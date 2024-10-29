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

// A Message stores the important properties of a AIS message, including only information useful
// for decoding: Type, Payload, Padding Bits
// A Message should come after processing one or more AIS radio sentences (checksum check,
// concatenate payloads spanning across sentences, etc).
type Message struct {
    Type    uint8
    Padding uint8
    Payload string
}
