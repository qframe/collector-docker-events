package qtypes_plugin

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"github.com/pkg/errors"
	"github.com/zpatrick/go-config"

	"github.com/qframe/types/helper"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/ticker"
)

const (
	version = "0.2.0"
)

type Base struct {
	BaseVersion 	string
	Version 		string
	QChan 			qtypes_qchannel.QChan
	ErrChan			chan error
	Cfg 			*config.Config
	MyID			int
	Typ				string
	Pkg				string
	Name 			string
	LogOnlyPlugs 	[]string
	MsgCount		map[string]float64
	LocalCfg 		map[string]string

}

func NewBase(qChan qtypes_qchannel.QChan, cfg *config.Config) Base {
	b := Base{
		BaseVersion: version,
		QChan: qChan,
		ErrChan: make(chan error),
		Cfg: cfg,
		LogOnlyPlugs:   []string{},
		MsgCount:       map[string]float64{
			"received": 0.0,
			"loopDrop": 0.0,
			"inputDrop": 0.0,
			"successDrop": 0.0,
		},
	}
	b.LocalCfg, _  = cfg.Settings()
	logPlugs, err := cfg.String("log.only-plugins")
	if err == nil {
		b.LogOnlyPlugs = strings.Split(logPlugs, ",")
	}
	return b
}

func (p *Base) CfgString(path string) (string, error) {
	key := fmt.Sprintf("%s.%s.%s", p.Typ, p.Name, path)
	if res, ok := p.LocalCfg[key]; ok {
		return res, nil
	}
	if res, ok := p.LocalCfg[path]; ok {
		return res, nil
	}
	return "", errors.New("Could not find "+key)
}

func (p *Base) CfgStringOr(path, alt string) string {
	res, err := p.CfgString(path)
	if err != nil {
		return alt
	}
	return res
}

func (p *Base) CfgInt(path string) (int, error) {
	key := fmt.Sprintf("%s.%s.%s", p.Typ, p.Name, path)
	if res, ok := p.LocalCfg[key]; ok {
		return strconv.Atoi(res)
	}
	return 0, errors.New("Could not find "+key)
}

func (p *Base) CfgIntOr(path string, alt int) int {
	res, err := p.CfgInt(path)
	if err != nil {
		return alt
	}
	return res
}

func (p *Base) CfgBool(path string) (bool, error) {
	key := fmt.Sprintf("%s.%s.%s", p.Typ, p.Name, path)
	if res, ok := p.LocalCfg[key]; ok {
		switch res {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New(fmt.Sprintf("Key '%s' neither false not true, but %s: ", key, res))

		}
	}
	return false, errors.New("Could not find "+key)
}

func (p *Base) CfgBoolOr(path string, alt bool) bool {
	res, err := p.CfgBool(path)
	if err != nil {
		return alt
	}
	return res
}

func (p *Base) GetInputs() (res []string) {
	inStr, err := p.CfgString("inputs")
	if err == nil {
		res = strings.Split(inStr, ",")
	}
	return res
}

func (p *Base) GetCfgItems(key string) []string {
	inStr, err := p.CfgString(key)
	if err != nil {
		inStr = ""
	}
	return strings.Split(inStr, ",")
}

func (p *Base) Log(logLevel, msg string) {
	if len(p.LogOnlyPlugs) != 0 && ! qtypes_helper.IsItem(p.LogOnlyPlugs, p.Name) {
		return
	}
	// TODO: Setup in each Log() invocation seems rude
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	dL, _ := p.Cfg.StringOr("log.level", "info")
	dI := qtypes_helper.LogStrToInt(dL)
	lI := qtypes_helper.LogStrToInt(logLevel)
	lMsg := fmt.Sprintf("[%+6s] %15s Name:%-10s >> %s", strings.ToUpper(logLevel), p.Pkg, p.Name, msg)
	if lI == 0 {
		log.Panic(lMsg)
	} else if dI >= lI {
		log.Println(lMsg)
	}
}

func (p *Base) StartTicker(name string, durMs int) qtypes_ticker.Ticker {
	p.Log("info", fmt.Sprintf("Start ticker '%s' with duration of %dms", name, durMs))
	ticker := qtypes_ticker.NewTicker(name, durMs)
	go ticker.DispatchTicker(p.QChan)
	return ticker
}

func (b *Base) SendData(msg interface{}) {
	b.QChan.SendData(msg)
}

/*
func (p *Base) DispatchMsgCount() {

	tickMs := p.CfgIntOr("count-ticker-ms", 5000)
	p.Log("info", fmt.Sprintf("Dispatch goroutine to send MsgCount every %dms", tickMs))
	ticker := time.NewTicker(time.Duration(tickMs)*time.Millisecond).C
	pre := map[string]float64{}
	for {
		tick := <-ticker
		pre = p.SendMsgCount(tick, pre)
	}
}

func (p *Base) SendMsgCount(tick time.Time, pre map[string]float64) map[string]float64 {
	dims := map[string]string{
		"plugin_name": p.Name,
		"plugin_version": p.Version,
		"plugin_type": p.Typ,
	}
	qm := NewExt(p.Name, "none", Counter, 0.0, dims, tick, false)
	for k,v := range p.MsgCount {
		if _, ok := pre[k]; !ok {
			pre[k] = v
		} else if pre[k] == v {
			continue
		}
		qm.Name = fmt.Sprintf("msg.%s", k)
		p.Log("debug", fmt.Sprintf("Send MsgCount %s=%f", qm.Name,v))
		qm.Value = float64(v)
		p.QChan.SendData(qm)
	}
	return pre
}
*/