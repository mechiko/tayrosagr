package process

import (
	"fmt"
	"tayrosagr/repo/configdb"
	"tayrosagr/repo/znakdb"
	"tayrosagr/utility"

	"github.com/mechiko/dbscan"
	"go.uber.org/zap"
)

// const startSSCC = "1462709225" // gs1 rus id zapivkom для памяти запивком

type Krinica struct {
	logger   *zap.SugaredLogger
	dbCfg    *configdb.DbConfig
	dbZnak   *znakdb.DbZnak
	gtin     string
	inn      string
	Sscc     []string
	Cis      []*utility.CisInfo
	Pallet   map[string][]*utility.CisInfo
	warnings []string
	errors   []string
	Koroba   map[string][]string
	CisAll   map[string]string
}

func New(dbs *dbscan.Dbs) (*Krinica, error) {
	logger, _ := zap.NewProduction()
	if dbs == nil {
		return nil, fmt.Errorf("отсутствует описатель бд")
	}
	cfgInfo := dbs.Info(dbscan.Config)
	if !cfgInfo.Exists {
		return nil, fmt.Errorf("отсутствует база конфиг.дб")
	}
	dbCfg, err := configdb.New(cfgInfo)
	if err != nil {
		return nil, fmt.Errorf(" %w", err)
	}
	znakInfo := dbs.Info(dbscan.TrueZnak)
	if !znakInfo.Exists {
		return nil, fmt.Errorf("отсутствует база чз")
	}
	dbZnak, err := znakdb.New(znakInfo)
	if err != nil {
		return nil, fmt.Errorf(" %w", err)
	}

	inn, err := dbCfg.Key("inn")
	if err != nil {
		return nil, fmt.Errorf("find inn error %w", err)
	}
	if inn == "" {
		return nil, fmt.Errorf("inn config.db empty error %w", err)
	}
	k := &Krinica{
		logger:   logger.Sugar(),
		dbCfg:    dbCfg,
		dbZnak:   dbZnak,
		inn:      inn,
		Pallet:   make(map[string][]*utility.CisInfo),
		warnings: make([]string, 0),
		errors:   make([]string, 0),
		Sscc:     make([]string, 0),
		Cis:      make([]*utility.CisInfo, 0),
	}
	return k, nil
}

func (k *Krinica) AddWarn(warn string) {
	k.warnings = append(k.warnings, warn)
}

func (k *Krinica) Warnings() []string {
	return k.warnings
}

func (k *Krinica) AddError(err string) {
	k.errors = append(k.errors, err)
}

func (k *Krinica) Errors() []string {
	return k.errors
}

func (k *Krinica) ResetPalletMap() {
	for key := range k.Pallet {
		delete(k.Pallet, key)
	}
}

func (k *Krinica) Reset() {
	k.ResetPalletMap()
	k.Cis = make([]*utility.CisInfo, 0)
	k.Sscc = make([]string, 0)
	k.errors = make([]string, 0)
	k.warnings = make([]string, 0)
}
