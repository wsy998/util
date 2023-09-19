package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"util/internal/model"
	"util/internal/service"
)

func RunDev(ctx context.Context, config *model.TypeConfig) error {
	if config.Backend.Mirror != "" {
		origin := genv.Get("GOPROXY").String()
		err := genv.Set("GOPROXY", config.Backend.Mirror)
		if err != nil {
			return err
		}
		defer func() {
			if origin != "" {
				if err := genv.Set("GOPROXY", origin); err != nil {
					return
				}
			} else {
				if err := genv.Remove("GOPROXY"); err != nil {
					return
				}
			}
		}()
	}
	errChan := make(chan error, 10)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
		}()
		if err := service.Os().Run(ctx, "go mod tidy", config.Backend.Path); err != nil {
			errChan <- err
			return
		}
		if err := service.Os().Run(ctx, config.Backend.Dev, config.Backend.Path); err != nil {
			errChan <- err
			return
		}
	}()
	go func() {
		defer func() {
			wg.Done()
		}()
		if err := service.Os().Run(ctx, config.Frontend.Install, config.Frontend.Path); err != nil {
			errChan <- err
			return
		}
		if err := service.Os().Run(ctx, config.Frontend.Dev, config.Frontend.Path); err != nil {
			errChan <- err
			return
		}
	}()
	wg.Wait()
	close(errChan)
	var err error
	for e := range errChan {
		if e != nil {
			err = e
			break
		}
	}
	return err
}
func RunBuild(ctx context.Context, config *model.TypeConfig) error {
	if config.Backend.Mirror != "" {
		origin := genv.Get("GOPROXY").String()
		err := genv.Set("GOPROXY", config.Backend.Mirror)
		if err != nil {
			return err
		}
		defer func() {
			if origin != "" {
				if err := genv.Set("GOPROXY", origin); err != nil {
					return
				}
			} else {
				if err := genv.Remove("GOPROXY"); err != nil {
					return
				}
			}
		}()
	}
	errChan := make(chan error, 10)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
		}()
		if err := service.Os().Run(ctx, "go mod tidy", config.Backend.Path); err != nil {
			errChan <- err
			return
		}
		if err := service.Os().Run(ctx, config.Backend.Build, config.Backend.Path); err != nil {
			errChan <- err
			return
		}
	}()
	go func() {
		defer func() {
			wg.Done()
		}()
		if err := service.Os().Run(ctx, config.Frontend.Install, config.Frontend.Path); err != nil {
			errChan <- err
			return
		}
		if err := service.Os().Run(ctx, config.Frontend.Build, config.Frontend.Path); err != nil {
			errChan <- err
			return
		}
	}()
	wg.Wait()
	close(errChan)
	var err error
	for e := range errChan {
		if e != nil {
			err = e
			break
		}
	}
	return err
}
func RunGo(ctx context.Context, config *model.TypeConfig) error {
	s, ok := config.Backend.Alias[os.Args[2]]
	if ok {
		cmd := []string{s}
		cmd = append(cmd, os.Args[3:]...)
		if err := service.Os().Run(ctx, gstr.Join(cmd, " "), config.Backend.Path); err != nil {
			return err
		}
	} else {
		if err := service.Os().Run(ctx, gstr.Join(os.Args[2:], " "), config.Backend.Path); err != nil {
			return err
		}
	}
	return nil
}
func RunFe(ctx context.Context, config *model.TypeConfig) error {
	s, ok := config.Frontend.Alias[os.Args[2]]
	if ok {
		cmd := []string{s}
		cmd = append(cmd, os.Args[3:]...)
		if err := service.Os().Run(ctx, gstr.Join(cmd, " "), config.Frontend.Path); err != nil {
			return err
		}
	} else {

		if err := service.Os().Run(ctx, gstr.Join(os.Args[2:], " "), config.Frontend.Path); err != nil {
			return err
		}
	}
	return nil
}
func RunAlias(ctx context.Context, config *model.TypeConfig) error {
	s, exist := config.Alias[os.Args[1]]
	if !exist {
		fmt.Println("Don't found this Command!")
		return nil
	}
	err := service.Os().Run(ctx, s, gfile.Pwd())
	if err != nil {
		return err
	}
	return nil
}
