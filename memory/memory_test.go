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

package memory_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/iTrellis/pslimit"
	"github.com/iTrellis/pslimit/memory"
)

func TestCalculateTotalMemory(t *testing.T) {
	var mstats runtime.MemStats
	runtime.ReadMemStats(&mstats)

	expectedTotal := pslimit.Unit(mstats.HeapInuse + mstats.StackInuse +
		mstats.MSpanInuse + mstats.MCacheInuse + mstats.BuckHashSys)

	totalmem := memory.CalculateTotalMemory(mstats)

	if totalmem != expectedTotal {
		t.Fatal("Expected total memory to match, but it didn't")
	}
}

func TestMemoryWatcher(t *testing.T) {
	t.Run("new", func(t *testing.T) {

		defOptions := pslimit.Options{
			WarningLimit:  512 * pslimit.MegaByte,
			CriticalLimit: 768 * pslimit.MegaByte,
			Cycle:         10,
			Interval:      5 * time.Second,
			ExitCode:      119,
			Exit:          true,
			ExitTime:      10 * time.Second,
		}
		t.Log("When no option is provided")
		{
			mLimit := memory.New()
			if mLimit.Options() != defOptions {
				t.Fatal("Expected default options")
			}
		}

		t.Log("When complete new options is provided")
		{
			wcfgCustom := pslimit.Options{
				WarningLimit:  256 * pslimit.MegaByte,
				CriticalLimit: 512 * pslimit.MegaByte,
				Cycle:         5,
				Interval:      15 * time.Second,
				ExitCode:      102,
				ExitTime:      20 * time.Second,
			}

			mLimit := memory.New()

			mLimit.Init(
				pslimit.WarningLimit(256*pslimit.MegaByte),
				pslimit.CriticalLimit(512*pslimit.MegaByte),
				pslimit.CycleLimit(5),
				pslimit.Interval(15*time.Second),
				pslimit.ExitCode(102),
				pslimit.ExitTime(20*time.Second),
				pslimit.Exit(false),
			)
			if mLimit.Options() != wcfgCustom {
				t.Fatal("Expected custom options")
			}
		}

		t.Log("When partial options is provided")
		{
			wcfgExpected1 := pslimit.Options{
				WarningLimit:  256 * pslimit.MegaByte,
				CriticalLimit: 512 * pslimit.MegaByte,
				Cycle:         defOptions.Cycle,
				Interval:      defOptions.Interval,
				ExitCode:      defOptions.ExitCode,
				Exit:          defOptions.Exit,
				ExitTime:      defOptions.ExitTime,
			}

			wcfgExpected2 := pslimit.Options{
				WarningLimit:  defOptions.WarningLimit,
				CriticalLimit: defOptions.CriticalLimit,
				Cycle:         15,
				Interval:      15 * time.Second,
				ExitCode:      102,
				Exit:          true,
				ExitTime:      20 * time.Second,
			}

			mLimit := memory.New()

			mLimit.Init(
				pslimit.WarningLimit(256*pslimit.MegaByte),
				pslimit.CriticalLimit(512*pslimit.MegaByte),
				pslimit.Exit(),
			)
			if mLimit.Options() != wcfgExpected1 {
				t.Fatal("Expected first merged option to match, but it didn't")
			}

			mLimit = memory.New()
			mLimit.Init(
				pslimit.CycleLimit(15),
				pslimit.Interval(15*time.Second),
				pslimit.ExitCode(102),
				pslimit.ExitTime(20*time.Second),
			)

			if mLimit.Options() != wcfgExpected2 {
				t.Fatal("Expected second merged options to match, but it didn't")
			}
		}
	})

	t.Run("reach_critical", func(t *testing.T) {

		mLimit := memory.New()

		mLimit.Init(
			pslimit.CriticalLimit(512 * pslimit.MegaByte),
		)

		notReachedTotal := 512 * pslimit.MegaByte
		if mLimit.ReachCritical(notReachedTotal) {
			t.Fatal("Expected total NOT to reach critical, but it did")
		}

		reachedTotal := 513 * pslimit.MegaByte
		if !mLimit.ReachCritical(reachedTotal) {
			t.Fatal("Expected total to reach critical, but it didn't")
		}
	})

	t.Run("reach_warning", func(t *testing.T) {
		mLimit := memory.New()

		mLimit.Init(
			pslimit.WarningLimit(256 * pslimit.MegaByte),
		)

		notReachedTotal := 256 * pslimit.MegaByte
		if mLimit.ReachWarning(notReachedTotal) {
			t.Fatal("Expected total NOT to reach warning, but it did")
		}

		reachedTotal := 257 * pslimit.MegaByte
		if !mLimit.ReachWarning(reachedTotal) {
			t.Fatal("Expected total to reach warning, but it didn't")
		}
	})
}
