package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	Name string
	X    int
	Y    int
}

type Path struct {
	Rooms  []string
	Length int
}

type Tunnel struct {
	Room1 string
	Room2 string
}

type AntFarm struct {
	AntCount int
	Rooms    map[string]Room
	Tunnels  []Tunnel
	Start    string
	End      string
}

type AntAllPath struct {
	ArrayRooms     [][]string
	MaxL           int
	MaxAnts        int
	ExcessAnts     int
	AdditionTunels int
	FinnalTunnels  int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}
	filename := os.Args[1]
	antFarm, err := ReadAntFarm(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	allPaths, err := FindAllPathsWrapper(antFarm)
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	allPaths3 := make([][][]string, 0)
	arip := make([]string, 0)
	for _, path := range allPaths {
		if !AripT(path[1], arip) {
			arip = append(arip, path[1])
			allPaths3 = append(allPaths3, RazdeitMassiv(path[1], allPaths))
		}
	}

	allOnly := make([][]string, 0)
	allOnly2 := make([][][]string, 0)
	co := 0
	p := 0
	aaall := SerchAll2(co, allPaths3, allOnly, allOnly2, p)
	sort3DArray(aaall)
	aaall = removeDuplicateArrays(aaall)
	countAnts := antFarm.AntCount
	allAnts := AntAllPaths(aaall, countAnts)
	for _, e := range allAnts {
		fmt.Println(e)
	}

}

func AntAllPaths(aaall [][][]string, n int) []*AntAllPath {
	antAllPaths := make([]*AntAllPath, len(aaall))

	for i, paths := range aaall {
		antAllPath := &AntAllPath{}
		for _, path := range paths {
			l := len(path)
			if l > antAllPath.MaxL {
				antAllPath.MaxL = l
			}
		}
		for _, path := range paths {
			antAllPath.MaxAnts += antAllPath.MaxL - len(path) + 1
		}

		antAllPath.ExcessAnts = n - antAllPath.MaxAnts
		antAllPath.AdditionTunels = (antAllPath.ExcessAnts + len(paths) - 1) / len(paths)
		if n <= antAllPath.MaxAnts {
			antAllPath.FinnalTunnels = antAllPath.MaxL
		} else {
			antAllPath.FinnalTunnels = antAllPath.AdditionTunels + antAllPath.MaxL
		}

		antAllPath.ArrayRooms = paths
		antAllPaths[i] = antAllPath
	}

	return antAllPaths
}

func ReadAntFarm(filename string) (*AntFarm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	antFarm := &AntFarm{
		Rooms: make(map[string]Room),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts_t := strings.Fields(line)
		if strings.HasPrefix(line, "##start") || strings.HasPrefix(line, "##end") {
			if !scanner.Scan() {
				return nil, fmt.Errorf("no next line after %s", line)
			}
			nextLine := scanner.Text()
			parts := strings.Fields(nextLine)
			if len(parts) < 3 {
				return nil, fmt.Errorf("invalid room format: %s", nextLine)
			}
			name := parts[0]
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid X coordinate: %s", parts[1])
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("invalid Y coordinate: %s", parts[2])
			}
			if strings.HasPrefix(line, "##start") {
				antFarm.Start = name
			} else {
				antFarm.End = name
			}
			antFarm.Rooms[name] = Room{Name: name, X: x, Y: y}
		} else if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "L") && len(parts_t) == 3 {
			name := parts_t[0]
			x, _ := strconv.Atoi(parts_t[1])
			y, _ := strconv.Atoi(parts_t[2])
			antFarm.Rooms[name] = Room{Name: name, X: x, Y: y}
		} else if strings.Count(line, "-") == 1 {
			parts := strings.Split(line, "-")
			tunnel := Tunnel{Room1: parts[0], Room2: parts[1]}
			antFarm.Tunnels = append(antFarm.Tunnels, tunnel)
		} else if len(parts_t) == 1 && !strings.HasPrefix(line, "#") {
			antCount, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("invalid ant count format: %s", line)
			}
			if antCount <= 0 {
				return nil, fmt.Errorf("invalid ant count format: %s", line)
			}
			antFarm.AntCount = antCount
		}
	}
	return antFarm, nil
}

