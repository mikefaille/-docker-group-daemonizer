package unixGrp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Group struct {
	Name    string
	Members []string
}

func TakeAllGroups() (chanGroup chan Group) {

	currentChanGroup := make(chan Group, 2)
	file, err := os.Open("/etc/group")
	check(err)

	scanner := bufio.NewScanner(file)

	go func() {

		for scanner.Scan() {
			fmt.Println(1)
			out := scanner.Bytes()
			lineB := make([]byte, len(out))
			copy(lineB, out)
			line := string(lineB[:])

			currentG := takeGroup(line)

			members, err := takeGroupMember(line)
			if err == nil {
				currentG.Members = members
			}

			currentChanGroup <- currentG

		}
		fmt.Println(2)
		close(currentChanGroup)
		file.Close()
	}()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return currentChanGroup

}

func takeGroup(s string) Group {
	groupEnd := strings.Index(s, ":")
	g := Group{Name: s[:groupEnd], Members: []string{}}
	return g
}

func takeGroupMember(s string) (members []string, err error) {

	membersString := strings.LastIndex(s, ":")
	members = strings.Split(s[membersString+1:], ",")
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
