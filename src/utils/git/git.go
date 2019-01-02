package utils

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"gopkg.in/src-d/go-git.v4"
	gitSsh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"log"
	"strings"
)

// Git holds git repo data
type Git struct {
	Repo *git.Repository
	Path string
	url  string
}

// GitOps does stuff with git
func GitOps(c *configuration.Config) error {

	// Set up git instance
	repo := new(Git)
	repo.Path = c.Workspace

	// Open Repo for use by Contra
	err := GitOpen(repo)

	if err != nil {
		return err
	}

	worktree, err := repo.Repo.Worktree()

	if err != nil {
		return err
	}

	// Grab status and changes
	status, changes, err := GitStatus(*worktree)
	// Status will evaluate to true if something has changed
	if changes {
		// Commit if changes detected
		changesOut, changedFiles, err := Commit(repo.Path, status, *worktree)
		if err != nil {
			return err
		}
		log.Println("GIT changes committed.")

		err = gitSendEmail(c, changesOut, changedFiles)
		if err != nil {
			// Log the error, but carry on.
			log.Printf("WARNING: GIT notification email error: %v\n", err)
		}

		// push to remote if configured
		if c.GitPush {
			// If private key file is set, init public key auth.
			auth, err := gitSSHAuth(c)
			if err != nil {
				log.Printf(`WARNING: Error establish GIT authentication: "%s" changes will not be pushed.`, err)
				log.Printf(`WARNING: Verify GitAuth and GitPrivateKey configuration in %s`, c.ConfigFile)
				return nil
			}
			err = repo.Repo.Push(&git.PushOptions{Auth: auth})
			if err != nil {
				return err
			}
			log.Println("GIT Push successful.")
		}
	}

	return err
}

//gitSendEmail sends git related email notifications
func gitSendEmail(c *configuration.Config, changes, changedFiles []string) error {

	// Bail out if email is disabled
	if !c.EmailEnabled {
		log.Println("Email notifications are disabled.")
		return nil
	}

	// Convert slice of changes to a comma separated string
	changesString := strings.Join(changes, "\n")

	log.Printf("%s changed, sending email\n", strings.Join(changedFiles, ","))

	// Send email with changes
	err := utils.SendEmail(c, c.EmailSubject, changesString)

	if err != nil {
		return err
	}

	return nil
}

// gitSSHAuth sets up authentication for git a git remote
func gitSSHAuth(c *configuration.Config) (gitSsh.AuthMethod, error) {
	// If auth is disabled, there is nothing to do here.
	if !c.GitAuth {
		return nil, nil
	}

	auth, err := gitSsh.NewPublicKeysFromFile(c.GitUser, c.GitPrivateKey, "")
	return auth, err

}
