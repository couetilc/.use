package main

// I think now, some type of testing framework would be good
// - I want a small VM: alpine linux (musl and busybox)
//	- reset VM state then load different files/filesystem for each test
//	- assert against single process stdout/stderr
//	- I think Go unit tests will drive this.

import (
	"testing"
	// "log"
	// "fmt"
	"io"
	"github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
	"path/filepath"
	"io/fs"
	"compress/gzip"
	"archive/tar"
	"regexp"
	"os"
)

func init() {
	// I want to create (if not already created) the alpine linux container
	// for my tests. 
	// E.g. docker build -f Dockerfile.test github.com/couetilc/.use/testimage
	//	docker create github.com/couetilc/.use/testimage

	// so image options
	// - alpine:latest (around 5MB)
	// - busybox:musl (5x smaller than alpine)
	// let's start with busybox:musl, is smaller and I think enough.

	// all my test files should probably be named volumes
	// - So my unit test will choose what named volume to mount, then
	// assert against output

	// OK so 
	// - base image busybox:musl
	// - run tests as user "nobody"
	// - /home is test working directory (mount named volume here)

	// go client for docker https://pkg.go.dev/github.com/docker/docker/client
	// will drive all my unit tests.

}

func TestMain(t *testing.T) {
	var opt imageBuildOptions
	opt.includes = regexp.MustCompile(`^\.$|.*\.go|go\.mod|go\.sum`)
	opt.excludes = regexp.MustCompile(`\.git|\.envrc|bin`)
	opt.Dockerfile = "Dockerfile.test"
	opt.ContextFrom(os.Getenv("PWD"))

	// dc, err := client.NewClientWithOpts(client.FromEnv)
	// if err != nil { log.Fatalln(err) }
	// defer dc.Close()
	//
	// dc.ImageBuild(t.Context(), opt.Context, opt.ImageBuildOptions)

	t.Fatalf("TODO: run an alpine linux container")
	// // e.g.
	// // set CMD/ENTRYPOINT for image to be the test command
	// // MOUNT filesystem?
}

type imageBuildOptions struct {
	types.ImageBuildOptions
	// includes and excludes rules for adding to the build context
	includes *regexp.Regexp
	excludes *regexp.Regexp
}

func (opt *imageBuildOptions ) ContextFrom(rootDir string) {
	pr, pw := io.Pipe()
	opt.Context = pr

	go func() {
		defer pw.Close()
		gz := gzip.NewWriter(pw)
		defer gz.Close()
		tw := tar.NewWriter(gz)
		defer tw.Close()

		filepath.Walk(rootDir , func (path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// build context filepaths must be relative paths
			rel, err := filepath.Rel(os.Getenv("PWD"), path)
			if err != nil {
				return err
			}

			// establish which files should be in the build context
			if !opt.includes.MatchString(rel) || opt.excludes.MatchString(rel) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				// skip file
				return nil
			}

			hdr, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			hdr.Name = rel

			err = tw.WriteHeader(hdr)
			if err != nil {
				return err
			}

			if info.Mode().IsRegular() {
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()

				_, err = io.Copy(tw, f)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}()
}
