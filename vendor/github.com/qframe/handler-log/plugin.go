package qhandler_log

import (
	"fmt"
	"github.com/zpatrick/go-config"

	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/plugin"
	"github.com/qframe/types/constants"
	"github.com/qframe/types/messages"
	"reflect"
	"github.com/qframe/cache-inventory"
	"strings"
)

const (
	version   = "0.2.0"
	pluginTyp = qtypes_constants.HANDLER
	pluginPkg = "log"
)

type Plugin struct {
	*qtypes_plugin.Plugin
}

func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	p := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	p.Version = version
	p.Name = name
	return p, nil
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start log handler v%s", p.Version))
	bg := p.QChan.Data.Join()
	for {
		select {
		case val := <-bg.Read:
			switch val.(type) {
			case qtypes_messages.Message:
				qm := val.(qtypes_messages.Message)
				if qm.StopProcessing(p.Plugin, false) {
					continue
				}
				p.Log("info" , fmt.Sprintf("%-15s : %s", strings.Join(qm.SourcePath, "->"), qm.Message))
			case qcache_inventory.ContainerRequest:
				continue
			default:
				p.Log("info", fmt.Sprintf("Dunno type '%s': %v", reflect.TypeOf(val), val))
			}
		}
	}
}
