package typicli

import (
	"bufio"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/typienv"
	"gopkg.in/urfave/cli.v1"
)

func (t *TypicalCli) updateTypical(ctx *cli.Context) {
	runOrFatal(goCommand(), "build", "-o", typienv.TypicalBinaryPath(), typienv.TypicalMainPackage())
}

func (t *TypicalCli) buildBinary(ctx *cli.Context) {
	runOrFatal(goCommand(), "build", "-o",
		typienv.BinaryPath(t.TypiApp.BinaryName),
		typienv.MainPackage(t.TypiApp.ApplicationPkg))
}

func (t *TypicalCli) runBinary(ctx *cli.Context) {
	// TODO: .env should be store in typienv
	err := setEnv(".env")
	if err != nil {
		fmt.Printf("Can't load .env:%s", err.Error())
		fmt.Printf("Typically generate new .env file, please check the environment variable in .env before run again")
		t.generateDotEnv(".env")
	} else {
		runOrFatal(typienv.BinaryPath(t.TypiApp.BinaryName), []string(ctx.Args())...)
	}

}

func (t *TypicalCli) runTest(ctx *cli.Context) {
	args := []string{"test"}
	args = append(args, t.TypiApp.TestTargets...)
	args = append(args, "-coverprofile=cover.out")
	runOrFatal(goCommand(), args...)
}

func (t *TypicalCli) releaseDistribution(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func (t *TypicalCli) generateMock(ctx *cli.Context) {
	runOrFatal(goCommand(), "get", "github.com/golang/mock/mockgen")
	for _, mockTarget := range t.MockTargets {
		dest := t.MockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]
		runOrFatal(goBinary("mockgen"),
			"-source", mockTarget,
			"-destination", dest,
			"-package", t.MockPkg)
	}
}

func (t *TypicalCli) generateDotEnv(filename string) (err error) {
	buf, err := os.Create(filename)
	if err != nil {
		return
	}

	// TODO: move it to constant
	template := `{{range .}}export {{usage_key .}}={{usage_default .}}
{{end}}`

	buf.WriteString("// .env automtically reload when using `./typicalw run`\n\n")
	envconfig.Usagef(t.TypiApp.ConfigPrefix, t.TypiApp.Config, buf, template)

	for i := range t.Modules {
		module := t.Modules[i]
		envconfig.Usagef(module.ConfigPrefix, module.Config, buf, template)
	}

	return
}

func (t *TypicalCli) appPath(name string) string {
	return fmt.Sprintf("./%s/%s", t.ApplicationPkg, name)
}

func setEnv(envfile string) (err error) {
	file, err := os.Open(envfile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "export") {
			args := strings.TrimSpace(text[len("export"):])
			pair := strings.Split(args, "=")
			if len(pair) > 1 {
				os.Setenv(pair[0], pair[1])
				log.Printf("Set Environment %s\n", pair[0])
			}
		}
	}
	return
}

func goBinary(name string) string {
	return fmt.Sprintf("%s/%s/%s", build.Default.GOPATH, "bin", name)
}

func goCommand() string {
	return fmt.Sprintf("%s/bin/go", build.Default.GOROOT)
}

func runOrFatal(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
