package qtypes_docker_events

import (
	"github.com/docker/docker/api/types"
	"github.com/qframe/types/helper"
)

type NetworkEvent struct {
	DockerEvent
	Network 	types.NetworkResource
}

func NewNetworkEvent(de DockerEvent, net types.NetworkResource) NetworkEvent {
	return NetworkEvent{
		DockerEvent: de,
		Network: net,
	}
}

func (e *NetworkEvent) GetNetworkName() string {
	if e.Network.Name != "" {
		return e.Network.Name
	} else {
		return "<none>"
	}
}

// NetworkToJSON create a nested JSON object.
func (e *NetworkEvent) NetworkToJSON() (map[string]interface{}) {
	res := e.Base.ToJSON()
	res["msg_message"] = e.Message
	res["network"] = e.Network
	return res
}

// NetworkToFlatJSON create a key/val JSON map, which can be consumed by KSQL.
func (e *NetworkEvent) NetworkToFlatJSON() (res map[string]interface{}) {
	res = e.Base.ToFlatJSON()
	res["msg_message"] = e.Message
	res["engine_id"] = e.Engine.ID
	res["network_id"] = e.Network.ID
	res["network_driver"] = e.Network.Driver
	res["network_scope"] = e.Network.Scope
	res["network_created"] = e.Network.Created.String()
	rLab, err := qtypes_helper.PrefixFlatKV(e.Network.Labels, res, "network_label")
	if err == nil {
		res = rLab
	}
	return res
}
