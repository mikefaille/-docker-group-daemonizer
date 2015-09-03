package unixGrp

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Group struct {
	Name    string
	Members []string
	Guid    int64
}

func TakeAllGroups() (chanGroup chan Group) {

	currentChanGroup := make(chan Group, 2)
	file, err := os.Open("/etc/group")
	check(err)

	scanner := bufio.NewScanner(file)

	go func() {

		for scanner.Scan() {

			out := scanner.Bytes()
			lineB := make([]byte, len(out))
			copy(lineB, out)
			line := string(lineB[:])

			currentG := takeGroupArray(line)

			currentChanGroup <- currentG

		}

		close(currentChanGroup)
		file.Close()
	}()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return currentChanGroup

}

func takeGroupArray(groupLine string) Group {
	group := strings.Split(groupLine, ":")
	guid := group[2]
	intGuid, _ := strconv.ParseInt(guid, 10, 64)
	members := strings.Split(group[3], ",")
	g := Group{Name: group[0], Guid: intGuid, Members: members}
	return g
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
