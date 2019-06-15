package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	handler "simpleRestCache/pkg/srcctl/grpcclient"
)

type control struct {
	handler Handler
}

// Handler asks a service and output to a console
type Handler interface {
	All()
	TopN(n int)
	LastN(n int)
	Refresh()
	Clean()
	Settings()
}

func main() {
	// parse arguments
	// h := NewHandler()
	fs := flag.NewFlagSet("simpleNews", flag.ExitOnError)
	var (
		port = fs.Int("p", 8081, "A control port of a simpleRestCache instance")
		host = fs.String("h", "srcsvc", "A address of a simpleRestCache instance")
	)
	fs.Parse(os.Args[1:])

	// trim arguments
	arr := os.Args[1:]
	h := find(arr, "-h")
	p := find(arr, "-p")

	max, min := p, h
	if p < h {
		max, min = h, p
	}
	if max != len(arr) {
		arr = arr[max+2:]
	} else if min != len(arr) {
		arr = arr[min+2:]
	}

	addr := *host + ":" + strconv.Itoa(*port)
	c := &control{
		handler: handler.New(addr),
	}
	// parse commands
	c.rootPath(arr)
}

// =============ROOT===============
// parse commands
func (c *control) rootPath(arr []string) {
	if len(arr) == 0 {
		c.usageRoot()
		os.Exit(0)
	}
	switch arr[0] {
	case "stat":
		c.statPath(arr[1:])
	case "cache":
		c.cachePath(arr[1:])
	case "settings":
		c.settingsPath(arr[1:])
	default:
		c.usageRoot()
		os.Exit(0)
	}
}

func (c *control) usageRoot() {
	fmt.Println("Usage: \t srcctl COMMAND")
	fmt.Println("Commands:")
	fmt.Println("\tstat\tDisplay statistic of cache usage")
	fmt.Println("\tcache\tManage cache")
	fmt.Println("\tsettings\tDisplay settings of a cache system")
}

// =============STAT===============
func (c *control) statPath(arr []string) {
	if len(arr) == 0 {
		c.usageStat()
		os.Exit(0)
	}
	switch arr[0] {
	case "top":
		c.topPath(arr[1:])
	case "last":
		c.lastPath(arr[1:])
	case "all":
		c.allPath(arr[1:])
	default:
		c.usageStat()
		os.Exit(0)
	}
}

func (c *control) usageStat() {
	fmt.Println("Usage: \t srcctl stat COMMAND")
	fmt.Println("Commands:")
	fmt.Println("\tall\t\tDisplay all from cache. <N> number")
	fmt.Println("\ttop <N>\t\tDisplay top <N> popular requests to cache. <N> number")
	fmt.Println("\tlast <N>\t\tDisplay last <N> unpopular requests to cache. <N> number")
}

// =============TOP===============
func (c *control) topPath(arr []string) {
	if len(arr) == 0 || len(arr) > 1 {
		c.usageTop()
		os.Exit(0)
	}

	n, err := strconv.Atoi(arr[0])
	if err != nil {
		c.usageTop()
		os.Exit(0)
	}

	c.handler.TopN(n)
}

func (c *control) usageTop() {
	fmt.Println("Usage: \t srcctl stat top <N>")
	fmt.Println("\tDisplay top <N> popular requests to cache. <N> number")
}

// =============LAST===============
func (c *control) lastPath(arr []string) {
	if len(arr) == 0 || len(arr) > 1 {
		c.usageTop()
		os.Exit(0)
	}

	n, err := strconv.Atoi(arr[0])
	if err != nil {
		c.usageTop()
		os.Exit(0)
	}

	c.handler.LastN(n)
}

func (c *control) usageLast() {
	fmt.Println("Usage: \t srcctl stat last <N>")
	fmt.Println("\tDisplay last <N> unpopular requests to cache. <N> number")
}

// =============CACHE===============
func (c *control) cachePath(arr []string) {
	if len(arr) == 0 {
		c.usageCache()
		os.Exit(0)
	}
	switch arr[0] {
	case "all":
		c.allPath(arr[1:])
	case "clean":
		c.cleanPath(arr[1:])
	case "refresh":
		c.refreshPath(arr[1:])
	default:
		c.usageCache()
		os.Exit(0)
	}
}

func (c *control) usageCache() {
	fmt.Println("Usage: \t srcctl cache COMMAND")
	fmt.Println("Commands:")
	fmt.Println("\tall\t\tDisplay all from cache. <N> number")
	fmt.Println("\tclean\t\tDelete all cache records")
	fmt.Println("\trefresh\t\tRefresh all cache records")
}

// =============ALL===============
func (c *control) allPath(arr []string) {
	if len(arr) != 0 {
		c.usageAll()
		os.Exit(0)
	}

	c.handler.All()
}

func (c *control) usageAll() {
	fmt.Println("Usage: \t srcctl stat all")
	fmt.Println("\tDisplay all requests stored in cache")
}

// =============REFRESH===============
func (c *control) refreshPath(arr []string) {
	if len(arr) != 0 {
		c.usageRefresh()
		os.Exit(0)
	}

	c.handler.Refresh()
}

func (c *control) usageRefresh() {
	fmt.Println("Usage: \t srcctl cache refresh")
}

// =============CLEAN===============
func (c *control) cleanPath(arr []string) {
	if len(arr) != 0 {
		c.usageClean()
		os.Exit(0)
	}

	c.handler.Clean()
}

func (c *control) usageClean() {
	fmt.Println("Usage: \t srcctl cache clean")
}

// =============SETTINGS===============
func (c *control) settingsPath(arr []string) {
	if len(arr) != 0 {
		c.usageSettings()
		os.Exit(0)
	}

	c.handler.Settings()

}

func (c *control) usageSettings() {
	fmt.Println("Usage: \t srcctl settings")
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}
