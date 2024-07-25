package metadata

import (
	"fmt"
	"os"
	"time"

	"github.com/fioncat/ks/pkg/utils"
	"gopkg.in/yaml.v3"
)

type History struct {
	Records []*HistoryRecord

	Path string
}

type HistoryRecord struct {
	Timestamp int64  `yaml:"timestamp" json:"timestamp"`
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`

	// TODO: Add more info?
}

func loadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Path: path}, nil
		}
		return nil, fmt.Errorf("read history file: %w", err)
	}

	var records []*HistoryRecord
	err = yaml.Unmarshal(data, &records)
	if err != nil {
		return nil, fmt.Errorf("parse history file %q: %w (HINT: you might need to remove bad history file manually)", path, err)
	}

	return &History{
		Records: records,
		Path:    path,
	}, nil
}

func (h *History) Add(name, namespace string) {
	record := &HistoryRecord{
		Name:      name,
		Namespace: namespace,
		Timestamp: time.Now().Unix(),
	}
	// The latest record should be at the beginning of the history
	h.Records = append([]*HistoryRecord{record}, h.Records...)
}

func (h *History) GetLastKubeConfig(current string) *string {
	for _, record := range h.Records {
		if record.Name != current {
			return &record.Name
		}
	}
	return nil
}

func (h *History) GetLastNamespace(name, currentNamespace string) *string {
	for _, record := range h.Records {
		if record.Namespace == "" {
			continue
		}
		if record.Name == name && record.Namespace != currentNamespace {
			return &record.Namespace
		}
	}
	return nil
}

func (h *History) ClearKubeConfig(name string) {
	newRecords := make([]*HistoryRecord, 0)
	for _, record := range h.Records {
		if record.Name == name {
			continue
		}
		newRecords = append(newRecords, record)
	}
	h.Records = newRecords
}

func (h *History) ClearAll() {
	h.Records = nil
}

func (h *History) Save() error {
	data, err := yaml.Marshal(h.Records)
	if err != nil {
		return err
	}

	err = utils.WriteFile(h.Path, data)
	if err != nil {
		return fmt.Errorf("write history file: %w", err)
	}

	return nil
}
