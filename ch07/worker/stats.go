package worker

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type CPUStat struct {
	Guest     uint64 `json:"guest"`
	GuestNice uint64 `json:"guest_nice"`
	Id        string `json:"id"`
	Idle      uint64 `json:"idle"`
	IOWait    uint64 `json:"iowait"`
	IRQ       uint64 `json:"irq"`
	Nice      uint64 `json:"nice"`
	SoftIRQ   uint64 `json:"softirq"`
	Steal     uint64 `json:"steal"`
	System    uint64 `json:"system"`
	User      uint64 `json:"user"`
}

type LoadAvg struct {
	Last15Min      float64 `json:"last15min"`
	Last1Min       float64 `json:"last1min"`
	Last5Min       float64 `json:"last5min"`
	LastPID        uint64  `json:"last_pid"`
	ProcessRunning uint64  `json:"process_running"`
	ProcessTotal   uint64  `json:"process_total"`
}

type DiskStat struct {
	All        uint64 `json:"all"`
	Free       uint64 `json:"free"`
	FreeInodes uint64 `json:"freeInodes"`
	Used       uint64 `json:"used"`
}

type MemInfo struct {
	Active            uint64 `json:"active"`
	AnonHugePages     uint64 `json:"anon_huge_pages"`
	AnonPages         uint64 `json:"anon_pages"`
	Bounce            uint64 `json:"bounce"`
	Buffers           uint64 `json:"buffers"`
	Cached            uint64 `json:"cached"`
	CommitLimit       uint64 `json:"commit_limit"`
	Committed_AS      uint64 `json:"committed_as"`
	DirectMap1G       uint64 `json:"direct_map_1G"`
	DirectMap2M       uint64 `json:"direct_map_2M"`
	DirectMap4k       uint64 `json:"direct_map_4k"`
	Dirty             uint64 `json:"dirty"`
	HardwareCorrupted uint64 `json:"hardware_corrupted"`
	HugePages_Free    uint64 `json:"huge_pages_free"`
	HugePages_Rsvd    uint64 `json:"huge_pages_rsvd"`
	HugePages_Surp    uint64 `json:"huge_pages_surp"`
	HugePages_Total   uint64 `json:"huge_pages_total"`
	Hugepagesize      uint64 `json:"hugepagesize"`
	Inactive          uint64 `json:"inactive"`
	KernelStack       uint64 `json:"kernel_stack"`
	Mapped            uint64 `json:"mapped"`
	MemAvailable      uint64 `json:"mem_available"`
	MemFree           uint64 `json:"mem_free"`
	MemTotal          uint64 `json:"mem_total"`
	NFS_Unstable      uint64 `json:"nfs_unstable"`
	PageTables        uint64 `json:"page_tables"`
	Shmem             uint64 `json:"shmem"`
	Slab              uint64 `json:"slab"`
	SReclaimable      uint64 `json:"s_reclaimable"`
	SUnreclaim        uint64 `json:"s_unclaim"`
	SwapCached        uint64 `json:"swap_cached"`
	SwapFree          uint64 `json:"swap_free"`
	SwapTotal         uint64 `json:"swap_total"`
	VmallocChunk      uint64 `json:"vmalloc_chunk"`
	VmallocTotal      uint64 `json:"vmalloc_total"`
	VmallocUsed       uint64 `json:"vmalloc_used"`
	Writeback         uint64 `json:"write_back"`
	WritebackTmp      uint64 `json:"writeback_tmp"`
}

type Stats struct {
	MemStats  *MemInfo
	DiskStats *DiskStat
	CpuStats  *CPUStat
	LoadStats *LoadAvg
	TaskCount int
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats:  GetCpuStats(),
		LoadStats: GetLoadAvg(),
	}
}

func GetMemoryInfo() *MemInfo {
	mem, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error reading memory stats")
		return &MemInfo{}
	}

	return &MemInfo{
		MemTotal:          mem.Total,
		MemFree:           mem.Free,
		MemAvailable:      mem.Available,
		Buffers:           mem.Buffers,
		Cached:            mem.Cached,
		SwapCached:        mem.SwapCached,
		Active:            mem.Active,
		Inactive:          mem.Inactive,
		SwapTotal:         mem.SwapTotal,
		SwapFree:          mem.SwapFree,
		Dirty:             mem.Dirty,
		Writeback:         mem.WriteBack,
		Mapped:            mem.Mapped,
		Shmem:             mem.Shared,
		Slab:              mem.Slab,
		SReclaimable:      mem.Sreclaimable,
		SUnreclaim:        mem.Sunreclaim,
		WritebackTmp:      mem.WriteBackTmp,
		CommitLimit:       mem.CommitLimit,
		Committed_AS:      mem.CommittedAS,
		VmallocTotal:      mem.VmallocTotal,
		VmallocUsed:       mem.VmallocUsed,
		VmallocChunk:      mem.VmallocChunk,
		AnonHugePages:     mem.AnonHugePages,
		HugePages_Total:   mem.HugePagesTotal,
		HugePages_Free:    mem.HugePagesFree,
		HugePages_Rsvd:    mem.HugePagesRsvd,
		HugePages_Surp:    mem.HugePagesSurp,
		Hugepagesize:      mem.HugePageSize,
	}
}

func GetDiskInfo() *DiskStat {
	usage, err := disk.Usage("/")
	if err != nil {
		log.Printf("Error reading volume stats /")
		return &DiskStat{}
	}

	return &DiskStat{
		All:        usage.Total,
		Used:       usage.Used,
		Free:       usage.Free,
		FreeInodes: usage.InodesFree,
	}
}

func GetCpuStats() *CPUStat {
	cpuTimes, _ := cpu.Times(false)
	t := cpuTimes[0]
	return &CPUStat{
		Id:        t.CPU,
		User:      uint64(t.User),
		Nice:      uint64(t.Nice),
		System:    uint64(t.System),
		Idle:      uint64(t.Idle),
		IOWait:    uint64(t.Iowait),
		IRQ:       uint64(t.Irq),
		SoftIRQ:   uint64(t.Softirq),
		Steal:     uint64(t.Steal),
		Guest:     uint64(t.Guest),
		GuestNice: uint64(t.GuestNice),
	}
}

func GetLoadAvg() *LoadAvg {
	loadavg, err := load.Avg()
	if err != nil {
		fmt.Println("Load average not supported on this platform.")
		return &LoadAvg{}
	}

	miscStat, err := load.Misc()
	if err != nil {
		fmt.Println("Load average not supported on this platform.")
		return &LoadAvg{}
	}

	return &LoadAvg{
		Last1Min:       loadavg.Load1,
		Last5Min:       loadavg.Load5,
		Last15Min:      loadavg.Load15,
		ProcessRunning: uint64(miscStat.ProcsRunning),
		ProcessTotal:   uint64(miscStat.ProcsTotal),
		LastPID:        uint64(miscStat.ProcsCreated),
	}
}
