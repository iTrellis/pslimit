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

import "time"

// Options is used to configure the memory watcher
type Options struct {
	// The amount of memory units that are required to trigger a Warning
	WarningLimit Unit

	// The amount of memory units that are required to trigger a Critical
	CriticalLimit Unit

	// Consecutive warnings that need to continuously meet, or a Critical is Triggered
	Cycle int

	// Memory check interval
	Interval time.Duration

	// Exit the Process
	Exit bool

	// Exit after ExitTime
	ExitTime time.Duration

	// Exit with ExitCode
	ExitCode int
}

func WarningLimit(u Unit) Option {
	return func(o *Options) {
		o.WarningLimit = u
	}
}

func CriticalLimit(u Unit) Option {
	return func(o *Options) {
		o.CriticalLimit = u
	}
}

func CycleLimit(limit int) Option {
	return func(o *Options) {
		o.Cycle = limit
	}
}

func Interval(interval time.Duration) Option {
	return func(o *Options) {
		o.Interval = interval
	}
}

func Exit(ps ...bool) Option {
	return func(o *Options) {
		if len(ps) > 0 {
			o.Exit = ps[0]
			return
		}
		o.Exit = true
	}
}

func ExitTime(t time.Duration) Option {
	return func(o *Options) {
		o.ExitTime = t
	}
}

func ExitCode(code int) Option {
	return func(o *Options) {
		o.ExitCode = code
	}
}
