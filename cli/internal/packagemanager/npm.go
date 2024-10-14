package packagemanager

import (
	"fmt"

	"github.com/khulnasoft/titanrepo/cli/internal/fs"
	"github.com/khulnasoft/titanrepo/cli/internal/lockfile"
	"github.com/khulnasoft/titanrepo/cli/internal/titanpath"
)

var nodejsNpm = PackageManager{
	Name:         "nodejs-npm",
	Slug:         "npm",
	Command:      "npm",
	Specfile:     "package.json",
	Lockfile:     "package-lock.json",
	PackageDir:   "node_modules",
	ArgSeparator: []string{"--"},

	getWorkspaceGlobs: func(rootpath titanpath.AbsoluteSystemPath) ([]string, error) {
		pkg, err := fs.ReadPackageJSON(rootpath.UntypedJoin("package.json"))
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. Titanrepo requires npm workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},

	getWorkspaceIgnores: func(pm PackageManager, rootpath titanpath.AbsoluteSystemPath) ([]string, error) {
		// Matches upstream values:
		// function: https://github.com/npm/map-workspaces/blob/a46503543982cb35f51cc2d6253d4dcc6bca9b32/lib/index.js#L73
		// key code: https://github.com/npm/map-workspaces/blob/a46503543982cb35f51cc2d6253d4dcc6bca9b32/lib/index.js#L90-L96
		// call site: https://github.com/npm/cli/blob/7a858277171813b37d46a032e49db44c8624f78f/lib/workspaces/get-workspaces.js#L14
		return []string{
			"**/node_modules/**",
		}, nil
	},

	Matches: func(manager string, version string) (bool, error) {
		return manager == "npm", nil
	},

	detect: func(projectDirectory titanpath.AbsoluteSystemPath, packageManager *PackageManager) (bool, error) {
		specfileExists := projectDirectory.UntypedJoin(packageManager.Specfile).FileExists()
		lockfileExists := projectDirectory.UntypedJoin(packageManager.Lockfile).FileExists()

		return (specfileExists && lockfileExists), nil
	},

	canPrune: func(cwd titanpath.AbsoluteSystemPath) (bool, error) {
		return true, nil
	},

	readLockfile: func(contents []byte) (lockfile.Lockfile, error) {
		return lockfile.DecodeNpmLockfile(contents)
	},
}
