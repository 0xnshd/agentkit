//go:build mage

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"charm.land/huh/v2"
)

func Lint(ctx context.Context) error {
	var target string
	theme := huh.ThemeFunc(func(isDark bool) *huh.Styles {
		return huh.ThemeDracula(true)
	})

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Lint scope").
				Description("Choose what you want to lint.").
				Options(
					huh.NewOption("Recent Changes", "recent-changes"),
					huh.NewOption("Package", "package-search"),
					huh.NewOption("Full project", "full-project"),
				).
				Value(&target),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		return err
	}

	switch target {
	case "recent-changes":
		return lintRecentChanges(ctx)
	case "package-search":
		return lintPackageSearch(ctx, theme)
	case "full-project":
		return lintFullProject(ctx)
	default:
		fmt.Println(target)
		return nil
	}
}

func lintRecentChanges(ctx context.Context) error {
	if !hasHead(ctx) {
		fmt.Println("No commits found yet.")
		return nil
	}

	if !hasParentCommit(ctx) {
		fmt.Println("No parent commit found, linting full project instead.")
		return runGolangCILint(ctx, "./...")
	}

	fmt.Println("Linting new issues for recent changes")
	return runGolangCILintArgs(ctx, "run", "--build-tags", "mage", "--new-from-rev=HEAD~1")
}

func lintFullProject(ctx context.Context) error {
	fmt.Println("Linting full project")
	return runGolangCILint(ctx, "./...")
}

func lintPackageSearch(ctx context.Context, theme huh.Theme) error {
	packages, err := listGoPackages(ctx)
	if err != nil {
		return err
	}
	if len(packages) == 0 {
		fmt.Println("No Go packages found.")
		return nil
	}

	var query string
	queryForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Package search").
				Description("Filter packages by name or path.").
				Value(&query),
		),
	).WithTheme(theme)
	if err := queryForm.Run(); err != nil {
		return err
	}

	matches := filterPackages(packages, query)
	if len(matches) == 0 {
		fmt.Printf("No Go packages matched %q.\n", query)
		return nil
	}

	var selected string
	options := make([]huh.Option[string], 0, len(matches))
	for _, pkg := range matches {
		options = append(options, huh.NewOption(pkg, pkg))
	}

	selectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select package").
				Description(fmt.Sprintf("%d matching package(s)", len(matches))).
				Options(options...).
				Value(&selected),
		),
	).WithTheme(theme)
	if err := selectForm.Run(); err != nil {
		return err
	}

	fmt.Printf("Linting package: %s\n", selected)
	return runGolangCILint(ctx, selected)
}

func listGoPackages(ctx context.Context) ([]string, error) {
	root, err := repoRoot(ctx)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "go", "list", "-tags", "mage", "-f", "{{.Dir}}", "./...")
	cmd.Dir = root
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("go list failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, err
	}

	text := strings.TrimSpace(string(output))
	if text == "" {
		return nil, nil
	}

	dirs := strings.Split(text, "\n")
	packages := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		rel, err := filepath.Rel(root, strings.TrimSpace(dir))
		if err != nil {
			return nil, err
		}
		if rel == "." {
			packages = append(packages, ".")
			continue
		}
		packages = append(packages, "./"+filepath.ToSlash(rel))
	}

	sort.Strings(packages)
	return packages, nil
}

func filterPackages(packages []string, query string) []string {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return packages
	}

	matches := make([]string, 0, len(packages))
	for _, pkg := range packages {
		if strings.Contains(strings.ToLower(pkg), query) {
			matches = append(matches, pkg)
		}
	}
	return matches
}

func runGolangCILint(ctx context.Context, packages ...string) error {
	args := append([]string{"run", "--build-tags", "mage"}, packages...)
	return runGolangCILintArgs(ctx, args...)
}

func runGolangCILintArgs(ctx context.Context, args ...string) error {
	root, err := repoRoot(ctx)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "golangci-lint", args...)
	cmd.Dir = root

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		text := strings.TrimSpace(output.String())
		if text != "" {
			fmt.Fprintln(os.Stderr, "Lint issues found:")
			fmt.Fprintln(os.Stderr, text)
		}
		return fmt.Errorf("golangci-lint failed")
	}

	text := strings.TrimSpace(output.String())
	if text != "" {
		fmt.Println(text)
	}
	fmt.Println("Lint passed with no issues.")
	return nil
}

func hasHead(ctx context.Context) bool {
	root, err := repoRoot(ctx)
	if err != nil {
		return false
	}

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "HEAD")
	cmd.Dir = root
	return cmd.Run() == nil
}

func hasParentCommit(ctx context.Context) bool {
	root, err := repoRoot(ctx)
	if err != nil {
		return false
	}

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "HEAD~1")
	cmd.Dir = root
	return cmd.Run() == nil
}

func repoRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return filepath.Abs(".")
	}

	root := strings.TrimSpace(string(output))
	if root == "" {
		return filepath.Abs(".")
	}
	return root, nil
}
