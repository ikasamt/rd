package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ikasamt/rd/pkg/config"
	"github.com/ikasamt/rd/pkg/redmine"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Redmine issues",
	Long:  `List issues from Redmine with various filters and options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		urlFlag, _ := cmd.Root().Flags().GetString("url")
		keyFlag, _ := cmd.Root().Flags().GetString("key")
		debugFlag, _ := cmd.Root().Flags().GetBool("debug")
		
		cfg, err := config.Load(urlFlag, keyFlag)
		if err != nil {
			return err
		}

		client := redmine.NewClient(cfg.RedmineURL, cfg.APIKey)
		client.Debug = debugFlag

		// フィルタの設定
		filter := &redmine.IssueFilter{}
		
		project, _ := cmd.Flags().GetString("project")
		if project != "" {
			filter.ProjectID = project
		}

		status, _ := cmd.Flags().GetString("status")
		if status != "" {
			filter.StatusID = status
		}

		assignee, _ := cmd.Flags().GetString("assignee")
		if assignee != "" {
			if assignee == "me" {
				assignee = "me"
			}
			filter.AssignedTo = assignee
		}

		// 取得
		issues, err := client.ListIssues(filter)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		// 出力形式の判定
		jsonFlag, _ := cmd.Root().Flags().GetBool("json")
		if jsonFlag {
			return outputJSON(issues)
		}

		oneline, _ := cmd.Flags().GetBool("oneline")
		if oneline {
			return outputOneline(issues.Issues)
		}

		csv, _ := cmd.Flags().GetBool("csv")
		if csv {
			return outputCSV(issues.Issues)
		}

		return outputTable(issues.Issues)
	},
}

func outputJSON(issues *redmine.IssuesResponse) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(issues)
}

func outputOneline(issues []redmine.Issue) error {
	for _, issue := range issues {
		fmt.Printf("#%d %s\n", issue.ID, issue.Subject)
	}
	return nil
}

func outputCSV(issues []redmine.Issue) error {
	fmt.Println("ID,Project,Status,Priority,Subject,Assignee")
	for _, issue := range issues {
		assignee := ""
		if issue.AssignedTo != nil {
			assignee = issue.AssignedTo.Name
		}
		fmt.Printf("%d,%s,%s,%s,%q,%s\n",
			issue.ID,
			issue.Project.Name,
			issue.Status.Name,
			issue.Priority.Name,
			issue.Subject,
			assignee,
		)
	}
	return nil
}

func outputTable(issues []redmine.Issue) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tProject\tStatus\tPriority\tSubject\tAssignee")
	fmt.Fprintln(w, strings.Repeat("-", 80))
	
	for _, issue := range issues {
		assignee := "-"
		if issue.AssignedTo != nil {
			assignee = issue.AssignedTo.Name
		}
		
		subject := issue.Subject
		if len(subject) > 40 {
			subject = subject[:37] + "..."
		}
		
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			issue.ID,
			issue.Project.Name,
			issue.Status.Name,
			issue.Priority.Name,
			subject,
			assignee,
		)
	}
	
	return w.Flush()
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().Bool("all", false, "Show all issues")
	listCmd.Flags().String("project", "", "Filter by project ID")
	listCmd.Flags().String("status", "", "Filter by status")
	listCmd.Flags().String("assignee", "", "Filter by assignee")
	listCmd.Flags().StringSlice("field", []string{}, "Filter by custom field (format: name=value)")
	listCmd.Flags().Bool("oneline", false, "Display in one line format")
	listCmd.Flags().Bool("csv", false, "Output in CSV format")
}