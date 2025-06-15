package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ikasamt/rd/pkg/config"
	"github.com/ikasamt/rd/pkg/redmine"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <issue-id>",
	Short: "Update a Redmine issue",
	Long:  `Update an existing Redmine issue with various options.`,
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

		update := &redmine.IssueUpdate{}
		hasUpdate := false

		// ステータス更新
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			if id, err := strconv.Atoi(status); err == nil {
				statusID := id
				update.StatusID = &statusID
				hasUpdate = true
			}
		}

		// 担当者更新
		if assignee, _ := cmd.Flags().GetString("assign"); assignee != "" {
			if assignee == "me" {
				// TODO: 現在のユーザーIDを取得
				assigneeID := 1 // 仮実装
				update.AssignedToID = &assigneeID
			} else if id, err := strconv.Atoi(assignee); err == nil {
				assigneeID := id
				update.AssignedToID = &assigneeID
			}
			hasUpdate = true
		}

		// 優先度更新
		if priority, _ := cmd.Flags().GetInt("priority"); priority > 0 {
			update.PriorityID = &priority
			hasUpdate = true
		}

		// 進捗率更新
		if doneRatio, _ := cmd.Flags().GetInt("done-ratio"); cmd.Flags().Changed("done-ratio") {
			update.DoneRatio = &doneRatio
			hasUpdate = true
		}

		// 開始日更新
		if startDate, _ := cmd.Flags().GetString("start-date"); startDate != "" {
			update.StartDate = &startDate
			hasUpdate = true
		}

		// 期限日更新
		if dueDate, _ := cmd.Flags().GetString("due-date"); dueDate != "" {
			update.DueDate = &dueDate
			hasUpdate = true
		}

		// コメント追加
		if note, _ := cmd.Flags().GetString("note"); note != "" {
			update.Notes = note
			hasUpdate = true
		}

		// カスタムフィールド更新
		fields, _ := cmd.Flags().GetStringSlice("field")
		if len(fields) > 0 {
			customFields := []redmine.CustomFieldValue{}
			for _, field := range fields {
				parts := strings.SplitN(field, "=", 2)
				if len(parts) == 2 {
					// TODO: カスタムフィールドIDの解決
					// 現在は仮実装
				}
			}
			if len(customFields) > 0 {
				update.CustomFields = customFields
				hasUpdate = true
			}
		}

		if !hasUpdate {
			return fmt.Errorf("no updates specified")
		}

		// 更新実行
		if err := client.UpdateIssue(issueID, update); err != nil {
			return fmt.Errorf("failed to update issue: %w", err)
		}

		fmt.Printf("Issue #%d updated successfully\n", issueID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().String("status", "", "Update status ID")
	updateCmd.Flags().String("assign", "", "Assign to user ID (or 'me')")
	updateCmd.Flags().Int("priority", 0, "Update priority ID")
	updateCmd.Flags().Int("done-ratio", 0, "Update done ratio (0-100)")
	updateCmd.Flags().String("start-date", "", "Update start date (YYYY-MM-DD)")
	updateCmd.Flags().String("due-date", "", "Update due date (YYYY-MM-DD)")
	updateCmd.Flags().String("note", "", "Add a note/comment")
	updateCmd.Flags().StringSlice("field", []string{}, "Update custom field (format: name=value)")
	updateCmd.Flags().Bool("interactive", false, "Interactive mode")
}