package nfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"time"
)

type Driver struct {
	name     string
	nodeID   string
	version  string
	endpoint string

	perm *uint32
	// CSI init parameter

	parameter1  string
	parameter2  int
	parameter3  time.Duration
	ns          *NodeServer
	cap         map[csi.VolumeCapability_AccessMode_Mode]bool
	cscap       []*csi.ControllerServiceCapability
	nscap       []*csi.NodeServiceCapability
	acp         []*csi.VolumeCapability_AccessMode
	volumeLocks *VolumeLocks
}

const (
	DriverName = "nfs.csi.k8s.io"
	// Address of the NFS server
	paramServer = "server"
	// Base directory of the NFS server to create volumes under.
	// The base directory must be a direct child of the root directory.
	// The root directory is omitted from the string, for example:
	//     "base" instead of "/base"
	paramShare = "share"
)

func NewCSIDriver(name, version, nodeID, endpoint, parameter1 string, parameter2 int, parameter3 time.Duration) *Driver {
	logrus.Info("Driver: %s version is %s", name, version)

	// parameter check  todo

	driver := &Driver{
		name:       name,
		nodeID:     nodeID,
		version:    version,
		endpoint:   endpoint,
		parameter1: parameter1,
		parameter2: parameter2,
		parameter3: parameter3,
	}
	// Specify how a volume can be accessed.
	driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		// writer
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		// reader_only
		csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
		csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
		csi.VolumeCapability_AccessMode_MULTI_NODE_SINGLE_WRITER,
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	})
	// you can  define the  capablities
	driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		// create volume
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		// delete snapshot
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
	})
	return driver
}

func (d *Driver) AddVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) {
	// accessMode
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		logrus.Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, &csi.VolumeCapability_AccessMode{Mode: c})
	}
	d.acp = vca
}

func (d *Driver) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability
	for _, c := range cl {
		logrus.Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}
	d.cscap = csc
}

// register  IdentityServer NodeServer ControllerServer  and Start the grpc server
func (d *Driver) Run() {
	server := NewNonBlockingGRPCServer()
	// start server
	server.Start(d.endpoint,
		NewIdentityServer(d),
		NewControllerServer(d),
		NewNodeServer(d))
	server.Wait()
}
