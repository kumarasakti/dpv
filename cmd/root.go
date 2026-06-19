package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kumarasakti/dpv/internal/color"
	"github.com/kumarasakti/dpv/internal/docker"
	"github.com/kumarasakti/dpv/internal/filter"
	"github.com/kumarasakti/dpv/internal/formatter"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var version = "dev"

const banner = ` ██████╗ ██████╗ ██╗   ██╗
 ██╔══██╗██╔══██╗██║   ██║
 ██║  ██║██████╔╝██║   ██║
 ██║  ██║██╔═══╝ ╚██╗ ██╔╝
 ██████╔╝██║      ╚████╔╝
 ╚═════╝ ╚═╝       ╚═══╝`

var bannerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

var (
	flagAll     bool
	flagSlim    bool
	flagInclude string
	flagOrder   string
	flagReverse bool
	flagJSON    bool
)

var rootCmd = &cobra.Command{
	Use:   "dpv [search]",
	Short: "Docker Pretty View -- a prettier docker ps",
	Long: bannerStyle.Render(banner) + "\n" +
		color.Bold.Render(" Docker Pretty View") +
		color.Dim.Render(" — a prettier docker ps\n") +
		"\n Displays Docker container info in a vertical, colored format.\n" +
		" Pass a comma-separated search to filter by container name.",
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print dpv version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(bannerStyle.Render(banner))
		fmt.Printf(" %s  %s\n", color.Bold.Render("Docker Pretty View"), color.Dim.Render(versionLabel(version)))
		fmt.Printf(" %s\n\n", color.Dim.Render("https://github.com/kumarasakti/dpv"))
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&flagAll, "all", "a", false, "Include stopped containers")
	rootCmd.Flags().BoolVarP(&flagSlim, "slim", "s", false, "Show slim minimal output")
	rootCmd.Flags().StringVarP(&flagInclude, "include", "i", "", "Columns to show: (n)id, (i)mage, co(m)mand, (c)reated, (s)tatus, (p)orts, (h)ealth")
	rootCmd.Flags().StringVarP(&flagOrder, "order", "o", "", "Sort by: name, created, status (default)")
	rootCmd.Flags().BoolVarP(&flagReverse, "reverse", "r", false, "Reverse display order")
	rootCmd.Flags().BoolVarP(&flagJSON, "json", "j", false, "Output as JSON")

	rootCmd.AddCommand(versionCmd)
	rootCmd.SetVersionTemplate(bannerStyle.Render(banner) + "\n" +
		" Docker Pretty View  {{.Version}}\n\n")
}

// Execute is the CLI entry point called from main.
func Execute(v string) {
	version = v
	rootCmd.Version = versionLabel(v)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := docker.NewDockerClient()
	if err != nil {
		return fmt.Errorf("cannot connect to Docker: %w", err)
	}
	defer client.Close()

	allContainers, err := client.ListContainers(ctx, flagAll)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	searchTerms := parseSearch(args)

	filtered := filter.Apply(allContainers, filter.Options{
		SearchTerms: searchTerms,
		ShowAll:     flagAll,
		OrderBy:     flagOrder,
		Reverse:     flagReverse,
	})

	runningCount := 0
	for _, c := range allContainers {
		if c.Running {
			runningCount++
		}
	}
	stats := docker.Stats{
		Total:       len(allContainers),
		Running:     runningCount,
		Matched:     len(filtered),
		HasSearch:   len(searchTerms) > 0,
		SearchTerms: searchTerms,
		ShowAll:     flagAll,
	}

	f := pickFormatter()
	return formatter.WriteTo(f, filtered, stats)
}

func parseSearch(args []string) []string {
	if len(args) == 0 || args[0] == "" {
		return nil
	}
	var terms []string
	for _, t := range strings.Split(args[0], ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			terms = append(terms, t)
		}
	}
	return terms
}

func pickFormatter() formatter.Formatter {
	if flagJSON {
		return &formatter.JSONFormatter{}
	}
	inc := formatter.ParseIncludes(flagInclude, flagSlim)
	if flagSlim {
		return &formatter.SlimFormatter{Includes: inc}
	}
	return &formatter.PrettyFormatter{Includes: inc}
}

// versionLabel ensures the version string is always displayed with a single "v" prefix.
func versionLabel(v string) string {
	if strings.HasPrefix(v, "v") {
		return v
	}
	return "v" + v
}
