package SrvNet

import (
	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const (
	ServerLoad_mem = "serverload_mem"
	ServerLoad_cpu = "serverload_cpu"
)

const (
	ServerLoad_memMax = 100
	ServerLoad_cpuMax = 100
)

const ()

//getMemoryUsedPercent 获取当前内存使用比例
func (n *_Node) getMemoryUsedPercent() float64 {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Error("mem.VirtualMemory error: ", err)
		return 0
	}

	//fmt.Printf("mem UsedPercent: %v%%\n", m.UsedPercent)
	return m.UsedPercent / 100
}

//getCpuUsedPercent 获取当前cpu使用比例
func (n *_Node) getCpuUsedPercent() float64 {
	c, err := cpu.Percent(0, false)
	if err != nil {
		log.Error("cpu.Percent error: ", err)
		return 0
	}

	//fmt.Printf("cpu UsedPercent_2_: %v%%\n", c[0])
	return c[0] / 100
}

//serverLoadloop
func (n *_Node) serverLoadloop() {
	//暂定五秒处理一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			n.defaultLoadData.Store(ServerLoad_cpu, n.getCpuUsedPercent())
			n.defaultLoadData.Store(ServerLoad_mem, n.getMemoryUsedPercent())
		}
	}
}

//getDefaultLoadData 默认上报数据接口
func (n *_Node) getDefaultLoadData() float32 {
	var result float32

	n.defaultLoadData.Range(func(key, value interface{}) bool {
		//keyString, ok := key.(string)
		//if !ok {
		//	return true
		//}

		valueFloat, ok := value.(float32)
		if !ok {
			return true
		}

		if valueFloat > result {
			result = valueFloat
		}
		return true
	})

	return result
}
