package xray

import (
	"bytes"

	"github.com/mhsanaei/3x-ui/v2/util/json_util"
)

// OutboundConfig represents an Xray outbound configuration.
// It defines how Xray sends outgoing connections including protocol, destination, and settings.
type OutboundConfig struct {
	Protocol       string               `json:"protocol"`
	Tag            string               `json:"tag"`
	Settings       json_util.RawMessage `json:"settings,omitempty"`
	StreamSettings json_util.RawMessage `json:"streamSettings,omitempty"`
	ProxySettings  json_util.RawMessage `json:"proxySettings,omitempty"`
	Mux            json_util.RawMessage `json:"mux,omitempty"`
}

// Equals compares two OutboundConfig instances for deep equality.
func (c *OutboundConfig) Equals(other *OutboundConfig) bool {
	if c.Protocol != other.Protocol {
		return false
	}
	if c.Tag != other.Tag {
		return false
	}
	if !bytes.Equal(c.Settings, other.Settings) {
		return false
	}
	if !bytes.Equal(c.StreamSettings, other.StreamSettings) {
		return false
	}
	if !bytes.Equal(c.ProxySettings, other.ProxySettings) {
		return false
	}
	if !bytes.Equal(c.Mux, other.Mux) {
		return false
	}
	return true
}
