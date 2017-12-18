package main

type cpiContext struct {
	URL string
}

type stemcellCloudProps struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	Infrastructure string `json:"infrastructure"`
	DiskFormat     string `json:"disk_format"`
	OsType         string `json:"os_type"`
	OsDistro       string `json:"os_distro"`
	Architecture   string `json:"architecture"`
}
