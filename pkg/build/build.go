package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/cmd/buildctl/build"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/npitsillos/spinit/errors"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/sync/errgroup"
)

var (
	DOCKERFILE_FRONTEND = "dockerfile.v0"
	DEFAULT_PLATFORMS   = "platform=linux/amd64,linux/arm64"
	IMAGE_TYPE          = "oci"
)

type BuildOpt struct {
	ProjectDir string
	Name       string
	Tag        string
	Dockerfile string
}

func BuildImage(buildOpts *BuildOpt) error {
	ctx := appcontext.Context()

	c, err := client.New(ctx, appdefaults.Address)
	if err != nil {
		return err
	}

	solveOpt, err := newSolveOpt(buildOpts)
	if err != nil {
		return err
	}
	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = c.Solve(ctx, nil, *solveOpt, ch)
		return err
	})
	eg.Go(func() error {
		d, err := progressui.NewDisplay(os.Stderr, progressui.TtyMode)
		if err != nil {
			// If an error occurs while attempting to create the tty display,
			// fallback to using plain mode on stdout (in contrast to stderr).
			d, _ = progressui.NewDisplay(os.Stdout, progressui.PlainMode)
		}
		// not using shared context to not disrupt display but let is finish reporting errors
		_, err = d.UpdateFrom(context.TODO(), ch)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func newSolveOpt(buildOpts *BuildOpt) (*client.SolveOpt, error) {

	file := filepath.Join(buildOpts.ProjectDir, buildOpts.Dockerfile)

	mount, err := fsutil.NewFS(buildOpts.ProjectDir)
	if err != nil {
		return nil, errors.ErrInvalidBuildCtxLocalMount
	}

	dockerfileMount, err := fsutil.NewFS(filepath.Dir(file))
	if err != nil {
		return nil, errors.ErrInvalidDockerfileLocalMount
	}

	exportEntry, err := createClientExportEntry(buildOpts)
	if err != nil {
		return nil, err
	}

	frontendAttrs, err := build.ParseOpt([]string{DEFAULT_PLATFORMS})
	if err != nil {
		return nil, err
	}
	frontendAttrs["filename"] = filepath.Base(file)

	return &client.SolveOpt{
		Exports: exportEntry,
		LocalMounts: map[string]fsutil.FS{
			"context":    mount,
			"dockerfile": dockerfileMount,
		},
		Frontend:      DOCKERFILE_FRONTEND,
		FrontendAttrs: frontendAttrs,
	}, nil
}

func createClientExportEntry(buildOpts *BuildOpt) ([]client.ExportEntry, error) {
	exports := fmt.Sprintf("type=%s,name=%s,dest=%s.tar", IMAGE_TYPE, buildOpts.Name, buildOpts.Name)
	return build.ParseOutput([]string{exports})
}
