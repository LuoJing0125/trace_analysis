package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

var Podmp map[string]TaskLog

func Err_Handle(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func ReadFile(file string) *PrometheusInfo {
	File, err := os.Open(file)

	res := new(PrometheusInfo)
	if err != nil {
		fmt.Println("File reading error", err)
		return res
	}
	prdec := json.NewDecoder(File)

	err = prdec.Decode(&res)
	if err != nil {
		fmt.Println(err)
	}
	File.Close()
	return res
}

func ReadPodContainerResourceLimits(filename string) error {
	var err error

	// kube_pod_container_resource_limits := ReadFile("./" + dir + "/kube_pod_container_resource_limits_2020-11-19.log") //kube_pod_container_resource_limits.log
	kube_pod_container_resource_limits := ReadFile("./resource_limit_pod/" + filename)

	for _, v := range kube_pod_container_resource_limits.Data.Result {
		tmptask, _ := Podmp[v.Metric.Pod]
		// if !ok {					//没有该pod，则跳过本次循环
		// 	continue
		// }

		tmptask.Pod = v.Metric.Pod
		value := reflect.ValueOf(v.RValue)
		if v.Metric.Resource == "memory" {
			tmpmem := tmptask.Memory
			tmpmem.Pod = tmptask.Pod
			tmpmem.Node = tmptask.Node

			tmpmem.Limit, err = strconv.ParseInt(value.Index(1).Elem().String(), 10, 64)
			if Err_Handle(err) {
				return err
			}

			tmptask.Memory = tmpmem

		} else if v.Metric.Resource == "nvidia_com_gpu" {
			tmpgpu := tmptask.GPU
			tmpgpu.Pod = tmptask.Pod
			tmpgpu.Node = tmptask.Node

			tmpgpu.NumGPU, err = strconv.ParseInt(value.Index(1).Elem().String(), 10, 64)
			if Err_Handle(err) {
				return err
			}

			tmptask.GPU = tmpgpu

		} else if v.Metric.Resource == "cpu" {
			tmpcpu := tmptask.CPU
			tmpcpu.Node = tmptask.Node
			tmpcpu.Pod = tmptask.Pod

			// tmpcpu.Limit, err = strconv.ParseInt(value.Index(1).Elem().String(), 10, 64)
			tmpcpu.Limit, err = strconv.ParseFloat(value.Index(1).Elem().String(), 64)
			if Err_Handle(err) {
				return err
			}

			tmptask.CPU = tmpcpu
		}
		// namepos := strings.Index(v.Metric.Container, "-")
		// tmptask.Container = v.Metric.Container
		// tmptask.User = v.Metric.Container[:namepos]//访问切片时，越界
		Podmp[v.Metric.Pod] = tmptask
	}

	return nil
}

func OuttoFile() {
	// nodeinfo, err := os.Create("./nodeinfo.csv")
	// if err != nil {
	// 	fmt.Println("node File creating error", err)
	// 	return
	// }
	// nodeinfo.WriteString("Podname\n")
	// for _, v := range Podmp {
	// 	nodeinfo.WriteString(v.Pod + "\n")
	// }
	// nodeinfo.Close()

	resourceinfo, err := os.Create("./resourceLimit.csv")
	if err != nil {
		fmt.Println("resource limits File creating error", err)
		return
	}

	resourceinfo.WriteString("Pod_Name,")
	// resourceinfo.WriteString("CPURequest,MemoryRequest,GPURequest\n")
	resourceinfo.WriteString("CPULimit,MemoryLimit,GPULimit\n")
	for _, v := range Podmp {
		resourceinfo.WriteString(v.Pod + ",")
		resourceinfo.WriteString(fmt.Sprintf("%.1f,%d,%d\n", v.CPU.Limit, v.Memory.Limit, v.GPU.NumGPU))
	}
	resourceinfo.Close()

}

func main() {

	Podmp = make(map[string]TaskLog)

	Lfiles, _ := ioutil.ReadDir("./resource_limit_pod") //读取资源限制log日期列表
	rl, err := os.Create("./Limit_date.log")            //创建文件
	if err != nil {
		fmt.Println("data File creating error", err)
		return
	}
	for _, f := range Lfiles {
		// fmt.Println(f.Name())
		rl.WriteString(f.Name() + "\n") //写入文件(字节数组)
	}
	rl.Close()

	var date []string
	resFile, err := os.Open("./Limit_date.log")
	defer resFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	rd := bufio.NewReader(resFile)

	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		date = append(date, line[:len(line)-1])

	}

	for _, d := range date {
		ReadPodContainerResourceLimits(d)
	}

	// ReadPodContainerResourceLimits("resource_limit_pod")

	fmt.Println()

	OuttoFile()

	return
}
