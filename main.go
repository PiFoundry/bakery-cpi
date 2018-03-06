package main

import (
	"fmt"
	"os"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	bakeryclient "github.com/vchrisr/bakery-client"
)

type CPIFactory struct{}

type CPI struct {
	bakeryClient *bakeryclient.Client
}

var _ apiv1.CPIFactory = CPIFactory{}
var _ apiv1.CPI = CPI{}

func main() {
	logger := boshlog.NewLogger(boshlog.LevelNone)

	cli := rpc.NewFactory(logger).NewCLI(CPIFactory{})

	err := cli.ServeOnce()
	if err != nil {
		logger.Error("main", "Serving once: %s", err)
		os.Exit(1)
	}
}

func (f CPIFactory) New(context apiv1.CallContext) (apiv1.CPI, error) {
	var parsedContext cpiContext
	err := context.As(&parsedContext)
	if err != nil {
		return CPI{}, err
	}

	return CPI{
		bakeryClient: bakeryclient.New(parsedContext.URL),
	}, nil
}

func (c CPI) Info() (apiv1.Info, error) {
	return apiv1.Info{
		StemcellFormats: []string{"piimg"},
	}, nil
}

func (c CPI) CreateStemcell(imagePath string, cp apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	var scProps stemcellCloudProps
	err := cp.As(&scProps)
	if err != nil {
		return apiv1.StemcellCID{}, err
	}

	name := fmt.Sprintf("%v-%v", scProps.Name, scProps.Version)

	cid, err := c.bakeryClient.UploadImage(imagePath, name)
	if err != nil {
		if strings.Contains(err.Error(), "403") { //when running bosh upload-stemcell --fix the existing image wmust be overwritten. Bakery returns 403 in that case so we delete it and then re-upload it.
			c.bakeryClient.DeleteImage(name)
			cid, err = c.bakeryClient.UploadImage(imagePath, name)
			if err != nil {
				return apiv1.StemcellCID{}, err
			}
		} else {
			return apiv1.StemcellCID{}, err
		}
	}

	return apiv1.NewStemcellCID(cid), nil
}

func (c CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	return c.bakeryClient.DeleteImage(cid.AsString())
}

func (c CPI) DeleteVM(cid apiv1.VMCID) error {
	return c.bakeryClient.UnbakePi(cid.AsString())
}

func (c CPI) CalculateVMCloudProperties(res apiv1.VMResources) (apiv1.VMCloudProps, error) {
	return apiv1.NewVMCloudPropsFromMap(map[string]interface{}{}), nil
}

func (c CPI) SetVMMetadata(cid apiv1.VMCID, metadata apiv1.VMMeta) error {
	return nil
}

func (c CPI) HasVM(cid apiv1.VMCID) (bool, error) {
	return c.bakeryClient.IsPiBaked(cid.AsString())
}

func (c CPI) RebootVM(cid apiv1.VMCID) error {
	return c.bakeryClient.PowerCyclePi(cid.AsString())
}

func (c CPI) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {
	disks, err := c.bakeryClient.GetDisks()
	if err != nil {
		return []apiv1.DiskCID{}, err
	}

	diskCids := make([]apiv1.DiskCID, len(disks.Disks))
	i := 0
	for _, disk := range disks.Disks {
		diskCids[i] = apiv1.NewDiskCID(disk.ID)
		i++
	}

	return diskCids, nil
}

func (c CPI) CreateDisk(size int,
	cloudProps apiv1.DiskCloudProps, associatedVMCID *apiv1.VMCID) (apiv1.DiskCID, error) {
	cid, err := c.bakeryClient.CreateDisk(size)
	if err != nil {
		return apiv1.DiskCID{}, err
	}

	return apiv1.NewDiskCID(cid), nil
}

func (c CPI) DeleteDisk(cid apiv1.DiskCID) error {
	return c.bakeryClient.DeleteDisk(cid.AsString())
}

func (c CPI) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	err := c.bakeryClient.AttachDisk(vmCID.AsString(), diskCID.AsString())
	if err != nil {
		return fmt.Errorf("Could not attach disk to Pi: %v", err)
	}

	ds, as, err := c.RegenerateSettings(vmCID)
	if err != nil {
		return fmt.Errorf("Unable to regenerate settings: %v", err)
	}

	return c.UploadSettings(vmCID, ds, as)
}

func (c CPI) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	err := c.bakeryClient.DetachDisk(vmCID.AsString(), diskCID.AsString())
	if err != nil {
		return fmt.Errorf("Could not deattach disk from Pi: %v", err)
	}

	ds, as, err := c.RegenerateSettings(vmCID)
	if err != nil {
		return fmt.Errorf("Unable to regenerate settings: %v", err)
	}

	return c.UploadSettings(vmCID, ds, as)
}

func (c CPI) HasDisk(cid apiv1.DiskCID) (bool, error) {
	return c.bakeryClient.DiskExists(cid.AsString())
}
