package main

import (
	"sync"

	"github.com/mikefaille/docker-group-daemonizer/dockerGrp"
	"github.com/mikefaille/docker-group-daemonizer/unixGrp"
	"github.com/mikefaille/tenus"
)

type error interface {
	Error() string
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	chanGroup := unixGrp.TakeAllGroups()

	for currentG := range chanGroup {

		gDocker, err := dockerGrp.CatchDockerEqGroup(currentG)

		if err == nil {
			wg.Add(1)
			go func() {
				gDocker.GenerateDockerDaemon()

				if ok, err := tenus.IsInterfaceExist(gDocker.Name); !ok || err != nil {
					check(err)
					gDocker.AddNewDockerBr()

				} else {
					err := tenus.DelBridgeWithName(gDocker.Name)

					check(err)
					gDocker.AddNewDockerBr()
				}
				wg.Done()
			}()
		}

	}
	wg.Done()
	wg.Wait()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
