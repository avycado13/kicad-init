// Package kicad handles all utilities related to creating and modifying KiCad files
package kicad

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/avycado13/kicad-init/files"
)

type BoardDefaults struct {
	Zones map[string]any `json:"zones,omitempty"`
}

type DesignSettings struct {
	ZonesUseNoOutline *bool `json:"zones_use_no_outline,omitempty"`
}

type IPC2581Settings struct {
	Dist       string `json:"dist"`
	DistPN     string `json:"distpn"`
	InternalID string `json:"internal_id"`
	Mfg        string `json:"mfg"`
	MPN        string `json:"mpn"`
}

type BoardSettings struct {
	DesignSettings  DesignSettings   `json:"design_settings"`
	IPC2581         IPC2581Settings  `json:"ipc2581"`
	LayerPairs      []map[string]any `json:"layer_pairs"`
	LayerPresets    []map[string]any `json:"layer_presets"`
	Viewports       []map[string]any `json:"viewports"`
	ThreeDViewports []map[string]any `json:"3dviewports"`
}

type NetClass struct {
	Name             string  `json:"name"`
	BusWidth         float64 `json:"bus_width"`
	Clearance        float64 `json:"clearance"`
	DiffPairGap      float64 `json:"diff_pair_gap"`
	DiffPairViaGap   float64 `json:"diff_pair_via_gap"`
	DiffPairWidth    float64 `json:"diff_pair_width"`
	LineStyle        int     `json:"line_style"`
	MicroviaDiameter float64 `json:"microvia_diameter"`
	MicroviaDrill    float64 `json:"microvia_drill"`
	PCBColor         string  `json:"pcb_color"`
	Priority         int     `json:"priority"`
	SchematicColor   string  `json:"schematic_color"`
	TrackWidth       float64 `json:"track_width"`
	ViaDiameter      float64 `json:"via_diameter"`
	ViaDrill         float64 `json:"via_drill"`
	WireWidth        float64 `json:"wire_width"`
}

type NetSettings struct {
	Classes             []NetClass       `json:"classes"`
	Meta                map[string]any   `json:"meta"`
	NetColors           map[string]any   `json:"net_colors,omitempty"`
	NetclassAssignments map[string]any   `json:"netclass_assignments,omitempty"`
	NetclassPatterns    []map[string]any `json:"netclass_patterns"`
}

type LibrarySettings struct {
	PinnedFootprintLibs []string `json:"pinned_footprint_libs"`
	PinnedSymbolLibs    []string `json:"pinned_symbol_libs"`
}

func NewBoardDefaults() BoardDefaults {
	return BoardDefaults{}
}

func NewDesignSettings() DesignSettings {
	return DesignSettings{}
}

func NewIPC2581Settings() IPC2581Settings {
	return IPC2581Settings{}
}

func NewBoardSettings() BoardSettings {
	return BoardSettings{
		DesignSettings:  NewDesignSettings(),
		IPC2581:         NewIPC2581Settings(),
		LayerPairs:      []map[string]any{},
		LayerPresets:    []map[string]any{},
		Viewports:       []map[string]any{},
		ThreeDViewports: []map[string]any{},
	}
}

func NewNetClass() NetClass {
	return NetClass{
		Name:             "Default",
		BusWidth:         12.0,
		Clearance:        0.2,
		DiffPairGap:      0.25,
		DiffPairViaGap:   0.25,
		DiffPairWidth:    0.2,
		LineStyle:        0,
		MicroviaDiameter: 0.3,
		MicroviaDrill:    0.1,
		PCBColor:         "rgba(0, 0, 0, 0.000)",
		Priority:         2147483647,
		SchematicColor:   "rgba(0, 0, 0, 0.000)",
		TrackWidth:       0.2,
		ViaDiameter:      0.6,
		ViaDrill:         0.3,
		WireWidth:        6.0,
	}
}

func NewNetSettings() NetSettings {
	return NetSettings{
		Classes:          []NetClass{NewNetClass()},
		Meta:             map[string]any{"version": 4},
		NetclassPatterns: []map[string]any{},
	}
}

func NewLibrarySettings() LibrarySettings {
	return LibrarySettings{
		PinnedFootprintLibs: []string{},
		PinnedSymbolLibs:    []string{},
	}
}

type KiCadProjectSettings struct {
	BoardDefaults   BoardDefaults   `json:"board_defaults"`
	Board           BoardSettings   `json:"board"`
	NetSettings     NetSettings     `json:"net_settings"`
	LibrarySettings LibrarySettings `json:"library"`
}

func GenerateKiCadProjectSettings() KiCadProjectSettings {
	return KiCadProjectSettings{
		BoardDefaults:   NewBoardDefaults(),
		Board:           NewBoardSettings(),
		NetSettings:     NewNetSettings(),
		LibrarySettings: NewLibrarySettings(),
	}
}

func CreateKiCadProject(mgr *files.Manager) (err error) {
	settings := GenerateKiCadProjectSettings()

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	file, err := mgr.CreateFile(fmt.Sprintf("%s.kicad_pro", filepath.Base(filepath.Clean(mgr.RootDir))))
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	_, err = file.Write(data)
	return err
}