func FindAllPaths(antFarm *AntFarm, currentRoom string, visited map[string]bool, path []string, allPaths *[][]string) {
	visited[currentRoom] = true
	path = append(path, currentRoom)

	if currentRoom == antFarm.End {
		*allPaths = append(*allPaths, append([]string{}, path...))
	} else {
		for _, tunnel := range antFarm.Tunnels {
			nextRoom := ""
			if tunnel.Room1 == currentRoom && !visited[tunnel.Room2] {
				nextRoom = tunnel.Room2
			} else if tunnel.Room2 == currentRoom && !visited[tunnel.Room1] {
				nextRoom = tunnel.Room1
			}

			if nextRoom != "" {
				FindAllPaths(antFarm, nextRoom, visited, path, allPaths)
			}
		}
	}

	visited[currentRoom] = false
}

func FindAllPathsWrapper(antFarm *AntFarm) ([][]string, error) {
	visited := make(map[string]bool)
	allPaths := make([][]string, 0)

	FindAllPaths(antFarm, antFarm.Start, visited, []string{}, &allPaths)

	if len(allPaths) == 0 {
		return nil, fmt.Errorf("no paths found")
	}

	return allPaths, nil
}

func sort3DArray(arr [][][]string) {
	for _, subArr := range arr {
		sort2DArray(subArr)
	}

	sort.Slice(arr, func(i, j int) bool {
		return countElements(arr[i]) < countElements(arr[j])
	})
}

func sort2DArray(arr [][]string) {
	sort.Slice(arr, func(i, j int) bool {
		return len(arr[i]) < len(arr[j])
	})
}

func countElements(arr [][]string) int {
	count := 0
	for _, subArr := range arr {
		for _, element := range subArr {
			if element != "" {
				count++
			}
		}
	}
	return count
}

func removeDuplicateArrays(allOnly2 [][][]string) [][][]string {
	uniqueArrays := make([][][]string, 0)
	encountered := make(map[string]bool)
	elementsCount := make(map[string]bool)

	for _, arr := range allOnly2 {
		key := arrayKey(arr)

		if elementsCount[key] {
			continue
		}

		if !encountered[key] {
			encountered[key] = true
			uniqueArrays = append(uniqueArrays, arr)
		}

		elementsCount[key] = true
	}

	return uniqueArrays
}

func arrayKey(arr [][]string) string {
	var key strings.Builder
	key.WriteString(fmt.Sprintf("%d:", len(arr)))

	for _, subArr := range arr {
		for _, element := range subArr {
			key.WriteString(fmt.Sprintf("%s,", element))
		}
	}

	return key.String()
}

func SerchAll2(co int, allPaths3 [][][]string, allOnly [][]string, allOnly2 [][][]string, p1 int) [][][]string {
	q := 0
	for i := 0; i < len(allPaths3); i++ {
		p := 0
		for k := 0; k < len(allPaths3[i]); k++ {
			if !proverkaAllOnly(allPaths3[i][k], allOnly) {
				allOnly = append(allOnly, allPaths3[i][k])
				p++
			}
		}
		if p == 0 {
			allOnly2 = append(allOnly2, append([][]string{}, allOnly...))
			allOnly = make([][]string, 0)
			q++
		}
	}
	if p1 != 1000 {
		allPaths3 = moveFirstToEnd(allPaths3)
		return SerchAll2(co, allPaths3, allOnly, allOnly2, p1+1)
	}
	return allOnly2
}

func moveFirstToEnd(arr [][][]string) [][][]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}

func proverkaAllOnly(a []string, b [][]string) bool {
	if len(b) == 0 {
		return false
	}
	for i := 0; i < len(b); i++ {
		for k := 1; k < len(b[i])-1; k++ {
			if !proverkaAllOnly1(b[i][k], a) {
				return true
			}
		}
	}
	return false
}

func proverkaAllOnly1(a string, b []string) bool {
	for i := 1; i < len(b)-1; i++ {
		if b[i] == a {
			return false
		}
	}
	return true
}

func RemoveMassiv(a string, b [][]string) bool {
	for _, path := range b {
		for i := 0; i < len(path); i++ {
			if path[i] == a {
				return false
			}
		}
	}
	return true
}

func AripT(a string, b []string) bool {
	for i := 0; i < len(b); i++ {
		if b[i] == a {
			return true
		}
	}
	return false
}

func RazdeitMassiv(a string, b [][]string) [][]string {
	allPaths := make([][]string, 0)
	for _, path := range b {
		if path[1] == a {
			allPaths = append(allPaths, path)
		}
	}
	return allPaths
}
