package foundation

import (
	"encoding/json"
	"fmt"
	"github.com/goravel/framework/contracts/filesystem"
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

type PackageManifest struct {
	Filesystem   filesystem.Driver
	BasePath     string
	ManifestPath string
	VendorPath   string
	Manifest     map[string]interface{}
}

func NewPackageManifest(fs filesystem.Driver, basePath, manifestPath string) *PackageManifest {
	vendorPath := os.Getenv("GOLANG_VENDOR_DIR")
	if vendorPath == "" {
		vendorPath = "vendor"
	}
	vendorPath = filepath.Join(basePath, vendorPath)
	return &PackageManifest{
		Filesystem:   fs,
		BasePath:     basePath,
		ManifestPath: manifestPath,
		VendorPath:   vendorPath,
		Manifest:     make(map[string]interface{}),
	}
}

func (pm *PackageManifest) config(key string) (interface{}, error) {
	var result []interface{}

	for _, manifestValue := range pm.Manifest {
		switch manifestType := manifestValue.(type) {
		case map[string]interface{}:
			if value, ok := manifestType[key]; ok {
				result = append(result, value)
			}
		}
	}

	return result, nil
}

func (pm *PackageManifest) Providers() ([]foundationcontract.ServiceProvider, error) {
	providers, err := pm.config("providers")
	if err != nil {
		return nil, err
	}

	if providersSlice, ok := providers.([]interface{}); ok {
		result := make([]foundationcontract.ServiceProvider, len(providersSlice))
		for i, provider := range providersSlice {
			if sp, ok := provider.(foundationcontract.ServiceProvider); ok {
				result[i] = sp
			} else {
				return nil, fmt.Errorf("provider is not of type foundationcontract.ServiceProvider")
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("config(\"providers\") did not return a slice of interface{}")
}

func (pm *PackageManifest) Aliases() (map[string]foundationcontract.ServiceProvider, error) {
	aliases, err := pm.config("aliases")
	if err != nil {
		return nil, err
	}

	if aliasesMap, ok := aliases.(map[string]interface{}); ok {
		result := make(map[string]foundationcontract.ServiceProvider)
		for alias, sp := range aliasesMap {
			if spTyped, ok := sp.(foundationcontract.ServiceProvider); ok {
				result[alias] = spTyped
			} else {
				return nil, fmt.Errorf("alias value is not of type foundationcontract.ServiceProvider")
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("config(\"aliases\") did not return a map[string]interface{}")
}

func (pm *PackageManifest) build() error {
	// intialize the providers and aliases of the Manifest
	pm.Manifest["providers"] = []foundationcontract.ServiceProvider{}
	pm.Manifest["aliases"] = map[string]foundationcontract.ServiceProvider{}

	err := filepath.Walk(pm.VendorPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == pm.VendorPath {
			return nil
		}

		if info.IsDir() || !strings.HasSuffix(path, "_manifest.so") {
			return filepath.SkipDir
		}

		_, err = pm.Filesystem.Get(path)
		if err != nil {
			return filepath.SkipDir // the application should not break because some 3rd party developer did a bad job
		}

		if strings.HasSuffix(path, "_manifest.so") {
			// load the manifest file here and the defined service providers
			manifestJsonPath := filepath.Join(filepath.Dir(path), "manifest.json")
			if !pm.Filesystem.Exists(manifestJsonPath) {
				return filepath.SkipDir // the two are required for lookups
			}
			declarations, manErr := pm.Filesystem.Get(manifestJsonPath) // containing packages and aliases
			if manErr != nil {
				return manErr // this is worth addressing
			}
			var declarationContents struct {
				Packages struct {
					Providers []string          `json:"provider"`
					Aliases   map[string]string `json:"aliases"`
				} `json:"packages"`
			}
			manErr = json.Unmarshal([]byte(declarations), &declarationContents)
			if manErr != nil {
				return manErr
			}

			// now use plugin to lookup these declared providers and aliases
			openedPlugin, pErr := plugin.Open(path)
			if pErr != nil {
				return pErr
			}
			for _, provider := range declarationContents.Packages.Providers {
				spImpl, pErr := openedPlugin.Lookup(provider)
				if pErr != nil {
					return fmt.Errorf("error loading %s in %s : %v", provider, path, pErr)
				}
				spFin, ok := spImpl.(foundationcontract.ServiceProvider)
				if !ok {
					return fmt.Errorf("%s does not implement foundation.ServiceProvider", provider)
				}
				currentProviders := pm.Manifest["providers"].([]foundationcontract.ServiceProvider)
				currentProviders = append(currentProviders, spFin)
				pm.Manifest["providers"] = currentProviders
			}

			// moving to the aliases
			for alias, provider := range declarationContents.Packages.Aliases {
				// first check that alias is found in the list of declared providers or skip it
				if !pm.arrayContainsString(declarationContents.Packages.Providers, provider) {
					continue
				}

				spImpl, pErr := openedPlugin.Lookup(provider)
				if pErr != nil {
					return fmt.Errorf("error loading %s in %s : %v", alias, path, pErr)
				}
				spFin, ok := spImpl.(foundationcontract.ServiceProvider)
				if !ok {
					return fmt.Errorf("%s does not implement foundation.ServiceProvider", provider)
				}
				currentAliases := pm.Manifest["aliases"].(map[string]foundationcontract.ServiceProvider)
				currentAliases[alias] = spFin
				pm.Manifest["aliases"] = currentAliases
			}
		}

		return nil
	})
	return err
}

func (pm *PackageManifest) arrayContainsString(array []string, target string) bool {
	for _, val := range array {
		if val == target {
			return true
		}
	}
	return false
}

func (pm *PackageManifest) GetManifest() (map[string]interface{}, error) {
	if pm.Manifest != nil {
		return pm.Manifest, nil
	}

	if len(pm.Manifest) > 0 {
		return pm.Manifest, nil
	}

	if _, err := os.Stat(pm.ManifestPath); os.IsNotExist(err) {
		err := pm.build()
		if err != nil {
			return nil, err
		}
	}

	return pm.Manifest, nil
}
