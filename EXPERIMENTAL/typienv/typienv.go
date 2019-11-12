package typienv

var (
	appVar        = EnvVar{"TYPICAL_APP", "app"}
	binVar        = EnvVar{"TYPICAL_BIN", "bin"}
	cmdVar        = EnvVar{"TYPICAL_CMD", "cmd"}
	mockVar       = EnvVar{"TYPICAL_MOCK", "mock"}
	releaseVar    = EnvVar{"TYPICAL_RELEASE", "release"}
	buildToolVar  = EnvVar{"TYPICAL_BUILD_TOOL", "build-tool"}
	dependencyVar = EnvVar{"TYPICAL_DEPENDENCY", "dependency"}
	metadataVar   = EnvVar{"TYPICAL_METADATA", ".typical-metadata"}
	readmeVar     = EnvVar{"TYPICAL_README", "README.md"}
)

var (
	App        *applicationFolder
	BuildTool  *applicationFolder
	Dependency *applicationFolder
	Bin        string
	Metadata   string
	Mock       string
	Release    string
	AppName    string
	Readme     string
)

type applicationFolder struct {
	Package string
	SrcPath string
	BinPath string
}

func init() {
	AppName = appVar.Value()
	cmd := cmdVar.Value()
	Bin = binVar.Value()
	Metadata = metadataVar.Value()
	buildTool := buildToolVar.Value()
	dependency := dependencyVar.Value()
	App = &applicationFolder{
		Package: "main",
		SrcPath: cmd + "/" + AppName,
		BinPath: Bin + "/" + AppName,
	}
	BuildTool = &applicationFolder{
		Package: "main",
		SrcPath: cmd + "/" + buildTool,
		BinPath: Bin + "/" + buildTool,
	}
	Dependency = &applicationFolder{
		Package: dependency,
		SrcPath: cmd + "/internal/" + dependency,
	}
	Mock = mockVar.Value()
	Release = releaseVar.Value()
	Readme = readmeVar.Value()
}
