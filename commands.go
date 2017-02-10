package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/command/agent"
	"github.com/hashicorp/consul/command/base"
	"github.com/hashicorp/consul/version"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &agent.Command{
				Revision:          version.GitCommit,
				Version:           version.Version,
				VersionPrerelease: version.VersionPrerelease,
				HumanVersion:      version.GetHumanVersion(),
				Ui:                ui,
				ShutdownCh:        make(chan struct{}),
			}, nil
		},

		"configtest": func() (cli.Command, error) {
			return &command.ConfigTestCommand{
				Command: base.Command{
					Flags: base.FlagSetNone,
					Ui:    ui,
				},
			}, nil
		},

		"event": func() (cli.Command, error) {
			return &command.EventCommand{
				Ui: ui,
			}, nil
		},

		"exec": func() (cli.Command, error) {
			return &command.ExecCommand{
				ShutdownCh: makeShutdownCh(),
				Ui:         ui,
			}, nil
		},

		"force-leave": func() (cli.Command, error) {
			return &command.ForceLeaveCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"kv": func() (cli.Command, error) {
			return &command.KVCommand{
				Ui: ui,
			}, nil
		},

		"kv delete": func() (cli.Command, error) {
			return &command.KVDeleteCommand{
				Ui: ui,
			}, nil
		},

		"kv get": func() (cli.Command, error) {
			return &command.KVGetCommand{
				Ui: ui,
			}, nil
		},

		"kv put": func() (cli.Command, error) {
			return &command.KVPutCommand{
				Ui: ui,
			}, nil
		},

		"kv export": func() (cli.Command, error) {
			return &command.KVExportCommand{
				Ui: ui,
			}, nil
		},

		"kv import": func() (cli.Command, error) {
			return &command.KVImportCommand{
				Ui: ui,
			}, nil
		},

		"join": func() (cli.Command, error) {
			return &command.JoinCommand{
				Ui: ui,
			}, nil
		},

		"keygen": func() (cli.Command, error) {
			return &command.KeygenCommand{
				Ui: ui,
			}, nil
		},

		"keyring": func() (cli.Command, error) {
			return &command.KeyringCommand{
				Ui: ui,
			}, nil
		},

		"leave": func() (cli.Command, error) {
			return &command.LeaveCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"lock": func() (cli.Command, error) {
			return &command.LockCommand{
				ShutdownCh: makeShutdownCh(),
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"maint": func() (cli.Command, error) {
			return &command.MaintCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"members": func() (cli.Command, error) {
			return &command.MembersCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"monitor": func() (cli.Command, error) {
			return &command.MonitorCommand{
				ShutdownCh: makeShutdownCh(),
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"operator": func() (cli.Command, error) {
			return &command.OperatorCommand{
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"info": func() (cli.Command, error) {
			return &command.InfoCommand{
				Ui: ui,
			}, nil
		},

		"reload": func() (cli.Command, error) {
			return &command.ReloadCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"rtt": func() (cli.Command, error) {
			return &command.RTTCommand{
				Command: base.Command{
					Flags: base.FlagSetClientHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"snapshot": func() (cli.Command, error) {
			return &command.SnapshotCommand{
				Ui: ui,
			}, nil
		},

		"snapshot restore": func() (cli.Command, error) {
			return &command.SnapshotRestoreCommand{
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"snapshot save": func() (cli.Command, error) {
			return &command.SnapshotSaveCommand{
				Command: base.Command{
					Flags: base.FlagSetHTTP,
					Ui:    ui,
				},
			}, nil
		},

		"snapshot inspect": func() (cli.Command, error) {
			return &command.SnapshotInspectCommand{
				Command: base.Command{
					Flags: base.FlagSetNone,
					Ui:    ui,
				},
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				HumanVersion: version.GetHumanVersion(),
				Ui:           ui,
			}, nil
		},

		"watch": func() (cli.Command, error) {
			return &command.WatchCommand{
				ShutdownCh: makeShutdownCh(),
				Ui:         ui,
			}, nil
		},
	}
}

// makeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// interrupt or SIGTERM received.
func makeShutdownCh() <-chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()

	return resultCh
}
