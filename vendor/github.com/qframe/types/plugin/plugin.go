package qtypes_plugin

import (
	"github.com/zpatrick/go-config"
	"github.com/qframe/types/qchannel"
)

type Plugin struct {
	Base
	MyID		int
	Typ			string
	Pkg			string
	Version 	string
	Name 		string
}


func NewNamedPlugin(qChan qtypes_qchannel.QChan, cfg *config.Config, typ, pkg, name, version string) *Plugin {
	b := NewBase(qChan, cfg)
	return NewBasePlugin(b, typ, pkg, name, version)
}


func NewBasePlugin(b Base, typ, pkg, name, version string) *Plugin {
	p := &Plugin{
		Base: b,
		Typ:   		typ,
		Pkg:  		pkg,
		Version:	version,
		Name: 		name,
	}
	return p
}

func (p *Plugin) GetInfo() (typ,pkg,name string) {
	return p.Typ, p.Pkg, p.Name
}

func (p *Plugin) GetLogOnlyPlugs() []string {
	return p.LogOnlyPlugs
}
