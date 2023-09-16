package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
	"os"
)

func main() {
	var json *gjson.Json
	if gfile.Exists(".config.json") {
		bytes := gfile.GetBytes(".config.json")
		var err error
		json, err = gjson.DecodeToJson(bytes)
		if err != nil {
			panic(err)
		}
	} else {
		json = gjson.New(nil)
	}

	ctx := gctx.GetInitCtx()
	frontendPath := json.Get("frontend.path", "frontend").String()
	frontendInstall := json.Get("frontend.install", "pnpm install").String()
	frontendDev := json.Get("frontend.dev", "pnpm dev").String()
	frontendBuild := json.Get("frontend.build", "pnpm build").String()
	frontBuildPath := json.Get("frontend.buildPath", "").String()
	port := json.Get("front.port", 5173).Int()
	backendPath := json.Get("backend.path", ".").String()
	backendRun := json.Get("backend.dev", "go run main.go").String()
	backendBuild := json.Get("backend.build", "go build main.go").String()
	backendBuildPath := json.Get("backend.buildPath", "").String()
	backendOutName := json.Get("backend.outName", "main").String()
	switch os.Args[1] {
	case "dev":
		errs := make(chan error, 100)
		q := make(chan int, 100)
		exit := make(chan int, 100)

		runGO(ctx, backendPath, backendRun, exit, errs, q)

		runFrontend(ctx, frontendPath, frontendInstall, frontendDev, exit, errs, q, port)
		go func() {
			err := <-errs
			fmt.Println(err)
		}()
		go func() {
			o := <-q
			exit <- o
			exit <- o
			exit <- o
		}()
		gproc.AddSigHandlerShutdown(func(sig os.Signal) {
			exit <- 0
		})
		gproc.Listen()
	case "build":
		if err := buildGo(ctx, backendPath, backendBuild, backendBuildPath, backendOutName); err != nil {
			panic(err)
		}
		if err := buildFrontend(ctx, frontendPath, frontendInstall, frontendBuild, frontBuildPath); err != nil {
			panic(err)
		}
	}

}

func buildGo(ctx context.Context, backendPath string, build string, path string, outName string) error {
	cmd := gproc.NewProcessCmd("go mod tidy")
	cmd.Dir = backendPath
	err := cmd.Run(ctx)
	if err != nil {
		return err
	}
	s := ""

	if path != "" {
		if gstr.HasSuffix(path, "/") || gstr.HasSuffix(path, "\\") {
			p := gfile.Join(gfile.Pwd(), path)
			os.MkdirAll(p, os.ModePerm)
			s = fmt.Sprintf(" -o %s", gfile.Join(p, outName))
		} else {
			p := gfile.Join(gfile.Pwd(), gfile.Dir(path))
			os.MkdirAll(p, os.ModePerm)
			s = fmt.Sprintf(" -o %s", gfile.Join(gfile.Pwd(), path))
		}
	} else {
		s = fmt.Sprintf(" -o %s", outName)
	}
	processCmd := gproc.NewProcessCmd(build + s)
	processCmd.Dir = backendPath
	err = processCmd.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
func buildFrontend(ctx context.Context, frontendPath, install, build, path string) error {

	cmd := gproc.NewProcessCmd(install)
	cmd.Dir = frontendPath
	err := cmd.Run(ctx)
	if err != nil {
		return err
	}
	s := ""
	if path != "" {
		s = fmt.Sprintf(" --outDir %s", gfile.Join(gfile.Pwd(), path, "frontend"))
	}
	processCmd := gproc.NewProcessCmd(build + s)
	processCmd.Dir = frontendPath
	err = processCmd.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
func runGO(ctx context.Context, backendPath, backendRun string, exitSign chan int, errChan chan error, pexit chan int) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					errChan <- e

				}
			}
		}()
		cmd := gproc.NewProcessCmd("go mod tidy")
		cmd.Dir = backendPath
		err := cmd.Run(ctx)
		if err != nil {
			panic(err)
		}
		processCmd := gproc.NewProcessCmd(backendRun)
		processCmd.Dir = backendPath
		_, err = processCmd.Start(ctx)
		if err != nil {
			panic(err)
		}
		go func() {
			o := <-exitSign
			if o == 0 || o == 2 {
				err := cmd.Kill()
				if err != nil {
					errChan <- err
				}
			}

		}()
		err = processCmd.Wait()
		pexit <- 1

		if err != nil {
			panic(err)
		}
	}()
}
func runFrontend(ctx context.Context, frontendPath, frontendInstall, frontendDev string, exitSign chan int, errChan chan error, pexit chan int, port int) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					errChan <- e

				}
			}
		}()
		cmd := gproc.NewProcessCmd(frontendInstall)
		cmd.Dir = frontendPath
		err := cmd.Run(ctx)
		if err != nil {
			panic(err)
		}
		processCmd := gproc.NewProcessCmd(frontendDev + fmt.Sprintf(" --port %d --strictPort", port))
		processCmd.Dir = frontendPath
		_, err = processCmd.Start(ctx)
		if err != nil {
			panic(err)
		}
		go func() {
			o := <-exitSign
			if o == 0 || o == 2 {
				err := cmd.Kill()
				if err != nil {
					errChan <- err
				}
			}
		}()
		err = processCmd.Wait()
		pexit <- 2
		if err != nil {
			panic(err)
		}

	}()
}
