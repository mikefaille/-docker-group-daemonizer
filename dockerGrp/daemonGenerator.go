package dockerGrp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"

	c "github.com/alyu/configparser"

	"github.com/mikefaille/docker-group-daemonizer/unixGrp"
	"github.com/mikefaille/tenus"
)

type DockerGroup struct {
	unixGrp.Group
	Number int64
}

func (dGroup DockerGroup) GenerateDockerSocket() {
	// [Unit]
	// Description=Docker Socket for the API
	// PartOf=docker.service

	// [Socket]
	// ListenStream=/var/run/docker.sock
	// SocketMode=0660
	// SocketUser=root
	// SocketGroup=docker

	// [Install]
	// WantedBy=sockets.target

	conf := c.NewConfiguration()
	sectionUnit := conf.NewSection("Unit")
	descr := fmt.Sprint("Docker Socket for the API for Team no", dGroup.Number)
	sectionUnit.Add("Description", descr)
	servicePath := fmt.Sprint(dGroup.Name, ".service")
	sectionUnit.Add("PartOf", servicePath)

	sectionSocket := conf.NewSection("Socket")
	socketPath := fmt.Sprint("/var/run/", dGroup.Name, ".socket")
	sectionSocket.Add("ListenStream", socketPath)
	sectionSocket.Add("SocketMode", "0660")
	sectionSocket.Add("SocketUser", "root")
	sectionSocket.Add("SocketGroup", dGroup.Name)

	sectionInstall := conf.NewSection("Install")
	sectionInstall.Add("WantedBy", "socket.target")

	fmt.Println(conf)
	socketConfigPath := fmt.Sprint("/etc/systemd/system/", dGroup.Name, ".socket")
	err := c.Save(conf, socketConfigPath)
	if err != nil {
		log.Println(err)
	}
}

func (gDocker DockerGroup) AddNewDockerBr() {

	br, err := tenus.NewBridgeWithName(gDocker.Name)
	if err != nil {
		log.Fatal(err)
	}

	brNet := fmt.Sprint("192.168.", gDocker.Number, ".1/24")
	brIp, brIpNet, err := net.ParseCIDR(brNet)
	if err != nil {
		log.Fatal(err)
	}

	if err := br.SetLinkIp(brIp, brIpNet); err != nil {
		fmt.Println(err)
	}

	// Nécéssaire car il semble y avoir un «Race condition»
	time.Sleep(time.Millisecond * 100)
	if err = br.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

}

func (dGroup DockerGroup) GenerateDockerService() {
	// [Unit]
	// Description=Docker Application Container Engine
	// Documentation=https://docs.docker.com
	// After=network.target docker.socket
	// Requires=docker.socket

	// [Service]
	// Type=notify
	// ExecStart=/usr/bin/docker daemon -H fd://
	// MountFlags=slave
	// LimitNOFILE=1048576
	// LimitNPROC=1048576
	// LimitCORE=infinity

	// [Install]
	// WantedBy=multi-user.target

	conf := c.NewConfiguration()
	sectionUnit := conf.NewSection("Unit")
	descr := fmt.Sprint("Docker Application Container Engine for Team no", dGroup.Number)
	sectionUnit.Add("Description", descr)
	after := fmt.Sprint("network.target ", dGroup.Name, ".socket")
	sectionUnit.Add("After", after)

	sectionService := conf.NewSection("Service")
	sectionService.Add("Type", "notify")
	execStart := fmt.Sprint("/usr/bin/docker daemon -b ", dGroup.Name, " -g /var/lib/docker", strconv.FormatInt(dGroup.Number, 10), " -G ", dGroup.Name, " --exec-root=/var/run/docker", strconv.FormatInt(dGroup.Number, 10), " --pidfile=\"/var/run/", dGroup.Name, ".pid\" -H fd://")
	// execStart := fmt.Sprint("/usr/bin/docker daemon -b docker", strconv.FormatInt(dGroup.Number, 10), " -g /var/lib/docker", " -G ", dGroup.Name, " --exec-root=/var/run/docker --pidfile=\"/var/run/", dGroup.Name, ".pid\" --bip 192.168.", dGroup.Number, ".1/24", " -H fd://")
	sectionService.Add("ExecStart", execStart)
	sectionService.Add("MountFlags", "slave")
	sectionService.Add("LimitNOFILE", "1048576")
	sectionService.Add("LimitNPROC", "1048576")
	sectionService.Add("LimitCORE", "infinity")

	sectionInstall := conf.NewSection("Install")
	sectionInstall.Add("WantedBy", "multi-user.target")

	fmt.Println(conf)
	confPath := fmt.Sprint("/etc/systemd/system/", dGroup.Name, ".service")
	err := c.Save(conf, confPath)
	if err != nil {
		log.Println(err)
	}
}

func CatchDockerEqGroup(group unixGrp.Group) (DockerGroup, error) {
	var validDockerGroup = regexp.MustCompile(`^docker-eq([0-9]+)$`)
	if validDockerGroup.MatchString(group.Name) {
		num, err := strconv.ParseInt(validDockerGroup.FindStringSubmatch(group.Name)[1], 10, 64)
		check(err)
		return DockerGroup{Group: group, Number: int64(num)}, nil

	}
	return DockerGroup{Group: unixGrp.Group{Name: "", Members: nil}, Number: 0}, errors.New("This is not a docker group")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
