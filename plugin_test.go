package qcollector_docker_events

import (
	"testing"

	"github.com/zpatrick/go-config"
	"github.com/stretchr/testify/assert"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/plugin"
)


func TestUnitNew(t *testing.T) {
	qChan := qtypes_qchannel.NewQChan()
	kv := map[string]string{"log.level": "trace"}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(kv)})
	b := qtypes_plugin.NewBase(qChan, cfg)
	exp := &Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(b, pluginTyp, pluginPkg, "plugin", version),
	}
	got, err := New(b, "plugin")
	assert.NoError(t, err, "No error expected here")
	assert.Equal(t, exp, got)
}