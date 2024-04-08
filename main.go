package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Room struct {
	Name string
	X    int
	Y    int
}

type Path struct {
	Rooms  []string // список комнат, через которые проходит путь
	Length int      // общее количество комнат в пути
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
		// Добавляем найденный путь в список всех путей
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

	// После завершения поиска пути, снимаем метку посещения комнаты и убираем ее из пути
	visited[currentRoom] = false
	path = path[:len(path)-1]
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
	allPaths1 := make([][]string, 0)
	allPaths1 = append(allPaths1, allPaths[0])
	// Вывод всех найденных путей
	for _, path := range allPaths {
		co := 0
		for i := 1; i < len(path)-1; i++ {
			if !RemoveMassiv(path[i], allPaths1) {
				break
			} else {
				co++
			}
		}
		if co != 0 && co == len(path)-2 {
			allPaths1 = append(allPaths1, path)
		}
	}
	for _, path := range allPaths1 {
		fmt.Println(path)
	}
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
