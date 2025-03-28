package build

import (
	"path/filepath"

	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/errors"
	helpers "github.com/npitsillos/spinit/helpers"
	"github.com/npitsillos/spinit/pkg/build"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	NAME_FLAG       = "name"
	TAG_FLAG        = "tag"
	DOCKERFILE_FLAG = "dockerfile"
	CONFIG_FLAG     = "config"
)

func newBuildOps(dir string, flagSet *pflag.FlagSet) (*build.BuildOpt, error) {
	dockerfile, _ := flagSet.GetString(DOCKERFILE_FLAG)
	name, _ := flagSet.GetString(NAME_FLAG)
	tag, _ := flagSet.GetString(TAG_FLAG)

	return &build.BuildOpt{
		ProjectDir: dir,
		Name:       name,
		Tag:        tag,
		Dockerfile: dockerfile,
	}, nil
}

func resolveProjectName(dir string) (string, error) {
	path, err := filepath.Abs(dir)
	return filepath.Base(path), err
}

func resolveDockerFile(path string) (string, error) {

	if helpers.FileExists(filepath.Join(path, "Dockerfile")) {
		return "Dockerfile", nil
	}
	if helpers.FileExists(filepath.Join(path, "dockerfile")) {
		return "dockerfile", nil
	}
	return "", errors.ErrNoDockerFile
}

func NewBuildCommand() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:     "build [WORKING_DIRECTORY]",
		Short:   "build application image",
		Long:    `Builds an application's Docker image and loads it in docker daemon.`,
		Args:    cobra.MaximumNArgs(1),
		PreRunE: resolveNameAndDockerfile,
		RunE: func(cmd *cobra.Command, args []string) error {
			buildOpts, err := newBuildOps(args[0], cmd.Flags())
			cobra.CheckErr(err)
			err = build.BuildImage(buildOpts)
			cobra.CheckErr(err)

			configFile, _ := cmd.Flags().GetString(CONFIG_FLAG)
			appConfig := config.NewAppConfig()
			appConfig.SetConfigFilePath(getConfigFilePath(args[0], configFile))
			appConfig.Build = determineBuildInfo(cmd.Flags())
			appConfig.WriteToFile()
			return nil
		},
	}

	buildCmd.Flags().StringP(NAME_FLAG, "n", "", "Image name. If not passed this is resolved from the project directory.")
	buildCmd.Flags().StringP(TAG_FLAG, "t", "latest", "Image tag")
	buildCmd.Flags().StringP(DOCKERFILE_FLAG, "d", "", "Path to dockerfile")
	buildCmd.Flags().StringP(CONFIG_FLAG, "c", "", "App config file")
	return buildCmd
}

func resolveNameAndDockerfile(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = append(args, helpers.GetWorkingDir())
	}
	nameFlag, _ := cmd.Flags().GetString(NAME_FLAG)
	if nameFlag == "" {
		nameFlag, err := resolveProjectName(args[0])
		if err != nil {
			return err
		}
		cmd.Flags().Set(NAME_FLAG, nameFlag)
	}
	dockerfile, _ := cmd.Flags().GetString(DOCKERFILE_FLAG)
	if dockerfile == "" {
		dockerfile, err := resolveDockerFile(args[0])
		if err != nil {
			return err
		}
		cmd.Flags().Set(DOCKERFILE_FLAG, dockerfile)
	}
	return nil
}

func getConfigFilePath(workDir, configFileOverride string) string {

	configFile := filepath.Join(workDir, config.DefaultAppConfigFile)
	if configFileOverride != "" {
		configFile = configFileOverride
	}

	return configFile
}

func determineBuildInfo(flagSet *pflag.FlagSet) *config.Build {
	dockerfile, _ := flagSet.GetString(DOCKERFILE_FLAG)
	version, _ := flagSet.GetString(TAG_FLAG)
	name, _ := flagSet.GetString(NAME_FLAG)

	return &config.Build{
		Image:      name,
		Dockerfile: dockerfile,
		Version:    version,
	}
}
