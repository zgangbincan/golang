package main

import (
	"fmt"
	"strconv"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"sort"
	"math/rand"
)

// 定义虚拟节点的全局变量，每个target复制64个虚拟节点
var Replicas int = 64

// 目标节点的初始化个数
var TargetCount int = 0

// 位置到target的映射
var PositionToTarget = map[string]string{}

// target到position的映射
var TargetToPosition = map[string][64]string{}

// 1代表没有排序过，2表示排序过
var PositionToTargetSorted int = 1

// 搞一个数组，这个用于positions的排序
var PositionsRank = []string{}

// 做的一个函数
type Mymethod int

// 添加节点
func (this *Mymethod)AddTarget(Target string) bool{
		// 首先得判断这个节点是否存在
		_, ok := TargetToPosition[Target]
		if ok {
				fmt.Printf("target:%s already exist.", Target)
				return false
			}
		// 创建hash
		h := md5.New()
		// 然后添加节点
		var Positions [64]string
		for I := 0; I < Replicas; I++{
				// 首先转化为string
				StrI := strconv.Itoa(I)
				// 把值转化为最后可以hash的值
				Prehash := []string{Target, StrI}
				HashStr := strings.Join(Prehash, "")
				// 对值进行md5
				h.Write([]byte(HashStr))
				Position  := hex.EncodeToString(h.Sum(nil))
				// 赋值positiontotarget的对应关系
				PositionToTarget[Position] = Target
				Positions[I] = Position
			}
		TargetToPosition[Target] = Positions
		TargetCount = TargetCount + 1
		// 新增加节点之后需要重新排序，为1表示需要排序。
		PositionToTargetSorted = 1
		return true
	}

// 删除节点
func (this *Mymethod)RemoveTarget(Target string) bool{
		// 先判断是否有这个Target
		_, ok := TargetToPosition[Target]
		if !ok {
			fmt.Printf("target:%s not exist.", Target)
			return false
		}
		// 删除PositionToTarget里面的position
		for _, position := range TargetToPosition[Target]{
				delete(PositionToTarget, position)
			}
		// 删除TargetToPosition里面的Target
		delete(TargetToPosition, Target)
		TargetCount = TargetCount - 1
		// 删除节点，也是需要排序的
		PositionToTargetSorted = 1
		return true
	}

// 查找所有的节点
func (this *Mymethod)LookUp(Resource string)(string){
		// 寻找target,目前只是一个，如果需要多个，那么只需更改第二个参数即可
		Results := this.lookUplist(Resource, 1)
		return Results
	}

// 具体的查找节点的函数
func (this *Mymethod)lookUplist(Resource string, RequestCount int)(string){
		var Results string
		if RequestCount == 0{
				fmt.Printf("RequestCount is:%d", RequestCount)
				return Results
			}
		// 当只有一个target的时候，直接返回
		if TargetCount == 1{
				for Target, _ := range TargetToPosition{
						fmt.Printf("Target is :% + v", Target)
						Results = Target
						break
					}
					return Results
			}
		this.rankPositions()
		// 对resource 进行md5,然后进行查找
		hLook := md5.New()
		hLook.Write([]byte(Resource))
		// 查找资源的hash值
		resourceHash := hex.EncodeToString(hLook.Sum(nil))
		// 查看是否已经有这个节点
		target, ok := PositionToTarget[resourceHash]
		if ok {
				Results = target
				return Results
			}
		positionLen := len(PositionsRank)
		// 看自己的排名
		sourceRank := sort.SearchStrings(PositionsRank,resourceHash)
		// 返回前一个节点
		if sourceRank == 0{
				Results = PositionToTarget[PositionsRank[positionLen - 1]]
			}else{
					Results = PositionToTarget[PositionsRank[sourceRank - 1]]
				}

		return Results
	}

// 把position赋值给PositionsRank然后对其进行排序
func (this *Mymethod)rankPositions(){
		// 首先判断是否需要排序
		if PositionToTargetSorted == 1{
			for position, _ := range PositionToTarget{
					PositionsRank = append(PositionsRank, position)
				}
			sort.Strings(PositionsRank)
			PositionToTargetSorted = 2
		}
	}

func main(){
		var m Mymethod
		// 测试加入节点
		m.AddTarget("192.168.1.114")
		m.AddTarget("192.168.1.35")
		var count int
		var count1 int
		fmt.Println("================================")
		// 测试查找节点
		Results := m.LookUp("202.168.1.172")
		// 做一个命中率的测试
		for i:= 0;i<10000;i++{
				source := rand.Intn(1000000)
				sourceStr := strconv.Itoa(source)
				Results = m.LookUp(sourceStr)
				if Results == "192.168.1.114"{
						count = count + 1
					}
				if Results == "192.168.1.35"{
						count1 = count1 + 1
					}
			}
		fmt.Printf("count is:%d", count)
		fmt.Println("================================")
		fmt.Printf("count1 is:%d", count1)
		fmt.Println("================================")
		fmt.Printf("Results1: %+v", Results)
		Results = m.LookUp("192.168.1.56")
		fmt.Printf("Results2: %+v", Results)
		fmt.Println("================================")
		// 测试删除节点
		m.RemoveTarget("192.168.1.114")
		Results = m.LookUp("192.168.1.5638")
		fmt.Printf("Results3: %+v", Results)
		// 测试多个节点中，需要排序后找到position
	}
