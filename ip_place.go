package ip_place

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type IPRange struct {
	End          int64
	CityId       int64
	CityName     string
	ProvinceName string
	Mark         int
}

var IP_COUNT = 354861
var IPMap = make(map[int64]*IPRange, IP_COUNT)
var IPSlice = make([]int64, IP_COUNT)

/**
查找CityId
*/
func init() {
	err := LoadIPData()
	if err != nil {
		fmt.Println(err)
	}
}

func GetPlaceNameByIP(ipstr string) (string, string, error) {
	ip := IpToLong(ipstr)
	index := SearchStartIndex(IP_COUNT, func(i int) bool { return IPSlice[i] <= ip })
	startIp := IPSlice[index]
	tail := IPMap[startIp]

	if tail.Mark == 1 {
		return tail.ProvinceName, "no city", nil
	} else if tail.Mark == 2 {
		return "this ip is not in our library", "", nil
	} else {
		return tail.ProvinceName, tail.CityName, nil
	}
}

func SearchStartIndex(n int, f func(int) bool) int {
	i, j := 0, n
	var k, h int = 0, 0
	for i < j {
		k = h
		h = i + (j-i)/2
		if !f(h) {
			j = h - 1
		} else {
			i = h
		}
		if k == h {
			i = j
		}
	}
	return i
}

//加载IP数据到内存
func LoadIPData() error {
	return ReadStringLine("IP_Library.txt", processIPLine)
}
func ReadStringLine(filePth string, hookfn func(string)) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadString('\n')
		hookfn(line)
		if err != nil {
			if err == io.EOF{
				return nil
			}
			return err
		}
	}
}

var lineIndex = 0

func processIPLine(line string) {
	arr := strings.Fields(line)
	if len(arr) != 4 {
		panic(fmt.Sprintf("IP库错误,index:%d,line:%s", lineIndex, line))
	}
	ip_start, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("解析StartIP错误,index:%d,line:%s", lineIndex, line))
	}
	ip_end, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("解析EndIP错误,index:%d,line:%s", lineIndex, line))
	}
	var cityId int64
	if strings.EqualFold(arr[3], "-") {
		if strings.EqualFold(arr[2], "-") {
			cityId = 0
			IPMap[ip_start] = &IPRange{End: ip_end, CityId: cityId, Mark: 2}
		} else {
			cityId = GetCodeByAreaName(arr[2])
			provincename := arr[2]
			IPMap[ip_start] = &IPRange{End: ip_end, CityId: cityId, ProvinceName: provincename, Mark: 1}
		}
	} else {
		cityId := GetCodeByAreaName(arr[3])
		cityname := arr[3]
		provincename := arr[2]
		IPMap[ip_start] = &IPRange{End: ip_end, CityId: cityId, CityName: cityname, ProvinceName: provincename, Mark: 0}
	}

	IPSlice[lineIndex] = ip_start
	lineIndex++
}

func IpToLong(ip string) int64 {
	bits := strings.Split(ip, ".")
	if len(bits) != 4 {
		return 0
	}
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

func LongToIp(ipnr int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}
