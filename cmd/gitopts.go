package cmd

import (
	"fmt"
	"os"

	"github.com/khurshiduktamov/devlog/internal/git"
	"github.com/spf13/cobra"
)

// authorFlags holds the commit-author flags shared by today and standup.
type authorFlags struct {
	author     string
	allAuthors bool
}

func (f *authorFlags) register(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.author, "author", "",
		"only include commits whose author matches this pattern (default: your git user.email)")
	cmd.Flags().BoolVar(&f.allAuthors, "all-authors", false,
		"include commits from every author, not just you")
}

// apply fills in the author fields of opts. By default it limits commits
// to the current user's git email; if that isn't configured it warns and
// includes everyone rather than showing nothing.
func (f *authorFlags) apply(opts *git.Options) {
	switch {
	case f.allAuthors:
		return
	case f.author != "":
		opts.Author = f.author
	default:
		if git.CheckRepo(opts.Dir) != nil {
			return // no repo: GetCommits reports that, no author needed
		}
		email := git.CurrentUserEmail(opts.Dir)
		if email == "" {
			fmt.Fprintln(os.Stderr,
				"Warning: git user.email is not set, so commits from all authors are shown.\n"+
					"  Fix: git config --global user.email you@example.com")
			return
		}
		opts.Author = email
		opts.AuthorExact = true
	}
}
