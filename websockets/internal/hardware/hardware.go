package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetMemory() (string, error) {
	runtime := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Hostname: %s\nTotal Memory: %d\nUsed Memory: %d\nOS: %s", hostStat.Hostname, vmStat.Total, vmStat.Used, runtime)

	return output, nil

}

func GetCpu() (string, error) {
	cpuStat, err := cpu.Info()

	if err != nil {
		fmt.Println(err)
	}

	output := fmt.Sprintf("Cores: %d\nCPU model: %s\nmhz: %f", cpuStat[0].Cores, cpuStat[0].ModelName, cpuStat[0].Mhz)

	return output, err
}

func GetDisk() (string, error) {
	diskStat, err := disk.Usage("/")

	if err != nil {
		fmt.Println(err)
	}

	output := fmt.Sprintf("Total disk usage: %d\nFree disk space: %d", diskStat.Total, diskStat.Free)

	return output, err
}
