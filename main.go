package main

import (
	"bufio"
	"fmt"
	"os"
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

func FindShortestPath(antFarm *AntFarm) ([]string, error) {
	visited := make(map[string]bool)

	queue := [][]string{}

	queue = append(queue, []string{antFarm.Start})

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		fmt.Println("path", path)
		fmt.Println("queue", queue)

		currentRoom := path[len(path)-1]

		if currentRoom == antFarm.End {
			return path, nil
		}

		visited[currentRoom] = true

		fmt.Println("currentRoom", currentRoom)

		for _, tunnel := range antFarm.Tunnels {
			if tunnel.Room1 == currentRoom && !visited[tunnel.Room2] {
				fmt.Println("Room1", tunnel)
				newPath := append(path, tunnel.Room2)
				queue = append(queue, newPath)
				visited[tunnel.Room2] = true
			} else if tunnel.Room2 == currentRoom && !visited[tunnel.Room1] {
				fmt.Println("Room2", tunnel)
				newPath := append(path, tunnel.Room1)
				queue = append(queue, newPath)
				visited[tunnel.Room1] = true
			}
		}
	}

	return nil, fmt.Errorf("no path found")
}

func MoveAnts(antFarm *AntFarm, path []string) {
	antPos := make(map[int]int) // Карта для отслеживания позиций муравьев
	for i := 1; i <= antFarm.AntCount; i++ {
		antPos[i] = 0 // Изначально все муравьи находятся в начальной комнате
	}

	// Двигаем муравьев по пути
	for {
		moved := false
		for antID, pos := range antPos {
			if pos < len(path)-1 {
				// Если муравей еще не достиг конечной комнаты
				currentRoom := path[pos]
				nextRoom := path[pos+1]
				for _, tunnel := range antFarm.Tunnels {
					// Ищем туннел, который соединяет текущую комнату с следующей
					if (tunnel.Room1 == currentRoom && tunnel.Room2 == nextRoom) ||
						(tunnel.Room1 == nextRoom && tunnel.Room2 == currentRoom) {
						// Если туннел найден, перемещаем муравья
						fmt.Printf("L%d-%s ", antID, nextRoom)
						antPos[antID]++
						moved = true
						break
					}
				}
			}
		}
		fmt.Println() // Отделяем ходы муравьев друг от друга

		if !moved {
			break // Если ни один муравей не смог двинуться, заканчиваем цикл
		}
	}
}

func FindAndMoveAnts(antFarm *AntFarm) error {
	// Находим кратчайший путь
	shortestPath, err := FindShortestPath(antFarm)
	if err != nil {
		return err
	}

	// Двигаем муравьев по найденному пути
	MoveAnts(antFarm, shortestPath)

	return nil
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
	if antFarm != nil {
		err := FindAndMoveAnts(antFarm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
