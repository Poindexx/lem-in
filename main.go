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

type AntInfo struct {
	Id    int
	Room  []string
	Ves   []int
	Mesto int
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
	aaall = RemoveDuplicates(aaall)
	aaall = removeDuplicateArrays(aaall)

	countAnts := antFarm.AntCount
	antAllPaths := AntAllPaths(aaall, countAnts)

	var minAntAllPath *AntAllPath

	for _, antAllPath := range antAllPaths {
		if minAntAllPath == nil || antAllPath.FinnalTunnels < minAntAllPath.FinnalTunnels {
			minAntAllPath = antAllPath
		}
	}

	for i := 0; i < len(minAntAllPath.ArrayRooms); i++ {
		minAntAllPath.ArrayRooms[i] = minAntAllPath.ArrayRooms[i][1:]
	}
	ants := RaspredelenyeAnt(minAntAllPath, countAnts)
	printAntMovements(ants, minAntAllPath)

}

func RaspredelenyeAnt(minAntAllPath *AntAllPath, countAnts int) []AntInfo {
	ai := make([]AntInfo, countAnts)
	for i := 0; i < countAnts; i++ {
		ai[i].Ves = make([]int, len(minAntAllPath.ArrayRooms))
		for j := 0; j < len(minAntAllPath.ArrayRooms); j++ {
			ai[i].Ves[j] = len(minAntAllPath.ArrayRooms[j])
		}
	}

	for i := 0; i < countAnts; i++ {
		ai[i].Id = i
	}

	aa := ai[0].Ves
	for i := 0; i < countAnts; i++ {
		l := aa[0]
		id := 0
		for i1 := 0; i1 < len(aa); i1++ {
			if !(aa[i1] > l) {
				l = aa[i1]
				id = i1
			}
		}
		aa[id] = aa[id] + 1
		ai[i].Room = minAntAllPath.ArrayRooms[id]
	}
	return ai
}

func RemoveDuplicates(arr [][][]string) [][][]string {
	result := make([][][]string, len(arr))
	for i := 0; i < len(arr); i++ {
		uniqueMap := make(map[string]bool)
		var uniqueArr [][]string
		for j := 0; j < len(arr[i]); j++ {
			str := strings.Join(arr[i][j], "")
			if !uniqueMap[str] {
				uniqueMap[str] = true
				uniqueArr = append(uniqueArr, arr[i][j])
			}
		}
		result[i] = uniqueArr
	}

	return result
}

func printAntMovements(ants []AntInfo, minAntAllPath *AntAllPath) {

	for q := 0; q < minAntAllPath.FinnalTunnels-1; q++ {
		tekser := make(map[int]string)
		for i := 0; i < len(ants); i++ {
			if len(ants[i].Room) == 1 {
				if ants[i].Mesto < len(ants[i].Room) && ProverkaMassiva(tekser, ants[i].Room[ants[i].Mesto]) {
					tekser[ants[i].Id] = ants[i].Room[ants[i].Mesto]
					ants[i].Mesto++
				}
			} else {
				if ants[i].Mesto < len(ants[i].Room) && (ants[i].Mesto == len(ants[i].Room)-1 || ProverkaMassiva(tekser, ants[i].Room[ants[i].Mesto])) {
					tekser[ants[i].Id] = ants[i].Room[ants[i].Mesto]
					ants[i].Mesto++
				}
			}
		}
		for id, value := range tekser {
			fmt.Printf("l%d-%s ", id+1, value)
		}
		fmt.Println()
	}
}

func ProverkaMassiva(tekser map[int]string, soz string) bool {
	for _, value := range tekser {
		if value == soz {
			return false
		}
	}
	return true
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

// ReadAntFarm считывает данные муравейника из файла и возвращает объект AntFarm.
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
	var startSet, endSet bool // флаги для установки начальной и конечной комнат
	for scanner.Scan() {
		line := scanner.Text()
		parts_t := strings.Fields(line)
		if strings.HasPrefix(line, "##start") || strings.HasPrefix(line, "##end") {
			if !scanner.Scan() {
				return nil, fmt.Errorf("отсутствует следующая строка после %s", line)
			}
			nextLine := scanner.Text()
			parts := strings.Fields(nextLine)
			if len(parts) < 3 {
				return nil, fmt.Errorf("неверный формат комнаты: %s", nextLine)
			}
			name := parts[0]
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("неверная координата X: %s", parts[1])
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("неверная координата Y: %s", parts[2])
			}

			if _, ok := antFarm.Rooms[name]; ok {
				return nil, fmt.Errorf("комната с именем %s уже существует", name)
			}

			if strings.HasPrefix(line, "##start") {
				if startSet {
					return nil, fmt.Errorf("дублирующая строка ##start")
				}
				startSet = true
				antFarm.Start = name
			} else {
				if endSet {
					return nil, fmt.Errorf("дублирующая строка ##end")
				}
				endSet = true
				antFarm.End = name
			}
			antFarm.Rooms[name] = Room{Name: name, X: x, Y: y}
		} else if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "L") && len(parts_t) == 3 {
			name := parts_t[0]
			x, _ := strconv.Atoi(parts_t[1])
			y, _ := strconv.Atoi(parts_t[2])

			if _, ok := antFarm.Rooms[name]; ok {
				return nil, fmt.Errorf("комната с именем %s уже существует", name)
			}

			antFarm.Rooms[name] = Room{Name: name, X: x, Y: y}
		} else if strings.Count(line, "-") == 1 {
			parts := strings.Split(line, "-")
			room1 := parts[0]
			room2 := parts[1]

			// Проверка существования комнаты room1 и room2
			if _, ok := antFarm.Rooms[room1]; !ok {
				return nil, fmt.Errorf("комната %s не существует", room1)
			}
			if _, ok := antFarm.Rooms[room2]; !ok {
				return nil, fmt.Errorf("комната %s не существует", room2)
			}

			tunnel := Tunnel{Room1: room1, Room2: room2}
			antFarm.Tunnels = append(antFarm.Tunnels, tunnel)
		} else if len(parts_t) == 1 && !strings.HasPrefix(line, "#") {
			antCount, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("неверный формат количества муравьев: %s", line)
			}
			if antCount <= 0 {
				return nil, fmt.Errorf("неверный формат количества муравьев: %s", line)
			}
			antFarm.AntCount = antCount
		}
	}
	if !startSet {
		return nil, fmt.Errorf("не указана начальная комната ##start")
	}
	if !endSet {
		return nil, fmt.Errorf("не указана конечная комната ##end")
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
		return nil, fmt.Errorf("не правильный путь")
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
