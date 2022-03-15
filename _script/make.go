package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	flags := flag.NewFlagSet("build tool", flag.ExitOnError)
	all := flags.Bool("all", false, "build for all supported target")
	target := flags.String("target", runtime.GOOS, "set target for build, not used when -all is set")
	arch := flags.String("arch", runtime.GOARCH, "set arch for build, not used when -all is set")
	flags.Parse(os.Args[1:])
	if *all {
		for _, target := range [...]string{"windows", "linux", "freebsd", "darwin"} {
			for _, arch := range [...]string{"386", "amd64", "arm64", "arm"} {
				if err := build(target, arch); err != nil {
					log.Fatalln(err.Error())
				}
			}
		}
	} else {
		if err := build(*target, *arch); err != nil {
			log.Fatalln(err.Error())
		}
	}
}

const ldflags = "-s -w -buildid="

func build(target, arch string) error {
	out := "bin/access-network-" + target + "-" + arch
	if arch == "arm" {
		out += "7"
	}
	if target == "windows" {
		out += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", out, "-trimpath", "-ldflags", ldflags)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
