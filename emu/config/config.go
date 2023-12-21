// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type (
	// Config is the main configuration structure.
	Config struct {
		Version string
		Window  `yaml:"window"`
		CPU     `yaml:"cpu"`
		Disk    `yaml:"disk"`
		Render  `yaml:"render"`
	}

	// Window ...
	Window struct {
		Title string `yaml:"title"`
		Zoom  int    `yaml:"zoom"`
	}

	// CPU ...
	CPU struct {
		MHz float64 `yaml:"mhz"`
	}

	// Disk ...
	Disk struct {
		Drive1 string `yaml:"drive-1"`
		Drive2 string `yaml:"drive-2"`
	}

	// Render ...
	Render struct {
		Mono  `yaml:"mono"`
		LoRes `yaml:"lores"`
		HiRes `yaml:"hires"`
	}

	// Mono ...
	Mono struct {
		Color int `yaml:"color"`
	}

	// LoRes ...
	LoRes struct {
		Colors []int `yaml:"colors"`
	}

	// HiRes ...
	HiRes struct {
		Colors []int `yaml:"colors"`
	}
)

// DefaultConfig is a working configuration.
var DefaultConfig = &Config{

	Window: Window{
		Title: "RETRO Apple II    │    6502 @ %.2f MHz    │    RESET:  CTRL + SHIFT + R",

		// As the resolution of an Apple II is 280x192px,
		// the native presentation mode may be too small.
		Zoom: 3,
	},
	CPU: CPU{
		// Usual clock settings are 0.98 and 3.58 (aka 1MHz and 4MHz).
		MHz: 0.98,
	},

	// File paths of "inserted" Disk 1 and Disk 2 images.
	// Paths can be HTTP URLs. Fetched images are not saved.
	// Using the -1 and -2 options overrides this setting.
	Disk: Disk{},

	Render: Render{
		Mono: Mono{
			// The color of the monochrome text.
			Color: 0x00B500FF,
		},
		LoRes: LoRes{
			Colors: []int{
				0x000000FF, // 0 Black
				0x9A2110FF, // 1 Magenta
				0x3C22A5FF, // 2 Dark Blue
				0xA13FB7FF, // 3 Purple
				0x07653EFF, // 4 Dark Green
				0x7B7E80FF, // 5 Dark Gray
				0x308FE3FF, // 6 Medium Blue
				0xADD8E6FF, // 7 Light Blue
				0x8B4513FF, // 8 Brown
				0xC77028FF, // 9 Orange
				0x7B7E80FF, // A Gray
				0xF39AC2FF, // B Pink
				0x2FB81FFF, // C Green
				0xFFFF80FF, // D Yellow
				0x6EE1C0FF, // E Aqua
				0xFFFFFFFF, // F White
			},
		},
		HiRes: HiRes{
			Colors: []int{
				0x000000FF, // 0 Black
				0xA13FB7FF, // 2 Purple
				0x2FB81FFF, // 1 Green
				0xFFFFFFFF, // 3 White
				0x000000FF, // 4 Black
				0x308FE3FF, // 6 Med Blue
				0xC77028FF, // 5 Orange
				0xFFFFFFFF, // 7 White
			},
		},
	},
}

// FromYML fills the Config object from a YAML file.
// Overwrites existing values in this configuration.
func (c *Config) FromYML(name string) error {
	path := c.findConfig(name)

	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config not loadable: %w", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return fmt.Errorf("syntax error %s: %w", path, err)
	}

	// Ensure list sizes, when config is incomplete.
	lores := &c.Render.LoRes
	lores.Colors = append(lores.Colors, make([]int, 0x10)...)
	lores.Colors = lores.Colors[:0x10]

	hires := &c.Render.HiRes
	hires.Colors = append(hires.Colors, make([]int, 0x08)...)
	hires.Colors = hires.Colors[:0x08]

	return nil
}

func (*Config) findConfig(name string) string {
	exec, _ := os.Executable()
	self := filepath.Dir(exec)

	locations := []string{
		"%s",
		fmt.Sprintf("%s/%%s", self),
		"~/.config/retro/%s",
		"/usr/local/etc/retro/%s",
		"/etc/retro/%s",
	}

	// For development, when executable resides in /bin.
	if strings.HasSuffix(self, "/bin") {
		locations = append([]string{"../%s"}, locations...)
	}

	for _, location := range locations {
		path := fmt.Sprintf(location, name)

		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		return path
	}
	return name
}
