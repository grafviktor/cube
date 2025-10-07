package worker

import (
	"fmt"
	"log"

	"github.com/c9s/goprocinfo/linux"
	"github.com/shirou/gopsutil/load"
)

type CrossPlatformLoadAvg struct {
	Last1Min       float64 `json:"last1min"`
	Last5Min       float64 `json:"last5min"`
	Last15Min      float64 `json:"last15min"`
	ProcessRunning uint64  `json:"process_running"`
	ProcessTotal   uint64  `json:"process_total"`
	LastPID        uint64  `json:"last_pid"`
}

type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *CrossPlatformLoadAvg
	TaskCount int
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemAvailable
}

func (s *Stats) MemUsedPercent() uint64 {
	return s.MemStats.MemAvailable / s.MemStats.MemTotal
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemAvailable
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) CpuUsage() float64 {

	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle

	if total == 0 {
		return 0.00
	}

	return (float64(total) - float64(idle)) / float64(total)
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats:  GetCpuStats(),
		LoadStats: GetLoadAvg(),
	}
}

// GetMemoryInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#MemInfo
func GetMemoryInfo() *linux.MemInfo {
	memstats, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("Error reading from /proc/meminfo")
		return &linux.MemInfo{}
	}

	return memstats
}

// GetDiskInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#Disk
func GetDiskInfo() *linux.Disk {
	diskstats, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("Error reading from /")
		return &linux.Disk{}
	}

	return diskstats
}

// GetCpuInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#CPUStat
func GetCpuStats() *linux.CPUStat {
	stats, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error reading from /proc/stat")
		return &linux.CPUStat{}
	}

	return &stats.CPUStatAll
}

// GetLoadAvg See https://godoc.org/github.com/c9s/goprocinfo/linux#LoadAvg
// func GetLoadAvg() *linux.LoadAvg {
// 	loadavg, err := linux.ReadLoadAvg("/proc/loadavg")
// 	if err != nil {
// 		log.Printf("Error reading from /proc/loadavg")
// 		return &linux.LoadAvg{}
// 	}

// 	return loadavg
// }

func GetLoadAvg() *CrossPlatformLoadAvg {
	loadavg, err := load.Avg()
	if err != nil {
		fmt.Println("Load average not supported on this platform.")
		return &CrossPlatformLoadAvg{}
	}

	miscStat, err := load.Misc()
	if err != nil {
		fmt.Println("Load average not supported on this platform.")
		return &CrossPlatformLoadAvg{}
	}

	// loadavg.Load1, loadavg.Load5, loadavg.Load15

	return &CrossPlatformLoadAvg{
		Last1Min:       loadavg.Load1,
		Last5Min:       loadavg.Load5,
		Last15Min:      loadavg.Load15,
		ProcessRunning: uint64(miscStat.ProcsRunning),
		ProcessTotal:   uint64(miscStat.ProcsTotal),
		LastPID:        uint64(miscStat.ProcsCreated),
	}
}

/*
package main

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	// --- MemInfo ---
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Memory Info ===")
	fmt.Printf("Total: %v MB\n", vm.Total/1024/1024)
	fmt.Printf("Used:  %v MB\n", vm.Used/1024/1024)
	fmt.Printf("Free:  %v MB\n", vm.Free/1024/1024)
	fmt.Printf("Used Percent: %.2f%%\n\n", vm.UsedPercent)

	// --- Disk ---
	fmt.Println("=== Disk Info ===")
	partitions, _ := disk.Partitions(false)
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err == nil {
			fmt.Printf("%s -> Total: %v GB, Used: %v GB (%.2f%%)\n",
				p.Mountpoint,
				usage.Total/1024/1024/1024,
				usage.Used/1024/1024/1024,
				usage.UsedPercent)
		}
	}
	fmt.Println()

	// --- CPU Stat ---
	fmt.Println("=== CPU Info ===")
	cpuTimes, _ := cpu.Times(false)
	for _, t := range cpuTimes {
		fmt.Printf("User: %.2f  System: %.2f  Idle: %.2f\n",
			t.User, t.System, t.Idle)
	}
	percent, _ := cpu.Percent(0, false)
	fmt.Printf("CPU Usage: %.2f%%\n\n", percent[0])

	// --- LoadAvg ---
	fmt.Println("=== Load Average ===")
	loadavg, err := load.Avg()
	if err == nil {
		fmt.Printf("1-min: %.2f, 5-min: %.2f, 15-min: %.2f\n",
			loadavg.Load1, loadavg.Load5, loadavg.Load15)
	} else {
		fmt.Println("Load average not supported on this platform.")
	}
}
*/
