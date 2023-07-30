// Package deployer for handling deployments
package deployer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/codescalers/cloud4students/middlewares"
	"github.com/codescalers/cloud4students/models"
	"github.com/codescalers/cloud4students/streams"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
	"gorm.io/gorm"
)

func buildK8sCluster(node uint32, sshKey, network string, k models.K8sDeployInput) (workloads.K8sCluster, error) {
	master := workloads.K8sNode{
		Name:      k.MasterName,
		Flist:     k8sFlist,
		Planetary: true,
		Node:      node,
	}
	cru, mru, sru, ips, err := calcNodeResources(k.Resources, k.Public)
	if err != nil {
		return workloads.K8sCluster{}, err
	}
	master.CPU = int(cru)
	master.Memory = int(mru * 1024)
	master.DiskSize = int(sru)
	if ips == 1 {
		master.PublicIP = true
	}

	workers := []workloads.K8sNode{}
	for _, worker := range k.Workers {
		w := workloads.K8sNode{
			Name:  worker.Name,
			Flist: k8sFlist,
			Node:  node,
		}
		cru, mru, sru, _, err := calcNodeResources(k.Resources, false)
		if err != nil {
			return workloads.K8sCluster{}, err
		}
		w.CPU = int(cru)
		w.Memory = int(mru * 1024)
		w.DiskSize = int(sru)
		workers = append(workers, w)
	}
	k8sCluster := workloads.K8sCluster{
		Master:       &master,
		Workers:      workers,
		NetworkName:  network,
		Token:        token,
		SSHKey:       sshKey,
		SolutionType: k.MasterName,
	}

	return k8sCluster, nil
}

func (d *Deployer) deployK8sClusterWithNetwork(ctx context.Context, k8sDeployInput models.K8sDeployInput, sshKey string, adminSSHKey string) (uint32, uint64, uint64, error) {
	// get available nodes
	node, err := d.getK8sAvailableNode(ctx, k8sDeployInput)
	if err != nil {
		return 0, 0, 0, err
	}

	// build network
	network := buildNetwork(node, fmt.Sprintf("%sk8sNet", k8sDeployInput.MasterName))

	// build cluster
	cluster, err := buildK8sCluster(node,
		sshKey+"\n"+adminSSHKey,
		network.Name,
		k8sDeployInput,
	)
	if err != nil {
		return 0, 0, 0, err
	}

	// add network and cluster to be deployed
	err = d.Redis.PushK8s(streams.K8sDeployment{Net: &network, DL: &cluster})
	if err != nil {
		return 0, 0, 0, err
	}

	// wait for deployments
	for {
		if <-d.k8sDeployed {
			break
		}
	}

	// checks that network and k8s are deployed successfully
	loadedNet, err := d.tfPluginClient.State.LoadNetworkFromGrid(cluster.NetworkName)
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "failed to load network '%s' on nodes %v", cluster.NetworkName, network.Nodes)
	}

	loadedCluster, err := d.tfPluginClient.State.LoadK8sFromGrid([]uint32{node}, cluster.Master.Name)
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "failed to load kubernetes cluster '%s' on nodes %v", cluster.Master.Name, network.Nodes)
	}

	return node, loadedNet.NodeDeploymentID[node], loadedCluster.NodeDeploymentID[node], nil
}

func (d *Deployer) loadK8s(k8sDeployInput models.K8sDeployInput, userID string, node uint32, networkContractID uint64, k8sContractID uint64, expirationToleranceInDays int) (models.K8sCluster, error) {
	// load cluster
	resCluster, err := d.tfPluginClient.State.LoadK8sFromGrid([]uint32{node}, k8sDeployInput.MasterName)
	if err != nil {
		return models.K8sCluster{}, err
	}

	// save to db
	cru, mru, sru, _, err := calcNodeResources(k8sDeployInput.Resources, k8sDeployInput.Public)
	if err != nil {
		return models.K8sCluster{}, err
	}
	master := models.Master{
		CRU:       cru,
		MRU:       mru,
		SRU:       sru,
		Public:    k8sDeployInput.Public,
		PublicIP:  resCluster.Master.ComputedIP,
		Name:      k8sDeployInput.MasterName,
		YggIP:     resCluster.Master.YggIP,
		Resources: k8sDeployInput.Resources,
	}
	workers := []models.Worker{}
	for _, worker := range k8sDeployInput.Workers {

		cru, mru, sru, _, err := calcNodeResources(worker.Resources, false)
		if err != nil {
			return models.K8sCluster{}, err
		}
		workerModel := models.Worker{
			Name:      worker.Name,
			CRU:       cru,
			MRU:       mru,
			SRU:       sru,
			Resources: worker.Resources,
		}
		workers = append(workers, workerModel)
	}

	pkg, err := d.db.GetPkgByID(k8sDeployInput.PkgID)
	if err != nil {
		return models.K8sCluster{}, err
	}

	k8sCluster := models.K8sCluster{
		UserID:          userID,
		NetworkContract: int(networkContractID),
		ClusterContract: int(k8sContractID),
		Master:          master,
		Workers:         workers,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().AddDate(0, pkg.PeriodInMonth, expirationToleranceInDays),
	}

	return k8sCluster, nil
}

