package typicli

import (
	"bufio"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func buildBinary(ctx *cli.Context) {
	runOrFatal(goCommand(), "build", "-o", "bin/typical-rest-server", "./cmd/app")
}

func runBinary(ctx *cli.Context) {
	setEnv(".env")
	fmt.Println(os.Getenv("TEST"))
	runOrFatal("bin/typical-rest-server", []string(ctx.Args())...)
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
				log.Printf("Set Environment %s=%s\n", pair[0], pair[1])
			}
		}
	}
	return
}

func runTest(ctx *cli.Context) {
	testTarget := []string{
		"./app/controller",
		"./app/repository",
	}

	args := []string{"test"}
	args = append(args, testTarget...)
	args = append(args, "-coverprofile=cover.out")
	runOrFatal(goCommand(), args...)
}

func vendoringDependency(ctx *cli.Context) {
	runOrFatal(goCommand(), "get", "github.com/golang/dep/cmd/dep")
	runOrFatal(goBinary("dep"), "ensure")
}

func releaseDistribution(ctx *cli.Context) {
	fmt.Println("Not implemented")
}

func generateMock(ctx *cli.Context) {
	runOrFatal(goCommand(), "get", "github.com/golang/mock/mockgen")

	// TODO: retrieve file based on pattern
	runOrFatal(goBinary("mockgen"), "-source=app/repository/book_repo.go", "-destination=mock/book_repo.go", "-package=mock")
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
