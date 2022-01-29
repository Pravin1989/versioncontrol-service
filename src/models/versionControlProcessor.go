package models

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"time"
	"versioncontrol-service/src/common"

	"github.com/google/go-github/v42/github"
)

func ConnectToRepoAndStartOperation(ctx context.Context, client *github.Client) {
	var versionControl = new(GitVersionControl)
	versionControl.SourceRepoOwner = "Pravin1989"
	versionControl.SourceRepoName = "sample-repo"
	versionControl.CommitBranch = "test"
	versionControl.BaseBranch = "main"
	versionControl.AuthorName = "Pravin"
	versionControl.AuthorEmail = "Pravin.Budge@gmail.com"
	versionControl.CommitMessage = "This is my new commit for test"
	versionControl.PRTitle = "PR for new changes"
	versionControl.PrDescription = "Please review and approve"
	ref, err := versionControl.CreateBranch(ctx, client)
	if err != nil {
		log.Printf("Failed to create branch %s\n", err)
		return
	}
	tree, err := versionControl.GetTree(ctx, client, ref)
	if err != nil {
		log.Printf("Failed to create the tree based on the provided files: %s\n", err)
		return
	}
	if pushErr := versionControl.CreateAndPushCommit(ctx, client, ref, tree); pushErr != nil {
		log.Printf("Failed to create the commit and push: %s\n", pushErr)
		return
	}
	if createPRErr := versionControl.CreatePR(ctx, client); createPRErr != nil {
		log.Printf("Failed to create the PR: %s\n", createPRErr)
		return
	}

}

// This method is used to create Pull Request
func (gv GitVersionControl) CreatePR(ctx context.Context, client *github.Client) error {
	newPR := &github.NewPullRequest{
		Title:               &gv.PRTitle,
		Head:                &gv.CommitBranch,
		Base:                &gv.BaseBranch,
		Body:                &gv.PrDescription,
		MaintainerCanModify: github.Bool(true),
	}
	prRepoOwner := gv.SourceRepoOwner
	pr, _, err := client.PullRequests.Create(ctx, prRepoOwner, gv.SourceRepoName, newPR)
	if err != nil {
		return err
	}
	gv.PRLink = pr.GetHTMLURL()
	log.Printf("PR created: %s\n", gv.PRLink)
	return nil
}

// This method is used to create commit and push
func (gv GitVersionControl) CreateAndPushCommit(ctx context.Context, client *github.Client, ref *github.Reference, tree *github.Tree) error {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, gv.SourceRepoOwner, gv.SourceRepoName, *ref.Object.SHA, nil)
	if err != nil {
		return nil
	}
	parent.Commit.SHA = parent.SHA
	//Create commit using tree
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &gv.AuthorName, Email: &gv.AuthorEmail}
	commit := &github.Commit{Author: author, Message: &gv.CommitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, commitErr := client.Git.CreateCommit(ctx, gv.SourceRepoOwner, gv.SourceRepoName, commit)
	if commitErr != nil {
		return commitErr
	}
	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, updateErr := client.Git.UpdateRef(ctx, gv.SourceRepoOwner, gv.SourceRepoName, ref, false)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

// This method is used to create tree for existing changes
func (gv GitVersionControl) GetTree(ctx context.Context, client *github.Client, ref *github.Reference) (*github.Tree, error) {
	tree, _, err := client.Git.GetTree(ctx, gv.SourceRepoOwner, gv.SourceRepoName, *ref.Object.SHA, true)
	if err != nil {
		return nil, err
	}
	entries := []*github.TreeEntry{}
	for i, value := range tree.Entries {
		blob, _, blobErr := client.Git.GetBlob(ctx, gv.SourceRepoOwner, gv.SourceRepoName, *value.SHA)
		if blobErr != nil {
			continue
		}
		byteArray, _ := base64.StdEncoding.DecodeString(blob.GetContent())
		changes := string(byteArray) + "\n Added one line for testing " + strconv.Itoa(i)
		entries = append(entries, &github.TreeEntry{Path: value.Path, Type: github.String("blob"), Content: github.String(changes), Mode: github.String("100644")})
	}
	tree, _, err = client.Git.CreateTree(ctx, gv.SourceRepoOwner, gv.SourceRepoName, *ref.Object.SHA, entries)
	return tree, nil
}

// This method is used to create branch
func (gv GitVersionControl) CreateBranch(ctx context.Context, client *github.Client) (*github.Reference, error) {
	if ref, _, err := client.Git.GetRef(ctx, gv.SourceRepoOwner, gv.SourceRepoName, common.Refs+gv.CommitBranch); err == nil {
		return ref, nil
	}
	//If error occured above The branch is not created, need to create it
	if gv.CommitBranch == gv.BaseBranch {
		return nil, errors.New("The base branch main is same as commit branch")
	}

	baseRef, _, baseErr := client.Git.GetRef(ctx, gv.SourceRepoOwner, gv.SourceRepoName, common.Refs+gv.BaseBranch)
	if baseErr != nil {
		return nil, baseErr
	}
	newRef := &github.Reference{Ref: github.String(common.Refs + gv.CommitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err := client.Git.CreateRef(ctx, gv.SourceRepoOwner, gv.SourceRepoName, newRef) //Create Branch
	if err != nil {
		return nil, errors.New("Failed to create branch")
	}
	return ref, nil
}
