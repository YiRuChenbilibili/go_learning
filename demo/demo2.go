//获取cpu/内存信息

package main

import (
	f "fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	// "github.com/shirou/gopsutil/mem"  // to use v2
)

func main() {
	var input int
	for {
		f.Println("输入编码获取信息: 0:Percentcpu;1:Numcpu;2:MemTotal;3:MemUsed;4:MemFree;5:MemUsedPercet;6:getMemall;7:getCpuall;8:quit")
		f.Scanf("%d\n", &input)
		// almost every return value is a struct
		switch input {
		case 0:
			var cpu Cpuinfo
			f.Printf("Percentcpu:%f%%\n", cpu.getPercentcpu())
		case 1:
			var cpu Cpuinfo
			f.Printf("CpuNum:%v\n", cpu.getCpuNum())
		case 2:
			f.Printf("MemTotal:%v\n", getMemTotal())
		case 3:
			f.Printf("MemUsed:%v\n", getMemUsed())
		case 4:
			f.Printf("MemFree:%v\n", getMemFree())
		case 5:
			f.Printf("UsedPercent:%f%%\n", getMemUsedPercent())
		case 6:
			v := getMemAll()
			f.Printf("Total: %v, Free:%v,Used:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.Used, v.UsedPercent)
		case 7:
			var cpu Cpuinfo
			f.Printf("CpuNum: %v, Percentcpu:%v\n", cpu.getCpuNum(), cpu.getPercentcpu())

		}

		// convert to JSON. String() is also implemented
		//f.Println(v)
		if input == 8 {
			break
		}
	}
}

//cpu信息的结构体
type Cpuinfo struct {
	Percentcpu float64
	Numcpu     int
}

// type Meminfo struct {
// 	Total      uint64
// 	Used       uint64
// 	Free       uint64
// 	UsedPercet uint64
// }

//获得cpu占用率的方法
func (c *Cpuinfo) getPercentcpu() float64 {
	Percent, _ := cpu.Percent(time.Second, false)
	c.Percentcpu = Percent[0]
	return c.Percentcpu
}

//获得核数的方法
func (c *Cpuinfo) getCpuNum() int {
	c.Numcpu = runtime.NumCPU()
	return c.Numcpu
}

/*
在github.com/shirou/gopsutil中通过使用
Modkernel32 = windows.NewLazySystemDLL("kernel32.dll"）,
procGlobalMemoryStatusEx = common.Modkernel32.NewProc("GlobalMemoryStatusEx")及
mem, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
在windows上实现mem.VirtualMemory()，使用"golang.org/x/sys/windows"实现
*/
//获得内存相关内容的函数
func getMemTotal() uint64 {
	Meminfo, _ := mem.VirtualMemory()
	return Meminfo.Total
}

func getMemUsed() uint64 {
	Meminfo, _ := mem.VirtualMemory()
	return Meminfo.Used
}

func getMemFree() uint64 {
	Meminfo, _ := mem.VirtualMemory()
	return Meminfo.Free
}

func getMemUsedPercent() float64 {
	Meminfo, _ := mem.VirtualMemory()
	return Meminfo.UsedPercent
}

//获得内存所有信息的函数
func getMemAll() *mem.VirtualMemoryStat {
	Meminfo, _ := mem.VirtualMemory()
	return Meminfo
}
