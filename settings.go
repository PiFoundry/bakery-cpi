package main

import (
	"encoding/json"
	"fmt"

	bakeryclient "github.com/PiFoundry/bakery-client"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
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

func (c CPI) RegenerateSettings(vmCID apiv1.VMCID) ([]byte, []byte, error) {
	pi, err := c.bakeryClient.GetPi(vmCID.AsString())
	if err != nil {
		return nil, nil, fmt.Errorf("Could not find pi with id: %v. %v", vmCID.AsString(), err)
	}

	settingsBytes, err := c.bakeryClient.DownloadFileAsBytes(vmCID.AsString(), "settings.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get settings.json from Pi: %v\n", err)
	}

	ae, err := apiv1.NewAgentEnvFactory().FromBytes(settingsBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not parse settings: %v\n", err)
	}

	//attach persistent disks to agent settings
	for i, disk := range pi.Disks {
		if i >= 2 { //skip system and disk
			loopDevice := fmt.Sprintf("/dev/mapper/loop%v", i-1) //index 1 = loop0
			diskCID := apiv1.NewDiskCID(disk.ID)
			ae.DetachPersistentDisk(diskCID) //Just detach everything before attaching so we don't need a different regen func fr detach
			ae.AttachPersistentDisk(diskCID, loopDevice)
		}
	}

	settingsBytes, _ = ae.AsBytes()
	disksBytes, _ := json.Marshal(pi.Disks[1:]) //skip first disk, its not a disk actually

	return disksBytes, settingsBytes, nil
}
