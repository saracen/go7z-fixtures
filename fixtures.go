package fixtures

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type ReadSeekerCloser interface {
	io.ReaderAt
	io.ReadSeeker
	io.Closer
}

type Fixture struct {
	Archive string
	Tags    Tags

	Name    string
	Size    int64

	Entries int

	ReadSeekerCloser
}

type Tags []string

func (tags Tags) Has(tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (tags Tags) HasAny(any []string) bool {
	for _, t := range any {
		if tags.Has(t) {
			return true
		}
	}
	return false
}

var allFixtures = []Fixture{
	{ Archive: "executables", Entries: 10, Tags: []string{"executable", "bcj2", "386", "amd64", "arm", "ppc"} },
	{ Archive: "executables-bcj2-386-amd64", Entries: 2, Tags: []string{"executable", "bcj2", "386", "amd64"} },
	{ Archive: "executables-bcj2", Entries: 10, Tags: []string{"executable", "bcj2", "386", "amd64"} },
	{ Archive: "bzip2", Entries: 1, Tags: []string{"random", "bzip2"} },
	{ Archive: "copy", Entries: 1, Tags: []string{"random", "copy"} },
	{ Archive: "deflate", Entries: 1, Tags: []string{"random", "deflate"} },
	{ Archive: "delta-lzma", Entries: 2, Tags: []string{"random", "delta", "lzma"} },
	{ Archive: "delta", Entries: 2, Tags: []string{"random", "delta"} },
	{ Archive: "empty", Entries: 0, Tags: []string{"empty"} },
	{ Archive: "ppmd-bzip2-deflate-copy", Entries: 2, Tags: []string{"ppmd", "bzip2", "deflate", "copy" } },
	{ Archive: "ppmd", Entries: 1, Tags: []string{"ppmd" } },
}

var fixtureDir string

func init() {
	pkgs, err := packages.Load(nil, "github.com/saracen/go7z-fixtures")
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		if len(pkg.GoFiles) == 0 {
			panic("cannot find go7z-fixtures package on system")
		}

		fixtureDir = filepath.Dir(pkg.GoFiles[0])
		return
	}
}

func Fixtures(include []string, exclude []string) (fixtures []Fixture, closeAll io.Closer) {
	for _, fixture := range allFixtures {
		if !fixture.Tags.HasAny(include) {
			continue
		}

		if fixture.Tags.HasAny(exclude) {
			continue
		}
		
		name := filepath.Join(fixtureDir, "testdata", "archives", fixture.Archive + ".7z")
		f, err := os.Open(name)
		if err != nil {
			panic("fixture not found")
		}

		fixture.ReadSeekerCloser = f
		fixture.Name = name

		fi, err := f.Stat()
		if err != nil {
			panic(err)
		}

		fixture.Size = fi.Size()
		fixtures = append(fixtures, fixture)
	}

	return fixtures, &allcloser{fixtures}
}

type allcloser struct {
	fixtures []Fixture
}

func (close *allcloser) Close() error {
	var first error
	for _, c := range close.fixtures {
		err := c.Close()
		if err != nil && first == nil {
			first = err
		}
	}
	return first
}