func (d *Deployer) getK8sAvailableNode(ctx context.Context, k models.K8sDeployInput) (uint32, error) {
	_, mru, sru, ips, err := calcNodeResources(k.Resources, k.Public)
	if err != nil {
		return 0, err
	}

	for _, worker := range k.Workers {
		_, m, s, _, err := calcNodeResources(worker.Resources, false)
		if err != nil {
			return 0, err
		}
		mru += m
		sru += s
	}

	freeMRU := uint64(mru)
	freeSRU := uint64(sru)
	filter := types.NodeFilter{
		Status:  &statusUp,
		FreeMRU: &freeMRU,
		FreeSRU: &freeSRU,
		FreeIPs: &ips,
		FarmIDs: []uint64{1},
		IPv6:    &trueVal,
	}

	nodes, err := deployer.FilterNodes(ctx, d.tfPluginClient, filter)
	if err != nil {
		return 0, err
	}

	return uint32(nodes[0].NodeID), nil
}

// ValidateK8sQuota validates the quota a k8s deployment need
func ValidateK8sQuota(k models.K8sDeployInput, availableResourcesQuota, availablePublicIPsQuota int) (int, error) {
	neededQuota, err := calcNeededQuota(k.Resources)
	if err != nil {
		return 0, err
	}

	for _, worker := range k.Workers {
		workerQuota, err := calcNeededQuota(worker.Resources)
		if err != nil {
			return 0, err
		}
		neededQuota += workerQuota
	}

	if availableResourcesQuota < neededQuota {
		return 0, fmt.Errorf("no available quota %d for kubernetes deployment, you can request a new voucher", availableResourcesQuota)
	}
	if k.Public && availablePublicIPsQuota < publicQuota {
		return 0, fmt.Errorf("no available quota %d for public ips", availablePublicIPsQuota)
	}

	return neededQuota, nil
}

func (d *Deployer) deployK8sRequest(ctx context.Context, user models.User, k8sDeployInput models.K8sDeployInput, adminSSHKey string, expirationToleranceInDays int) (int, error) {
	pkg, err := d.db.GetPkgByID(k8sDeployInput.PkgID)
	if err != nil {
		log.Error().Err(err).Send()
		return http.StatusInternalServerError, errors.New(internalServerErrorMsg)
	}

	neededQuota, err := ValidateK8sQuota(k8sDeployInput, pkg.Vms, pkg.PublicIPs)
	if err != nil {
		log.Error().Err(err).Send()
		return http.StatusBadRequest, err
	}

	// deploy network and cluster
	node, networkContractID, k8sContractID, err := d.deployK8sClusterWithNetwork(ctx, k8sDeployInput, user.SSHKey, adminSSHKey)
	if err != nil {
		log.Error().Err(err).Send()
		return http.StatusInternalServerError, errors.New(internalServerErrorMsg)
	}

	k8sCluster, err := d.loadK8s(k8sDeployInput, user.ID.String(), node, networkContractID, k8sContractID, expirationToleranceInDays)
	if err != nil {
		log.Error().Err(err).Send()
		return http.StatusInternalServerError, errors.New(internalServerErrorMsg)
	}
	publicIPsQuota := pkg.PublicIPs
	if k8sDeployInput.Public {
		publicIPsQuota -= publicQuota
	}

	// update package
	err = d.db.UpdateUserPackage(user.ID.String(), pkg.Vms-neededQuota, publicIPsQuota)
	if err == gorm.ErrRecordNotFound {
		return http.StatusNotFound, errors.New("user quota is not found")
	}
	if err != nil {
		log.Error().Err(err).Send()
		return http.StatusInternalServerError, errors.New(internalServerErrorMsg)
	}

	err = d.db.CreateK8s(&k8sCluster)
	if err != nil {
		log.Error().Err(err).Send()
		return http.StatusInternalServerError, errors.New(internalServerErrorMsg)
	}

	// metrics
	middlewares.Deployments.WithLabelValues(user.ID.String(), k8sDeployInput.Resources, "master").Inc()
	for _, worker := range k8sDeployInput.Workers {
		middlewares.Deployments.WithLabelValues(user.ID.String(), worker.Resources, "worker").Inc()
	}

	return 0, nil
}
