// Package deployer for handling deployments
package deployer

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/codescalers/cloud4students/models"
	"github.com/codescalers/cloud4students/streams"
	"github.com/codescalers/cloud4students/validators"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"gopkg.in/validator.v2"
)

const internalServerErrorMsg = "Something went wrong"

var (
	vmEntryPoint = "/init.sh"

	k8sFlist = "https://hub.grid.tf/tf-official-apps/threefoldtech-k3s-latest.flist"
	vmFlist  = "https://hub.grid.tf/tf-official-vms/ubuntu-22.04.flist"

	smallCPU     = uint64(1)
	smallMemory  = uint64(2)
	smallDisk    = uint64(25)
	mediumCPU    = uint64(2)
	mediumMemory = uint64(4)
	mediumDisk   = uint64(50)
	largeCPU     = uint64(4)
	largeMemory  = uint64(8)
	largeDisk    = uint64(100)

	smallQuota  = 1
	mediumQuota = 2
	largeQuota  = 3
	publicQuota = 1

	trueVal  = true
	statusUp = "up"

	token = "random"
)

// Deployer struct holds deployments configuration
type Deployer struct {
	db             models.DB
	Redis          streams.RedisClient
	tfPluginClient deployer.TFPluginClient

	vmDeployed  chan bool
	k8sDeployed chan bool
}

// NewDeployer create new deployer
func NewDeployer(db models.DB, redis streams.RedisClient, tfPluginClient deployer.TFPluginClient) (Deployer, error) {
	// validations
	err := validator.SetValidationFunc("ssh", validators.ValidateSSHKey)
	if err != nil {
		return Deployer{}, err
	}
	err = validator.SetValidationFunc("password", validators.ValidatePassword)
	if err != nil {
		return Deployer{}, err
	}
	err = validator.SetValidationFunc("mail", validators.ValidateMail)
	if err != nil {
		return Deployer{}, err
	}

	return Deployer{
		db,
		redis,
		tfPluginClient,
		make(chan bool),
		make(chan bool),
	}, nil
}

// PeriodicRequests for executing deployment api requests
func (d *Deployer) PeriodicRequests(ctx context.Context, sec int) {
	ticker := time.NewTicker(time.Second * time.Duration(sec))
	for range ticker.C {
		d.ConsumeVMRequest(ctx, false)
		d.ConsumeK8sRequest(ctx, false)
	}
}

// PeriodicDeploy for executing deployments
func (d *Deployer) PeriodicDeploy(ctx context.Context, sec int) {
	ticker := time.NewTicker(time.Second * time.Duration(sec))

	for range ticker.C {
		vmNets, vms, err := d.consumeVMs(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to consume vms")
		}

		k8sNets, clusters, err := d.consumeK8s(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to consume clusters")
		}

		if len(vms) > 0 {
			err := d.tfPluginClient.NetworkDeployer.BatchDeploy(ctx, vmNets)
			if err != nil {
				log.Error().Err(err).Msg("failed to batch deploy network")
			}

			err = d.tfPluginClient.DeploymentDeployer.BatchDeploy(ctx, vms)
			if err != nil {
				log.Error().Err(err).Msg("failed to batch deploy vm")
			}

			for i := 0; i < len(vms); i++ {
				d.vmDeployed <- true
			}
		}

		if len(clusters) > 0 {
			err := d.tfPluginClient.NetworkDeployer.BatchDeploy(ctx, k8sNets)
			if err != nil {
				log.Error().Err(err).Msg("failed to batch deploy network")
			}

			err = d.tfPluginClient.K8sDeployer.BatchDeploy(ctx, clusters)
			if err != nil {
				log.Error().Err(err).Msg("failed to batch deploy clusters")
			}

			for i := 0; i < len(clusters); i++ {
				d.k8sDeployed <- true
			}
		}
	}
}

