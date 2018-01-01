package main

import (
	"fmt"
	"time"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	var vmProps vmCloudProps
	cloudProps.As(&vmProps)

	pi, err := c.bakeryClient.BakePi(stemcellCID.AsString())
	if err != nil {
		return apiv1.VMCID{}, err
	}

	cid := apiv1.NewVMCID(pi.Id)

	//wait for the provisioning to finish
	piReadyChannel := make(chan bool)
	quitChannel := make(chan bool)

	go func() {
		piready := false
		for {
			select {
			case <-quitChannel:
				return
			default:
				piready, err = c.bakeryClient.IsPiBaked(cid.AsString())
				if err != nil {
					piReadyChannel <- false
					return
				}

				if piready {
					piReadyChannel <- true
					return
				}
			}
			time.Sleep(4 * time.Second) // poll once every 4 seconds
		}
	}()

	select {
	case success := <-piReadyChannel:
		if !success {
			c.bakeryClient.UnbakePi(cid.AsString())
			return apiv1.VMCID{}, fmt.Errorf("Error occured while waiting for Pi to be read. Rolled back deployment.")
		}
	case <-time.After(5 * time.Minute):
		c.bakeryClient.UnbakePi(cid.AsString())
		quitChannel <- true
		return apiv1.VMCID{}, fmt.Errorf("Waiting for Pi to be ready timed out. Rolled back deployment.")
	}
	//////////////////

	diskCID, err := c.CreateDisk(vmProps.EphemeralDisk, nil, &cid)
	if err != nil {
		c.bakeryClient.UnbakePi(cid.AsString())
		return apiv1.VMCID{}, err
	}

	err = c.AttachDisk(cid, diskCID)
	if err != nil {
		c.bakeryClient.UnbakePi(cid.AsString())
		return apiv1.VMCID{}, err
	}

	mac := fmt.Sprintf("b8:27:eb:%v:%v:%v", pi.Id[2:3], pi.Id[4:5], pi.Id[6:7]) //piId == serial number. last 6 digits of serial number==last 6 digits of mac. first 6 are rPi foundation mac range
	for _, network := range networks {
		if !network.IsDynamic() {
			network.SetMAC(mac)
		}
		break //only 1 network supported
	}

	err = c.UploadSettings(agentID, cid, networks, env, nil)
	if err != nil {
		//roll back the deploy
		c.bakeryClient.UnbakePi(cid.AsString())
		return apiv1.VMCID{}, err
	}

	err = c.bakeryClient.PowerCyclePi(cid.AsString())
	if err != nil {
		err = c.bakeryClient.UnbakePi(cid.AsString())
		return apiv1.VMCID{}, fmt.Errorf("Powering on failed. Rolled back deployment. Rollback result: %v", err.Error())
	}
	return cid, nil
}
