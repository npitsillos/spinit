package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/appdefaults"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/npitsillos/spinit/errors"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/sync/errgroup"
)

var DOCKERFILE_FRONTEND = "dockerfile.v0"

type BuildOpt struct {
	ProjectDir string
	Name       string
	Tag        string
	Dockerfile string
	Load       bool
}

func BuildDockerImage(buildOpts *BuildOpt) error {
	ctx := appcontext.Context()

	c, err := client.New(ctx, appdefaults.Address)
	if err != nil {
		return err
	}
	pipeR, pipeW := io.Pipe()
	solveOpt, err := newSolveOpt(pipeW, buildOpts)
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

	if buildOpts.Load {
		eg.Go(func() error {
			if err := loadDockerTar(pipeR); err != nil {
				return err
			}
			return pipeR.Close()
		})
	} else {
		pipeR.Close()
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func newSolveOpt(w io.WriteCloser, buildOpts *BuildOpt) (*client.SolveOpt, error) {

	file := filepath.Join(buildOpts.ProjectDir, buildOpts.Dockerfile)

	mount, err := fsutil.NewFS(buildOpts.ProjectDir)
	if err != nil {
		return nil, errors.ErrInvalidBuildCtxLocalMount
	}

	dockerfileMount, err := fsutil.NewFS(filepath.Dir(file))
	if err != nil {
		return nil, errors.ErrInvalidDockerfileLocalMount
	}

	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: "docker", // TODO: use containerd image store when it is integrated to Docker
				Attrs: map[string]string{
					"name": fmt.Sprintf("%s:%s", buildOpts.Name, buildOpts.Tag),
				},
				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return w, nil
				},
			},
		},
		LocalMounts: map[string]fsutil.FS{
			"context":    mount,
			"dockerfile": dockerfileMount,
		},
		Frontend: DOCKERFILE_FRONTEND,
		FrontendAttrs: map[string]string{
			"filename": filepath.Base(file),
		},
	}, nil
}

func loadDockerTar(r io.Reader) error {
	cmd := exec.Command("docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
