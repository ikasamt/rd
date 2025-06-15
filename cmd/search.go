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

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search Redmine issues",
	Long:  `Search issues in Redmine using full-text search or filters.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		urlFlag, _ := cmd.Root().Flags().GetString("url")
		keyFlag, _ := cmd.Root().Flags().GetString("key")
		
		cfg, err := config.Load(urlFlag, keyFlag)
		if err != nil {
			return err
		}

		client := redmine.NewClient(cfg.RedmineURL, cfg.APIKey)

		// 検索用のフィルタを作成
		filter := &redmine.IssueFilter{}

		// プロジェクト指定
		project, _ := cmd.Flags().GetString("project")
		if project != "" {
			filter.ProjectID = project
		}

		// ステータス指定
		status, _ := cmd.Flags().GetString("status")
		if status != "" {
			filter.StatusID = status
		}

		// 検索の実装
		// Redmine APIには直接的な全文検索がないため、
		// subject や description にクエリを含むものを検索する
		// ここでは簡易的にリスト取得でクライアント側フィルタを行う
		issues, err := client.ListIssues(filter)
		if err != nil {
			return fmt.Errorf("failed to search issues: %w", err)
		}

		// クライアント側でのフィルタリング
		var filteredIssues []redmine.Issue
		queryLower := strings.ToLower(query)
		
		for _, issue := range issues.Issues {
			if strings.Contains(strings.ToLower(issue.Subject), queryLower) ||
			   strings.Contains(strings.ToLower(issue.Description), queryLower) {
				filteredIssues = append(filteredIssues, issue)
			}
		}

		// 出力形式の判定
		jsonFlag, _ := cmd.Root().Flags().GetBool("json")
		if jsonFlag {
			result := struct {
				Issues     []redmine.Issue `json:"issues"`
				TotalCount int             `json:"total_count"`
			}{
				Issues:     filteredIssues,
				TotalCount: len(filteredIssues),
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(result)
		}

		oneline, _ := cmd.Flags().GetBool("oneline")
		if oneline {
			for _, issue := range filteredIssues {
				fmt.Printf("#%d %s\n", issue.ID, issue.Subject)
			}
			return nil
		}

		// テーブル形式で出力
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tProject\tStatus\tPriority\tSubject\tAssignee")
		fmt.Fprintln(w, strings.Repeat("-", 80))
		
		for _, issue := range filteredIssues {
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
		
		if err := w.Flush(); err != nil {
			return err
		}

		fmt.Printf("\nFound %d issues matching '%s'\n", len(filteredIssues), query)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().String("project", "", "Filter by project ID")
	searchCmd.Flags().String("status", "", "Filter by status")
	searchCmd.Flags().Bool("oneline", false, "Display in one line format")
}