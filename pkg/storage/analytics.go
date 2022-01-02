package storage

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v2"
)

type Analytics struct {
	InstallID            string    `json:"install_id"`
	RunID                string    `json:"run_id"`
	Version              string    `json:"version"`
	Timestamp            time.Time `json:"timestamp"`
	UploadIndex          int       `json:"upload_index"`
	GOOS                 string    `json:"goos"`
	GOARCH               string    `json:"goarch"`
	GoVersion            string    `json:"go_version"`
	MemAlloc             int       `json:"mem_alloc"`
	MemTotalAlloc        int       `json:"mem_total_alloc"`
	MemSys               int       `json:"mem_sys"`
	MemNumGC             int       `json:"mem_num_gc"`
	BadgerMain           int       `json:"badger_main"`
	BadgerTrees          int       `json:"badger_trees"`
	BadgerDicts          int       `json:"badger_dicts"`
	BadgerDimensions     int       `json:"badger_dimensions"`
	BadgerSegments       int       `json:"badger_segments"`
	ControllerIndex      int       `json:"controller_index" kind:"cumulative"`
	ControllerComparison int       `json:"controller_comparison" kind:"cumulative"`
	ControllerDiff       int       `json:"controller_diff" kind:"cumulative"`
	ControllerIngest     int       `json:"controller_ingest" kind:"cumulative"`
	ControllerRender     int       `json:"controller_render" kind:"cumulative"`
	SpyRbspy             int       `json:"spy_rbspy" kind:"cumulative"`
	SpyPyspy             int       `json:"spy_pyspy" kind:"cumulative"`
	SpyGospy             int       `json:"spy_gospy" kind:"cumulative"`
	SpyEbpfspy           int       `json:"spy_ebpfspy" kind:"cumulative"`
	SpyPhpspy            int       `json:"spy_phpspy" kind:"cumulative"`
	SpyDotnetspy         int       `json:"spy_dotnetspy" kind:"cumulative"`
	SpyJavaspy           int       `json:"spy_javaspy" kind:"cumulative"`
	AppsCount            int       `json:"apps_count"`
}

const analyticsKey = "analytics"

func (s *Storage) SaveAnalytics(a *Analytics) error {
	v, err := json.Marshal(a)
	if err != nil {
		return err
	}
	err = s.main.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(analyticsKey), v)
		err := txn.SetEntry(entry)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) LoadAnalytics() (*Analytics, error) {
	a := &Analytics{}
	err := s.main.View(func(txn *badger.Txn) error {
		v, err := txn.Get([]byte(analyticsKey))
		if err != nil {
			return err
		}
		err = v.Value(func(val []byte) error {
			err = json.Unmarshal(val, a)
			if err != nil {
				return err
			}
			return nil
		})
		return nil
	})
	return a, err
}
