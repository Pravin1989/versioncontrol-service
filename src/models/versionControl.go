package models

import (
	"context"

	"github.com/google/go-github/v42/github"
)

type GitVersionControl struct {
	SourceRepoName  string `json:"repoName"`
	SourceRepoOwner string `json:"repoOwner"`
	CommitBranch    string `json:"commitBranch"`
	BaseBranch      string `json:"baseBranch"`
	CommitMessage   string `json:"commitMessage"`
	AuthorUser
	PullRequest
}

type AuthorUser struct {
	AuthorName  string `json:"authorName"`
	AuthorEmail string `json:"authorEmail"`
}

type PullRequest struct {
	PrTitle       string `json:"prTitle"`
	PrDescription string `json:"prDescription"`
	PrLink        string `json:"prLink"`
}

type VersionControlOperations interface {
	CreateBranch(ctx context.Context, client *github.Client) (*github.Reference, error)
	CreatePR(ctx context.Context, client *github.Client) error
	CreateAndPushCommit(ctx context.Context, client *github.Client, ref *github.Reference, tree *github.Tree) error
	GetTree(ctx context.Context, client *github.Client, ref *github.Reference) (*github.Tree, error)
}
