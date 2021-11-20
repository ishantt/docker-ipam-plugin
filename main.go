package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	ipamApi "github.com/docker/go-plugins-helpers/ipam"
	"net"
)

const socketAddress = "/run/docker/plugins/sdip.sock"
const localAddressSpace = "LOCAL"
const globalAddressSpace = "GLOBAL"
const localAddressPool = "192.168.10.2/24"

var scs = spew.ConfigState{Indent: "  "}

type ipamDriver struct {
	allocatedIPAddresses map[string]struct{}
	networkAllocated     bool
}

func (i *ipamDriver) GetCapabilities() (*ipamApi.CapabilitiesResponse, error) {
	logrus.Infof("GetCapabilities called")
	return &ipamApi.CapabilitiesResponse{RequiresMACAddress: true}, nil
}

func (i *ipamDriver) GetDefaultAddressSpaces() (*ipamApi.AddressSpacesResponse, error) {
	logrus.Infof("GetDefaultAddressSpaces called")

	logrus.Infof("Returing response LocalDefaultAddressSpace: %s, GlobalDefaultAddressSpace: %s",
		localAddressSpace, globalAddressSpace)

	return &ipamApi.AddressSpacesResponse{LocalDefaultAddressSpace: localAddressSpace,
		GlobalDefaultAddressSpace: globalAddressSpace}, nil
}

func (i *ipamDriver) RequestPool(r *ipamApi.RequestPoolRequest) (*ipamApi.RequestPoolResponse, error) {
	if !i.networkAllocated {
		logrus.Infof("RequestPool called")

		rFormatted := scs.Sdump(r)
		logrus.Infof(rFormatted)

		logrus.Infof("Pool: %s", localAddressPool)
		i.networkAllocated = true
		return &ipamApi.RequestPoolResponse{PoolID: "1234", Pool: localAddressPool}, nil
	}
	return &ipamApi.RequestPoolResponse{}, errors.New("Pool Already Allocated")
}

func (i *ipamDriver) ReleasePool(r *ipamApi.ReleasePoolRequest) error {
	logrus.Infof("ReleasePool called")

	rFormatted := scs.Sdump(r)
	logrus.Infof(rFormatted)

	if r.PoolID == "1234" {
		logrus.Infof("Releasing Pool")
		i.networkAllocated = false
		i.allocatedIPAddresses = make(map[string]struct{})
	}
	return nil
}

func (i *ipamDriver) RequestAddress(r *ipamApi.RequestAddressRequest) (*ipamApi.RequestAddressResponse, error) {
	logrus.Infof("RequestAddress called")

	rFormatted := scs.Sdump(r)
	logrus.Infof(rFormatted)

	addr := i.getNextIP()
	addr = fmt.Sprintf("%s/%s", addr, "24")
	logrus.Infof("Allocated IP %s", addr)
	return &ipamApi.RequestAddressResponse{Address: addr}, nil
}

func (i *ipamDriver) ReleaseAddress(r *ipamApi.ReleaseAddressRequest) error {
	logrus.Infof("ReleaseAddress called")

	rFormatted := scs.Sdump(r)
	logrus.Infof(rFormatted)

	delete(i.allocatedIPAddresses, r.Address)
	if _, ok := i.allocatedIPAddresses[r.Address]; !ok {
		logrus.Infof("IP %s Released from the Pool", r.Address)
	}
	return nil
}

func (i *ipamDriver) getNextIP() string {
	ipAddr, ipNet, _ := net.ParseCIDR(localAddressPool)

	ret := ""
	for ip := ipAddr; ipNet.Contains(ip); incrementIP(ip) {
		if _, ok := i.allocatedIPAddresses[ip.String()]; !ok {
			ret = ip.String()
			i.allocatedIPAddresses[ret] = struct{}{}
			break
		}
	}
	return ret
}
func incrementIP(ip net.IP) {

	// length of IP is 16 bytes. This is because IPv6 address is 16 bytes.
	// For IPv4 , take the last octet and increment it by one.
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func main() {
	logrus.Infof("Starting Docker IPAM Plugin")
	i := &ipamDriver{allocatedIPAddresses: make(map[string]struct{})}
	h := ipamApi.NewHandler(i)
	logrus.Infof("Listening on socket %s", socketAddress)
	h.ServeUnix(socketAddress, 0)
}
