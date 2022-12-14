package commons

import "crypto/rsa"

type Address string
type Package string
type MainSku string
type CountryCode string
type Percentile float32

type ApiResponseParameters struct {
	MainSku MainSku `json:"main_sku"`
}

type IpRequestAccessFormat string
type UseIpCachesFlag bool

type IpStackResponse struct {
	CountryCode CountryCode `json:"country_code"`
}

type ConfigurationRequirement struct {
	Package     Package     `json:"package"`
	CountryCode CountryCode `json:"country_code"`
}

type ConfigurationChance struct {
	PercentileMin Percentile `json:"percentile_min"`
	PercentileMax Percentile `json:"percentile_max"`
	MainSku       MainSku    `json:"main_sku"`
}

type ConfigurationRule struct {
	ConfigurationRequirement
	ConfigurationChance
}

type ConfigurationTable []ConfigurationRule

type UpdaterEnvironment struct {
	JwtPublicKey *rsa.PublicKey
}
