package main

import (
	"encoding/json"
	"fmt"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
	bakeryclient "github.com/vchrisr/bakery-client"
)

func (c CPI) UploadSettings(cid apiv1.VMCID, diskSettings, agentSettings []byte) error {
	err := c.bakeryClient.UploadBytesAsFile(cid.AsString(), "settings.json", agentSettings)
	if err == nil {
		return c.bakeryClient.UploadBytesAsFile(cid.AsString(), "disks.json", diskSettings)
	}

	return err
}

func (c CPI) GenerateNewSettings(pi bakeryclient.PiInfo, agentID apiv1.AgentID, cid apiv1.VMCID, networks apiv1.Networks, env apiv1.VMEnv) ([]byte, []byte, error) {
	ao, err := LoadConfig("/var/vcap/jobs/bakery_cpi/config/cpi.json")
	if err != nil {
		return nil, nil, err
	}

	ae := apiv1.NewAgentEnvFactory().ForVM(agentID, cid, networks, env, ao)

	ae.AttachEphemeralDisk("/dev/mapper/loop0")

	aeJson, _ := ae.AsBytes()
	disksBytes, _ := json.Marshal(pi.Disks[1:])

	return disksBytes, aeJson, nil
}

func (c CPI) RegenerateSettings(pi bakeryclient.PiInfo) ([]byte, []byte, error) {
	vmCID := apiv1.NewVMCID(pi.Id)
	settingsBytes, err := c.bakeryClient.DownloadFileAsBytes(vmCID.AsString(), "settings.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get settings.json from Pi: %v\n", err)
	}

	disksBytes, err := c.bakeryClient.DownloadFileAsBytes(vmCID.AsString(), "disks.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get disks.json from Pi: %v\n", err)
	}

	ae, err := apiv1.NewAgentEnvFactory().FromBytes(settingsBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not parse settings: %v\n", err)
	}

	for i, disk := range pi.Disks {
		if i <= 1 { //skip first 2 disks
			loopDevice := fmt.Sprintf("/dev/mapper/loop%v", i-1) //index 1 = loop0
			ae.AttachPersistentDisk(apiv1.NewDiskCID(disk.ID), loopDevice)
		}
	}

	settingsBytes, _ = ae.AsBytes()
	disksBytes, _ = json.Marshal(pi.Disks[2:])

	return disksBytes, settingsBytes, nil
}
