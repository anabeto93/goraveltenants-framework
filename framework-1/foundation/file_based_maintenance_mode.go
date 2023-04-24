package foundation

import (
	"encoding/json"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileBasedMaintenanceMode struct {
	app foundationcontract.Application
}

func NewFileBasedMaintenanceMode(app foundationcontract.Application) *FileBasedMaintenanceMode {
	return &FileBasedMaintenanceMode{
		app: app,
	}
}

func (f *FileBasedMaintenanceMode) Activate(payload map[string]interface{}) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.path(), data, 0644)
}

func (f *FileBasedMaintenanceMode) Deactivate() error {
	if f.Active() {
		return os.Remove(f.path())
	}
	return nil
}

func (f *FileBasedMaintenanceMode) Active() bool {
	_, err := os.Stat(f.path())
	return !os.IsNotExist(err)
}

func (f *FileBasedMaintenanceMode) Data() (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(f.path())
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (f *FileBasedMaintenanceMode) path() string {
	return filepath.Join(f.app.StoragePath(), "framework", "down")
}
