package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"simple-app/pkg/flag"
	"simple-app/pkg/options"
)

type App struct {
	basename    string
	name        string
	description string
	options     options.CliOptions
	runFunc     RunFunc
	silence     bool
	noVersion   bool
	noConfig    bool
	//commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

func NewApp(name string,basename string,opts ...Option) *App {
	a := &App{
		name: name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

type Option func(*App)

func WithOptions(opt options.CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

type RunFunc func(basename string) error

func WithRunFunc(runc RunFunc ) Option{
	return func(a *App) {
		a.runFunc = runc
	}
}

func (a *App) Run() {
	fmt.Printf("%+v\n", a)
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v",err)
		os.Exit(1)
	}
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   a.basename,
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true

	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets flag.NamedFlagSets
	if a.options != nil {
		//a.options为接口类型，实现了Flags()方法
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))
	a.cmd = &cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}
	return nil
}