// CancelDeployment cancel deployments from grid
func (d *Deployer) CancelDeployment(contractID uint64, netContractID uint64, dlType string, dlName string) error {
	// cancel deployment
	err := d.tfPluginClient.SubstrateConn.CancelContract(d.tfPluginClient.Identity, contractID)
	if err != nil {
		return err
	}

	// cancel network
	err = d.tfPluginClient.SubstrateConn.CancelContract(d.tfPluginClient.Identity, netContractID)
	if err != nil {
		return err
	}

	// update state
	for node, contracts := range d.tfPluginClient.State.CurrentNodeDeployments {
		contracts = workloads.Delete(contracts, contractID)
		contracts = workloads.Delete(contracts, netContractID)
		d.tfPluginClient.State.CurrentNodeDeployments[node] = contracts

		d.tfPluginClient.State.Networks.DeleteNetwork(fmt.Sprintf("%s%sNet", dlType, dlName))
	}

	return nil
}

func (d *Deployer) CleanExpiredVMs(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		users, err := d.db.ListAllUsers()
		if err != nil {
			log.Error().Err(err).Msg("failed to get all users")
			return
		}

		for _, user := range users {
			vms, err := d.db.GetAllVms(user.UserID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get all user vms")
				continue
			}

			for _, vm := range vms {
				if vm.ExpirationDate.Before(time.Now()) {
					err = d.CancelDeployment(vm.ContractID, vm.NetworkContractID, "vm", vm.Name)
					if err != nil {
						log.Error().Err(err).Msg("failed to cancel contract of expired vm")
					}
					err := d.db.DeleteVMByID(vm.ID)
					if err != nil {
						log.Error().Err(err).Msg("failed to delete expired vm")
					}
				}
			}
		}
	}
}

func (d *Deployer) CleanExpiredK8S(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		users, err := d.db.ListAllUsers()
		if err != nil {
			log.Error().Err(err).Msg("failed to get all users")
			return
		}

		for _, user := range users {
			k8s, err := d.db.GetAllK8s(user.UserID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get all user k8s clusters")
				continue
			}

			for _, k := range k8s {
				if k.ExpirationDate.Before(time.Now()) {
					err = d.CancelDeployment(uint64(k.ClusterContract), uint64(k.NetworkContract), "k8s", k.Master.Name)
					if err != nil {
						log.Error().Err(err).Msg("failed to cancel contract of expired k8s cluster")
					}
					err := d.db.DeleteVMByID(k.ID)
					if err != nil {
						log.Error().Err(err).Msg("failed to delete expired k8s cluster")
					}
				}
			}
		}
	}
}

func buildNetwork(node uint32, name string) workloads.ZNet {
	return workloads.ZNet{
		Name:  name,
		Nodes: []uint32{node},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		AddWGAccess: false,
	}
}

func calcNodeResources(resources string, public bool) (uint64, uint64, uint64, uint64, error) {
	var cru uint64
	var mru uint64
	var sru uint64
	var ips uint64
	switch resources {
	case "small":
		cru += smallCPU
		mru += smallMemory
		sru += smallDisk
	case "medium":
		cru += mediumCPU
		mru += mediumMemory
		sru += mediumDisk
	case "large":
		cru += largeCPU
		mru += largeMemory
		sru += largeDisk
	default:
		return 0, 0, 0, 0, fmt.Errorf("unknown resource type %s", resources)
	}
	if public {
		ips = 1
	}
	return cru, mru, sru, ips, nil
}

// choose suitable nodes based on needed resources
func filterNode(resource string, public bool) (types.NodeFilter, error) {
	cru, mru, sru, ips, err := calcNodeResources(resource, public)
	if err != nil {
		return types.NodeFilter{}, err
	}

	return types.NodeFilter{
		FarmIDs:  []uint64{1},
		TotalCRU: &cru,
		FreeSRU:  &sru,
		FreeMRU:  &mru,
		FreeIPs:  &ips,
		IPv4:     &trueVal,
		Status:   &statusUp,
		IPv6:     &trueVal,
	}, nil
}

func calcNeededQuota(resources string) (int, error) {
	var neededQuota int
	switch resources {
	case "small":
		neededQuota += smallQuota
	case "medium":
		neededQuota += mediumQuota
	case "large":
		neededQuota += largeQuota
	default:
		return 0, fmt.Errorf("unknown resource type %s", resources)
	}

	return neededQuota, nil
}
