package deploy

import "strings"

type ResourceType string

const (
	ResourceTypeDeployment ResourceType = "DEPLOYMENT"
)

var AllResourceTypes = []ResourceType{
	ResourceTypeDeployment,
}

func (rt ResourceType) String() string {
	return string(rt)
}

func (rt ResourceType) Upper() ResourceType {
	return ResourceType(strings.ToUpper(rt.String()))
}

func (rt ResourceType) IsDeployment() bool {
	return rt.Upper() == ResourceTypeDeployment
}

func (rt ResourceType) IsSupported() bool {
	for _, at := range AllResourceTypes {
		if at == rt.Upper() {
			return true
		}
	}

	return false
}
