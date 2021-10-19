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
	"strings"
)

var Podmp map[string]TaskLog
var Mypod map[string]MyLog

func Err_Handle(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func ReadFile(file string) *ResultInfo {
	File, err := os.Open(file)

	res := new(ResultInfo)
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

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func Cal_Average_int(src []int64) float64 {
	var sum int64 = 0
	if len(src) == 0 {
		return 0
	}
	for _, v := range src {
		sum = sum + v
	}

	return Decimal(float64(sum) / float64(len(src)))
}

func Cal_Average_float(src []float64) float64 {
	var sum float64 = 0
	if len(src) == 0 {
		return 0
	}
	for _, v := range src {
		sum = sum + v
	}

	return Decimal(float64(sum) / float64(len(src)))
}

func MAX(x, y interface{}) interface{} {
	if reflect.TypeOf(x).Name() == "int64" {
		if reflect.ValueOf(x).Int() > reflect.ValueOf(y).Int() {
			return x
		}
		return y
	} else if reflect.TypeOf(x).Name() == "float64" {
		if reflect.ValueOf(x).Float() > reflect.ValueOf(y).Float() {
			return x
		}
		return y
	}
	return x
}

//fb_used
func Readgpumem(path string) {
	File, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer File.Close()

	br := bufio.NewReader(File)

	for {
		line, err := br.ReadString('\n') //读行
		if err != nil || io.EOF == err {
			break
		}

		if line[0] == '[' {
			line = line[1:] //去第一行开头的[
		}

		if line[0] != '{' {
			continue
		}

		line = line[:len(line)-2] //去每行结尾的,或]

		res := new(ResultInfo)
		prdec := json.NewDecoder(strings.NewReader(line))
		err = prdec.Decode(&res)
		if err != nil {
			fmt.Println(err)

		}

		if res.Metric.Pod_name == "" || res.Metric.Pod_namespace == "" || res.Metric.Pod_namespace == "default" || res.Metric.Pod_namespace == "ingress-nginx" || res.Metric.Pod_namespace == "kube-system" || res.Metric.Pod_namespace == "lens-metrics" {
			continue
		}

		val, _ := Podmp[res.Metric.Pod_name]

		val.Pod_name = res.Metric.Pod_name

		for _, v := range res.RDetail {
			value := reflect.ValueOf(v.RValue)

			flag := false
			for k, _ := range val.GPU.GPUMem {
				if val.GPU.GPUMem[k].Pod == res.Metric.Pod_name {
					flag = true
					for i := 0; i < value.Len(); i++ {
						gpumem, _ := strconv.ParseInt(value.Index(i).Elem().String(), 10, 64)
						val.GPU.GPUMem[k].MaxR = MAX(val.GPU.GPUMem[k].MaxR, gpumem).(int64) //最大值
						val.GPU.GPUMem[k].History = append(val.GPU.GPUMem[k].History, gpumem)
					}
					//平均值
					var aver float64 = 0
					for _, vv := range val.GPU.GPUMem {
						aver += Cal_Average_int(vv.History)
					}
					aver = Decimal(aver / float64(len(val.GPU.GPUMem)))
					val.GPU.GPUMem[k].AveR = int64(aver)

					Podmp[res.Metric.Pod_name] = val
				}
			}
			if flag {
				continue
			}

			var tmpgmemhis GPUMemHistory
			tmpgmemhis.Uuid = res.Metric.Uuid
			tmpgmemhis.Pod = res.Metric.Pod_name
			tmpgmemhis.Total = 0
			tmpgmemhis.MaxR = 0
			tmpgmemhis.AveR = 0

			var tmpMax int64 = 0

			for i := 0; i < value.Len(); i++ {
				gpumem, _ := strconv.ParseInt(value.Index(i).Elem().String(), 10, 64)
				tmpMax = MAX(tmpMax, gpumem).(int64)
				tmpgmemhis.History = append(tmpgmemhis.History, gpumem)
			}
			tmpgmemhis.MaxR = tmpMax

			tmpgmemhis.AveR = int64(Cal_Average_int(tmpgmemhis.History))

			val.GPU.GPUMem = append(val.GPU.GPUMem, tmpgmemhis)
			Podmp[res.Metric.Pod_name] = val
		}
	}
}

//power
func Readgpupower(path string) {
	File, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer File.Close()

	br := bufio.NewReader(File)

	for {
		line, err := br.ReadString('\n') //读行
		if err != nil || io.EOF == err {
			break
		}

		if line[0] == '[' {
			line = line[1:] //去第一行开头的[
		}

		if line[0] != '{' {
			continue
		}

		line = line[:len(line)-2] //去每行结尾的,或]

		res := new(ResultInfo)
		prdec := json.NewDecoder(strings.NewReader(line))
		err = prdec.Decode(&res)
		if err != nil {
			fmt.Println(err)

		}

		if res.Metric.Pod_name == "" || res.Metric.Pod_namespace == "" || res.Metric.Pod_namespace == "default" || res.Metric.Pod_namespace == "ingress-nginx" || res.Metric.Pod_namespace == "kube-system" || res.Metric.Pod_namespace == "lens-metrics" {
			continue
		}

		val, _ := Podmp[res.Metric.Pod_name]

		val.Pod_name = res.Metric.Pod_name

		for _, v := range res.RDetail {
			value := reflect.ValueOf(v.RValue)

			flag := false
			for k, _ := range val.GPU.GPUPower {
				if val.GPU.GPUPower[k].Pod == res.Metric.Pod_name {
					flag = true
					for i := 0; i < value.Len(); i++ {
						gpupower, _ := strconv.ParseFloat(value.Index(i).Elem().String(), 64)
						val.GPU.GPUPower[k].MaxR = MAX(val.GPU.GPUPower[k].MaxR, gpupower).(float64) //最大值
						val.GPU.GPUPower[k].History = append(val.GPU.GPUPower[k].History, gpupower)
					}
					//平均值
					var aver float64 = 0
					for _, vv := range val.GPU.GPUPower {
						aver += Cal_Average_float(vv.History)
					}
					aver = Decimal(aver / float64(len(val.GPU.GPUPower)))
					val.GPU.GPUPower[k].AveR = aver

					Podmp[res.Metric.Pod_name] = val
				}
			}
			if flag {
				continue
			}

			var tmpgpowhis GPUPowHistory
			tmpgpowhis.Uuid = res.Metric.Uuid
			tmpgpowhis.Pod = res.Metric.Pod_name
			tmpgpowhis.MaxR = 0
			tmpgpowhis.AveR = 0

			var tmpMax float64 = 0

			for i := 0; i < value.Len(); i++ {
				gpupower, _ := strconv.ParseFloat(value.Index(i).Elem().String(), 64)
				tmpMax = MAX(tmpMax, gpupower).(float64)
				tmpgpowhis.History = append(tmpgpowhis.History, gpupower)
			}
			tmpgpowhis.MaxR = tmpMax

			tmpgpowhis.AveR = Cal_Average_float(tmpgpowhis.History)

			val.GPU.GPUPower = append(val.GPU.GPUPower, tmpgpowhis)
			Podmp[res.Metric.Pod_name] = val
		}
	}
}

//gpu_utilization
func Readgpuutil(path string) {
	File, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer File.Close()

	br := bufio.NewReader(File)

	for {
		line, err := br.ReadString('\n') //读行
		if err != nil || io.EOF == err {
			break
		}

		if line[0] == '[' {
			line = line[1:] //去第一行开头的[
		}

		if line[0] != '{' {
			continue
		}

		line = line[:len(line)-2] //去每行结尾的,或]

		res := new(ResultInfo)
		prdec := json.NewDecoder(strings.NewReader(line))
		err = prdec.Decode(&res)
		if err != nil {
			fmt.Println(err)

		}

		if res.Metric.Pod_name == "" || res.Metric.Pod_namespace == "" || res.Metric.Pod_namespace == "default" || res.Metric.Pod_namespace == "ingress-nginx" || res.Metric.Pod_namespace == "kube-system" || res.Metric.Pod_namespace == "lens-metrics" {
			continue
		}

		val, _ := Podmp[res.Metric.Pod_name]

		val.Pod_name = res.Metric.Pod_name

		for _, v := range res.RDetail {
			value := reflect.ValueOf(v.RValue)

			flag := false
			for k, _ := range val.GPU.GPUUtil {
				if val.GPU.GPUUtil[k].Pod == res.Metric.Pod_name {
					flag = true
					for i := 0; i < value.Len(); i++ {
						gpuutil, _ := strconv.ParseInt(value.Index(i).Elem().String(), 10, 64)
						val.GPU.GPUUtil[k].MaxR = MAX(val.GPU.GPUUtil[k].MaxR, gpuutil).(int64) //最大值
						val.GPU.GPUUtil[k].History = append(val.GPU.GPUUtil[k].History, gpuutil)
					}
					//平均值
					var aver float64 = 0
					for _, vv := range val.GPU.GPUUtil {
						aver += Cal_Average_int(vv.History)
					}
					aver = Decimal(aver / float64(len(val.GPU.GPUUtil)))
					val.GPU.GPUUtil[k].AveR = aver

					Podmp[res.Metric.Pod_name] = val
				}
			}
			if flag {
				continue
			}

			var tmpgutilhis GPUHistory
			tmpgutilhis.Uuid = res.Metric.Uuid
			tmpgutilhis.Pod = res.Metric.Pod_name
			tmpgutilhis.MaxR = 0
			tmpgutilhis.AveR = 0

			var tmpMax int64 = 0

			for i := 0; i < value.Len(); i++ {
				gpuutil, _ := strconv.ParseInt(value.Index(i).Elem().String(), 10, 64)
				tmpMax = MAX(tmpMax, gpuutil).(int64)
				tmpgutilhis.History = append(tmpgutilhis.History, gpuutil)
			}
			tmpgutilhis.MaxR = tmpMax

			tmpgutilhis.AveR = Cal_Average_int(tmpgutilhis.History)

			val.GPU.GPUUtil = append(val.GPU.GPUUtil, tmpgutilhis)
			Podmp[res.Metric.Pod_name] = val
		}
	}
}

func OuttoFile() {
	gpuinfo, err := os.Create("./GPUdata.csv")
	if err != nil {
		fmt.Println("File creating error", err)
		return
	}

	//转移到Mypod
	for _, v := range Podmp {
		//GPU显存
		for _, vv := range v.GPU.GPUMem {
			myp, _ := Mypod[v.Pod_name]
			myp.Pod_name = v.Pod_name
			myp.GPUMemAve = vv.AveR
			myp.GPUMemMax = vv.MaxR
			Mypod[v.Pod_name] = myp
		}
		//GPU利用率
		for _, vv := range v.GPU.GPUUtil {
			myp, _ := Mypod[v.Pod_name]
			myp.Pod_name = v.Pod_name
			myp.GPUUtilAve = vv.AveR
			// myp.GPUUtilMax = vv.MaxR
			Mypod[v.Pod_name] = myp
		}
		//GPU功率
		for _, vv := range v.GPU.GPUPower {
			myp, _ := Mypod[v.Pod_name]
			myp.Pod_name = v.Pod_name
			myp.GPUPowerAve = vv.AveR
			// myp.GPUPowerMax = vv.MaxR
			Mypod[v.Pod_name] = myp
		}
	}

	// gpuinfo.WriteString("Pod_Name,GPUMemAve,GPUMemMax\n")
	gpuinfo.WriteString("Pod_Name,GPUMemAve,GPUMemMax,GPUUtilAve,GPUPowerAve\n")

	for _, v := range Mypod {
		gpuinfo.WriteString(v.Pod_name + ",")
		gpuinfo.WriteString(fmt.Sprintf("%d,%d,", v.GPUMemAve, v.GPUMemMax))
		gpuinfo.WriteString(fmt.Sprintf("%.3f,%.3f\n", v.GPUUtilAve, v.GPUPowerAve))
	}

	gpuinfo.Close()
}

func main() {

	Podmp = make(map[string]TaskLog)
	Mypod = make(map[string]MyLog)

	// Mfiles, _ := ioutil.ReadDir("./fb_used")
	// Ufiles, _ := ioutil.ReadDir("./gpu_utilization")
	// Pfiles, _ := ioutil.ReadDir("./power")
	Folders, _ := ioutil.ReadDir("./data")
	for _, fo := range Folders {
		if fo.Name() == "fb_used" {
			Files, _ := ioutil.ReadDir("./data/" + fo.Name())
			for _, fi := range Files {
				Logs, _ := ioutil.ReadDir("./data/" + fo.Name() + "/" + fi.Name())
				for _, l := range Logs {
					Readgpumem("./data/" + fo.Name() + "/" + fi.Name() + "/" + l.Name())
				}
			}
		}
		if fo.Name() == "gpu_utilization" {
			Files, _ := ioutil.ReadDir("./data/" + fo.Name())
			for _, fi := range Files {
				Logs, _ := ioutil.ReadDir("./data/" + fo.Name() + "/" + fi.Name())
				for _, l := range Logs {
					Readgpuutil("./data/" + fo.Name() + "/" + fi.Name() + "/" + l.Name())
				}
			}
		}
		if fo.Name() == "power" {
			Files, _ := ioutil.ReadDir("./data/" + fo.Name())
			for _, fi := range Files {
				Logs, _ := ioutil.ReadDir("./data/" + fo.Name() + "/" + fi.Name())
				for _, l := range Logs {
					Readgpupower("./data/" + fo.Name() + "/" + fi.Name() + "/" + l.Name())
				}
			}
		}

	}

	fmt.Println()

	OuttoFile()

	return
}
