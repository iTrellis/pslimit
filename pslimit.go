/*
Copyright © 2021 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package pslimit

// Unit represents a size in bytes
type Unit int64

const (
	Byte     Unit = 1
	KiloByte      = 1024 * Byte
	MegaByte      = 1024 * KiloByte
	GigaByte      = 1024 * MegaByte
	TeraByte      = 1024 * GigaByte
)

// EventType represents a Watch Event Type
type EventType struct{}

// Boom Throws on overload
var Boom EventType

type Limit interface {
	Init(...Option)
	Options() Options
	TotalUnit() Unit
	ReachWarning(total Unit) bool
	ReachCritical(total Unit) bool
	Start() <-chan EventType
	Stop()
}

// Option option parameters
type Option func(*Options)
