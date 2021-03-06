//
// Copyright © 2016-2020 Solus Project <copyright@getsol.us>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"fmt"
	"github.com/getsolus/solbuild/builder"
	"github.com/getsolus/solbuild/builder/source"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var deleteCacheCmd = &cobra.Command{
	Use:     "delete-cache",
	Short:   "delete solbuild cached files",
	Long:    `Delete assets stored on disk by solbuild`,
	Aliases: []string{"dc"},
	Run:     deleteCache,
}

// Whether we nuke *all* assets, i.e. sources too
var purgeAll bool

// Whether we should nuke images
var purgeImages bool

func init() {
	deleteCacheCmd.Flags().BoolVarP(&purgeAll, "all", "a", false, "Also delete ccache, packages and sources. Does not delete images.")
	deleteCacheCmd.Flags().BoolVarP(&purgeImages, "images", "i", false, "Deletes solbuild cached images.")
	RootCmd.AddCommand(deleteCacheCmd)
}

func deleteCache(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		profile = strings.TrimSpace(args[0])
	}

	if CLIDebug {
		log.SetLevel(log.DebugLevel)
	}
	log.StandardLogger().Formatter.(*log.TextFormatter).DisableColors = builder.DisableColors

	if os.Geteuid() != 0 {
		fmt.Fprintf(os.Stderr, "You must be root to delete caches\n")
		os.Exit(1)
	}

	manager, err := builder.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create new Manager: %v", err)
		os.Exit(1)
	}

	// By default include /var/lib/solbuild
	nukeDirs := []string{
		manager.Config.OverlayRootDir,
	}

	if purgeAll {
		nukeDirs = append(nukeDirs, []string{
			builder.CcacheDirectory,
			builder.LegacyCcacheDirectory,
			builder.PackageCacheDirectory,
			source.SourceDir,
		}...)
	}

	if purgeImages {
		nukeDirs = append(nukeDirs, []string{builder.ImagesDir}...)
	}

	for _, p := range nukeDirs {
		if !builder.PathExists(p) {
			continue
		}
		log.WithFields(log.Fields{
			"dir": p,
		}).Info("Removing cache directory")
		if err := os.RemoveAll(p); err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"dir":   p,
			}).Error("Could not remove cache directory")
			os.Exit(1)
		}
	}
}
