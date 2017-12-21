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

	err = c.UploadEnvJson(agentID, cid, networks, env)
	if err != nil {
		//roll back the deploy
		c.bakeryClient.UnbakePi(cid.AsString())
		return apiv1.VMCID{}, err
	}

	err = c.bakeryClient.PowerCyclePi(cid.AsString())
	if err != nil {
		c.bakeryClient.UnbakePi(cid.AsString())
		return apiv1.VMCID{}, fmt.Errorf("Powering on failed. Rolled back deployment.")
	}
	return cid, nil
}
