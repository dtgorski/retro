// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"retro/emu"
	"retro/emu/config"
	"runtime"
)

type (
	// options are the CLI options and flags.
	options struct {
		otherConfigFile  *string
		image1FilePath   *string
		image2FilePath   *string
		cpuSpeedInMHz    *float64
		justPrintVersion *bool
		windowZoomLevel  *int
	}
)

var version = "dev"

func main() {
	runtime.GOMAXPROCS(1)

	log.SetPrefix("retro: ")
	log.SetFlags(0)

	version = fmt.Sprintf("%s (%s %s)", version, runtime.GOOS, runtime.GOARCH)

	// Load defaults, parse command line options and flags.
	conf := config.DefaultConfig
	opts := parseOpts()

	// Overwrite default config with command line options.
	applyOpts(conf, opts)

	conf.Version = version
	conf.Window.Title = fmt.Sprintf(conf.Window.Title, conf.CPU.MHz)

	// -v, prints version.
	if *opts.justPrintVersion {
		log.Printf(version)
		os.Exit(0)
	}

	// Emulator "OFF" button.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ... and go.
	if err := emu.Run(ctx, conf); err != nil {
		log.Fatal(err)
	}
}

func applyOpts(conf *config.Config, opts options) {

	// Overwrite configuration defaults with configuration file.
	if err := conf.FromYML(*opts.otherConfigFile); err != nil {
		log.Print(err)
		return
	}

	// Overwrite loaded disk config with provided image names.
	if *opts.image1FilePath != "" {
		conf.Disk.Drive1 = *opts.image1FilePath
	}
	if *opts.image2FilePath != "" {
		conf.Disk.Drive2 = *opts.image2FilePath
	}

	// Overwrite loaded config with command line options.
	if *opts.windowZoomLevel >= 1 && *opts.windowZoomLevel < 0x10 {
		conf.Window.Zoom = *opts.windowZoomLevel
	}

	// Overwrite CPU clock.
	if *opts.cpuSpeedInMHz >= 0.1 && *opts.cpuSpeedInMHz <= 20 {
		conf.CPU.MHz = *opts.cpuSpeedInMHz
	}
}

func parseOpts() options {
	flag.Usage = func() {
		writer := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(writer, help, version)
	}
	opts := options{
		otherConfigFile:  flag.String("c", "retro.config.yml", ""),
		image1FilePath:   flag.String("1", "", ""),
		image2FilePath:   flag.String("2", "", ""),
		cpuSpeedInMHz:    flag.Float64("m", 0.98, ""),
		justPrintVersion: flag.Bool("v", false, ""),
		windowZoomLevel:  flag.Int("z", 3, ""),
	}
	flag.Parse()
	return opts
}

// ---

var help = `
Usage: retro [options]
Retro Apple II Emulator %s

The default name of the configuration file is "retro.config.yml".
At program start, following configuration locations are checked:
first the directory, where the executable resides; then outside,
in ~/.config/retro/, /usr/local/etc/retro/, and /etc/retro/

When no disk image is provided (with -1), the emulator will boot
into the Applesoft BASIC prompt. Additionally, the ROM(s) of the
virtual Apple Disk II Interface will not be mounted into memory.

Command line options override their configuration counterparts. 

Options:
    -1 <path/to/image>
         The APPLE DISK II .dsk image to insert into Drive 1.
         Path can be a HTTP URL. Fetched images are not saved.
    
    -2 <path/to/image>
         The APPLE DISK II .dsk image to insert into Drive 2.
         Path can be a HTTP URL. Fetched images are not saved.

    -c <path/to/config>
         Path to an alternative configuration file. 

    -m <cpu-clock-in-mhz>
         MHz speed of the CPU clock. Usual clock settings are
         0.98 (1MHz) and 3.58 (4MHz). The first is the default.

    -z <window-zoom [1..n]>
         Window magnification. A zoom factor of 1 is equivalent
         to the Apple II native resolution of 280 x 192 pixels.
         Default value: 3 

More options:
    -h  Display this usage help and exit.
    -v  Print program version and exit.

Sources: <https://github.com/dtgorski/retro>

`
