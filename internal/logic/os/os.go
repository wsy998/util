package os

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"

	"util/internal/service"
)

func init() {
	service.RegisterOs(newSOs())
}

type sOs struct {
	manager *gproc.Manager
}

func newSOs() *sOs {
	return &sOs{manager: gproc.NewManager()}
}

func (receiver *sOs) Run(ctx context.Context, arg string, workSpace string) error {
	cmd := gproc.NewProcessCmd(arg)
	cmd.Dir = workSpace
	pid, err := cmd.Start(ctx)
	if err != nil {
		return err
	}
	receiver.manager.AddProcess(pid)
	g.Log().Infof(ctx, `Run the command "%s" in the "%s" directory with Pid %d.`, cmd.String(), workSpace, pid)
	err = cmd.Wait()
	receiver.manager.RemoveProcess(pid)
	if err != nil {
		return err
	}
	return nil
}
func (receiver sOs) KillAll() error {
	return receiver.manager.KillAll()
}
