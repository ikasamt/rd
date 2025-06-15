package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ikasamt/rd/pkg/config"
	"github.com/ikasamt/rd/pkg/redmine"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <issue-id>",
	Short: "Get a specific Redmine issue",
	Long:  `Display detailed information about a specific Redmine issue.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid issue ID: %s", args[0])
		}

		urlFlag, _ := cmd.Root().Flags().GetString("url")
		keyFlag, _ := cmd.Root().Flags().GetString("key")
		
		cfg, err := config.Load(urlFlag, keyFlag)
		if err != nil {
			return err
		}

		client := redmine.NewClient(cfg.RedmineURL, cfg.APIKey)

		includeComments, _ := cmd.Flags().GetBool("comments")
		issue, err := client.GetIssue(issueID, includeComments)
		if err != nil {
			return fmt.Errorf("failed to get issue: %w", err)
		}

		// 出力形式の判定
		jsonFlag, _ := cmd.Root().Flags().GetBool("json")
		if jsonFlag {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(issue)
		}

		return printIssueDetail(issue)
	},
}

func printIssueDetail(issue *redmine.Issue) error {
	fmt.Printf("Issue #%d\n", issue.ID)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Subject:     %s\n", issue.Subject)
	fmt.Printf("Project:     %s\n", issue.Project.Name)
	fmt.Printf("Tracker:     %s\n", issue.Tracker.Name)
	fmt.Printf("Status:      %s\n", issue.Status.Name)
	fmt.Printf("Priority:    %s\n", issue.Priority.Name)
	fmt.Printf("Author:      %s\n", issue.Author.Name)
	
	if issue.AssignedTo != nil {
		fmt.Printf("Assigned to: %s\n", issue.AssignedTo.Name)
	} else {
		fmt.Printf("Assigned to: -\n")
	}

	if issue.StartDate != nil {
		fmt.Printf("Start Date:  %s\n", *issue.StartDate)
	}
	if issue.DueDate != nil {
		fmt.Printf("Due Date:    %s\n", *issue.DueDate)
	}

	fmt.Printf("Done Ratio:  %d%%\n", issue.DoneRatio)
	if issue.EstimatedHours != nil {
		fmt.Printf("Estimated:   %.1f hours\n", *issue.EstimatedHours)
	}

	fmt.Printf("Created:     %s\n", issue.CreatedOn.Format(time.RFC3339))
	fmt.Printf("Updated:     %s\n", issue.UpdatedOn.Format(time.RFC3339))

	// カスタムフィールド
	if len(issue.CustomFields) > 0 {
		fmt.Println("\nCustom Fields:")
		for _, cf := range issue.CustomFields {
			fmt.Printf("  %s: %v\n", cf.Name, cf.Value)
		}
	}

	// 説明
	if issue.Description != "" {
		fmt.Println("\nDescription:")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Println(issue.Description)
	}

	// コメント（ジャーナル）
	if len(issue.Journals) > 0 {
		fmt.Println("\nComments:")
		fmt.Println(strings.Repeat("-", 80))
		for _, journal := range issue.Journals {
			if journal.Notes != "" {
				fmt.Printf("\n[%s] %s:\n%s\n",
					journal.CreatedOn.Format("2006-01-02 15:04"),
					journal.User.Name,
					journal.Notes)
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().Bool("comments", false, "Include comments")
	getCmd.Flags().Bool("fields", false, "Include custom fields")
}