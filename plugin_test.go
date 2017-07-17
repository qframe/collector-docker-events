package qcollector_docker_events

import (
	"testing"

	"github.com/zpatrick/go-config"
	"github.com/stretchr/testify/assert"
	"github.com/qnib/qframe-types"
)


func TestNew(t *testing.T) {
	qChan := qtypes.NewQChan()
	kv := map[string]string{"log.level": "trace"}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(kv)})
	exp := &Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, "plugin", version),
	}
	got, err := New(qChan, cfg, "plugin")
	assert.NoError(t, err, "No error expected here")
	assert.Equal(t, exp, got)
}