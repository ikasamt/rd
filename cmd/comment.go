package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ikasamt/rd/pkg/config"
	"github.com/ikasamt/rd/pkg/redmine"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment <issue-id> <comment>",
	Short: "Add a comment to a Redmine issue",
	Long:  `Add a comment (note) to an existing Redmine issue.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid issue ID: %s", args[0])
		}

		// 残りの引数をコメントとして結合
		comment := strings.Join(args[1:], " ")
		if comment == "" {
			return fmt.Errorf("comment cannot be empty")
		}

		urlFlag, _ := cmd.Root().Flags().GetString("url")
		keyFlag, _ := cmd.Root().Flags().GetString("key")
		debugFlag, _ := cmd.Root().Flags().GetBool("debug")
		
		cfg, err := config.Load(urlFlag, keyFlag)
		if err != nil {
			return err
		}

		client := redmine.NewClient(cfg.RedmineURL, cfg.APIKey)
		client.Debug = debugFlag

		// コメントだけの更新
		update := &redmine.IssueUpdate{
			Notes: comment,
		}

		if err := client.UpdateIssue(issueID, update); err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}

		fmt.Printf("Comment added to issue #%d\n", issueID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
}