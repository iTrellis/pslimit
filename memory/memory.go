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

package memory

import (
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/iTrellis/pslimit"
)

var defaultOptions = pslimit.Options{
	WarningLimit:  512 * pslimit.MegaByte,
	CriticalLimit: 768 * pslimit.MegaByte,
	Cycle:         10,
	Interval:      5 * time.Second,
	Exit:          true,
	ExitCode:      119,
	ExitTime:      10 * time.Second,
}

type memLimit struct {
	options pslimit.Options
	events  chan pslimit.EventType
	once    sync.Once
	count   int
	ticker  *time.Ticker
	stopped chan struct{}
}

// New 生成对象
func New() pslimit.Limit {
	p := &memLimit{
		options: defaultOptions,
		events:  make(chan pslimit.EventType),
		stopped: make(chan struct{}),
	}
	return p
}

func (p *memLimit) Boom() {
	p.events <- pslimit.Boom
	close(p.events)
}

func (p *memLimit) Init(opts ...pslimit.Option) {
	for _, o := range opts {
		o(&p.options)
	}
}
func (p *memLimit) Options() pslimit.Options {
	return p.options
}

func (p *memLimit) Start() <-chan pslimit.EventType {
	go func() {
		p.ticker = time.NewTicker(p.options.Interval)
		for {
			select {
			case <-p.ticker.C:
				p.tick()
			case <-p.stopped:
				return
			}
		}
	}()
	return p.events
}

func (p *memLimit) Stop() {
	p.count = 0
	p.ticker.Stop()
}

func (p *memLimit) TotalUnit() pslimit.Unit {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	return CalculateTotalMemory(rtm)
}

// CalculateTotalMemory returns the total memory usage
func CalculateTotalMemory(stats runtime.MemStats) pslimit.Unit {
	// Sys related stats may be released to the OS so the runtime
	// memory usage would not be close with the one observed via
	// ps or activity monitor
	return pslimit.Unit(stats.HeapInuse +
		stats.StackInuse +
		stats.MSpanInuse +
		stats.MCacheInuse +
		stats.BuckHashSys)
}

// ReachCritical returns whether total memory reached critical threshold
func (p *memLimit) ReachCritical(total pslimit.Unit) bool {
	return total > p.options.CriticalLimit
}

// ReachWarning returns whether total memory reached warning threshold
func (p *memLimit) ReachWarning(total pslimit.Unit) bool {
	return total > p.options.WarningLimit
}

// Call in every circle, check the memory usage
func (p *memLimit) tick() {
	total := p.TotalUnit()
	if p.ReachCritical(total) {
		p.once.Do(p.trigger)
	} else if p.ReachWarning(total) {
		p.count++
		if p.count >= p.options.Cycle {
			p.once.Do(p.trigger)
		}
	} else {
		p.count = 0
	}
}

// Trigger boom and exit after the ExitTime duration
func (p *memLimit) trigger() {
	p.Boom()
	p.Stop()
	p.stopped <- struct{}{}

	<-time.After(p.options.ExitTime)
	if p.options.Exit {
		os.Exit(p.options.ExitCode)
	}
}
