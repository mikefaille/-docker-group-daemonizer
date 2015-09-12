package dockerGrp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"
	"text/template"

	c "github.com/alyu/configparser"

	"github.com/mikefaille/docker-group-daemonizer/unixGrp"
	"github.com/mikefaille/tenus"
)

type DockerGroup struct {
	unixGrp.Group
	Number  int64
	Options string
}

func (gDocker DockerGroup) AddNewDockerBr() {

	br, err := tenus.NewBridgeWithName(gDocker.Name)
	if err != nil {
		log.Fatal(err)
	}
	teamSegmentIP := gDocker.Number + 100
	brNet := fmt.Sprint("192.168.", teamSegmentIP, ".1/24")
	brIp, brIpNet, err := net.ParseCIDR(brNet)
	if err != nil {
		log.Fatal(err)
	}

	if err := br.SetLinkIp(brIp, brIpNet); err != nil {
		fmt.Println(err)
	}

	if err = br.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

}

func (dGroup DockerGroup) GenerateDockerDaemon() {

	tmpl, err := template.New("dockerOpts").Parse("-b  {{.Name}}  -g /var/lib/{{.Name}}  -G {{.Name}}  --exec-root=/var/run/docker{{.Number}} --pidfile=\"/var/run/{{.Name}}.pid\" --label=[\"equipe={{.Number}}\"] -H unix:///var/run/{{.Name}}.sock")
	if err != nil {
		panic(err)
	}

	optsBuf := new(bytes.Buffer)

	err = tmpl.Execute(optsBuf, dGroup)
	dGroup.Options = optsBuf.String()

	dGroup.generateUpstartDaemon()

	// TODO : Detect Daemon Manager

	// lsb := exec.Command("cat", "/etc/lsb-release", "|", "grep", "DISTRIB_RELEASE")
	// lsbOutput, err := lsb.CombinedOutput()
	// check(err)

	// switch string(lsbOutput) {
	// case "DISTRIB_RELEASE=14.04":
	// 	dGroup.generateUpstartDaemon()

	// default:
	// 	dGroup.GenerateSystemdService()
	// 	dGroup.GenerateSystemdSocket()
	// }

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

func (dGroup DockerGroup) GenerateSystemdService() {
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
	execStart := fmt.Sprint("/usr/bin/docker daemon ", dGroup.Options)
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

func (dGroup DockerGroup) GenerateSystemdSocket() {
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

func (dGroup DockerGroup) generateUpstartDaemon() {

	// TODO remove old template as file ?
	// dat, err := ioutil.ReadFile("docker-upstart.tmpl")
	// if err != nil {
	// 	fmt.Println("File named \"docker-upstart.tmpl\" must be loaded")
	// }

	tmpl, err := template.New("upstartDaemon").Parse(string(GetUpstartTemplate()))
	if err != nil {
		panic(err)
	}

	check(err)
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, dGroup)

	upstartFilePath := fmt.Sprint("/etc/init/", dGroup.Name, ".conf")
	err = ioutil.WriteFile(upstartFilePath, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	dockerOPTS := fmt.Sprint("DOCKER_OPTS=\"", dGroup.Options, "\n")
	defaultFilePath := fmt.Sprint("/etc/default/", dGroup.Name)
	err = ioutil.WriteFile(defaultFilePath, []byte(dockerOPTS), 0644)
	if err != nil {
		panic(err)
	}

}
