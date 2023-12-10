package models

type AdvancedForm struct {
	Type      string
	Provider  SbercloudProvider
	Vpc       SbercloudVpc
	VpcSubnet SbercloudVpcSubnet
}

type SbercloudVpc struct {
	ResourceName         string
	Name                 string
	Cidr                 string
	EnterpriseProject_id string
}

type SbercloudProvider struct {
	Region    string
	AccessKey string
	SecretKey string
}

type SbercloudVpcSubnet struct {
	Name      string
	Resource  string
	Cidr      string
	GatewayIp string
	VpcId     string
}
