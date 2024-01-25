package watch

import (
	"fmt"
	"gs/generate"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

const binaryName = "gs-watcher"

// Builder composes of both runner and watcher. Whenever watcher gets notified, builder starts a build process, and forces the runner to restart
type Builder struct {
	runner  *Runner
	watcher *Watcher
}

// NewBuilder constructs the Builder instance
func NewBuilder(w *Watcher, r *Runner) *Builder {
	return &Builder{watcher: w, runner: r}
}

// Build listens watch events from Watcher and sends messages to Runner
// when new changes are built.
func (b *Builder) Build() {
	go b.registerSignalHandler()
	go func() {
		b.watcher.update <- "init"
	}()

	for range b.watcher.Wait() {
		err := generate.Generate()
		if err != nil {
			log.Fatalln(err)
		}
		fileName := generateBinaryName(fmt.Sprintf("gs/%s", b.watcher.gsConfig.Module))
		log.Info("Building service")
		removeFile(fileName)
		// build package
		cmd, err := runCommand("go", "build", "-o", fileName, path.Join(b.watcher.gsConfig.Paths.Gen, "cmd", "local", "main.go"))
		if err != nil {
			log.Fatalf("Could not run 'go build' command: %s", err)
		}

		if err := cmd.Wait(); err != nil {
			if err := interpretError(err); err != nil {
				log.Println(fmt.Sprintf("An error occurred while building: %s", err))
			} else {
				log.Println("A build error occurred. Please update your code...", err)
			}

			continue
		}
		log.Info("Running service")
		// and start the new process
		b.runner.restart(fileName)
	}
}

func (b *Builder) registerSignalHandler() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signals
	b.watcher.Close()
	b.runner.Close()
}

// interpretError checks the error, and returns nil if it is
// an exit code 2 error. Otherwise error is returned as it is.
// when a compilation error occurres, it returns with code 2.
func interpretError(err error) error {
	exiterr, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}

	status, ok := exiterr.Sys().(syscall.WaitStatus)
	if !ok {
		return err
	}

	if status.ExitStatus() == 2 {
		return nil
	}

	return err
}

func generateBinaryPrefix() string {
	pth := os.Getenv("GOPATH")
	if pth != "" {
		return fmt.Sprintf("%s/bin/%s", pth, binaryName)
	}

	return pth
}

func generateBinaryName(packagePath string) string {
	packageName := strings.Replace(packagePath, "/", "-", -1)

	return fmt.Sprintf("%s-%s", generateBinaryPrefix(), packageName)
}

func removeFile(fileName string) {
	// check if file exists
	_, err := os.Stat(fileName)
	if err != nil {
		return
	}
	if fileName != "" {
		cmd := exec.Command("rm", fileName)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}
}
