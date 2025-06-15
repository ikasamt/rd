package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ikasamt/rd/pkg/config"
	"github.com/ikasamt/rd/pkg/redmine"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Redmine issue",
	Long:  `Create a new issue in Redmine with various options or interactive mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		urlFlag, _ := cmd.Root().Flags().GetString("url")
		keyFlag, _ := cmd.Root().Flags().GetString("key")
		
		cfg, err := config.Load(urlFlag, keyFlag)
		if err != nil {
			return err
		}

		client := redmine.NewClient(cfg.RedmineURL, cfg.APIKey)

		interactive, _ := cmd.Flags().GetBool("interactive")
		if interactive {
			return createIssueInteractive(client)
		}

		// コマンドラインオプションから作成
		return createIssueFromFlags(cmd, client)
	},
}

func createIssueFromFlags(cmd *cobra.Command, client *redmine.Client) error {
	title, _ := cmd.Flags().GetString("title")
	if title == "" {
		return fmt.Errorf("title is required (use --title or --interactive)")
	}

	projectID, _ := cmd.Flags().GetString("project")
	if projectID == "" {
		return fmt.Errorf("project is required (use --project or --interactive)")
	}

	// プロジェクトIDの取得
	project, err := client.GetProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	issue := &redmine.IssueCreate{
		ProjectID: project.ID,
		Subject:   title,
	}

	// オプションの設定
	description, _ := cmd.Flags().GetString("description")
	if description != "" {
		issue.Description = description
	}

	assignee, _ := cmd.Flags().GetString("assignee")
	if assignee != "" {
		if id, err := strconv.Atoi(assignee); err == nil {
			issue.AssignedToID = id
		}
	}

	trackerID, _ := cmd.Flags().GetInt("tracker")
	if trackerID > 0 {
		issue.TrackerID = trackerID
	}

	priorityID, _ := cmd.Flags().GetInt("priority")
	if priorityID > 0 {
		issue.PriorityID = priorityID
	}

	statusID, _ := cmd.Flags().GetInt("status")
	if statusID > 0 {
		issue.StatusID = statusID
	}

	startDate, _ := cmd.Flags().GetString("start-date")
	if startDate != "" {
		issue.StartDate = startDate
	}

	dueDate, _ := cmd.Flags().GetString("due-date")
	if dueDate != "" {
		issue.DueDate = dueDate
	}

	// カスタムフィールド
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
			issue.CustomFields = customFields
		}
	}

	// チケット作成
	created, err := client.CreateIssue(issue)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	// 出力
	jsonFlag, _ := cmd.Root().Flags().GetBool("json")
	if jsonFlag {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(created)
	}

	fmt.Printf("Issue #%d created successfully\n", created.ID)
	fmt.Printf("URL: %s/issues/%d\n", client.BaseURL, created.ID)
	return nil
}

func createIssueInteractive(client *redmine.Client) error {
	scanner := bufio.NewScanner(os.Stdin)

	// プロジェクト選択
	projects, err := client.ListProjects()
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	fmt.Println("Available projects:")
	for i, p := range projects.Projects {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}

	fmt.Print("\nSelect project number: ")
	scanner.Scan()
	projectNum, err := strconv.Atoi(scanner.Text())
	if err != nil || projectNum < 1 || projectNum > len(projects.Projects) {
		return fmt.Errorf("invalid project number")
	}
	selectedProject := projects.Projects[projectNum-1]

	// タイトル入力
	fmt.Print("\nIssue title: ")
	scanner.Scan()
	title := scanner.Text()
	if title == "" {
		return fmt.Errorf("title is required")
	}

	// 説明入力
	fmt.Print("\nDescription (optional, press Enter to skip): ")
	scanner.Scan()
	description := scanner.Text()

	// チケット作成
	issue := &redmine.IssueCreate{
		ProjectID:   selectedProject.ID,
		Subject:     title,
		Description: description,
	}

	created, err := client.CreateIssue(issue)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Printf("\nIssue #%d created successfully\n", created.ID)
	fmt.Printf("URL: %s/issues/%d\n", client.BaseURL, created.ID)
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("title", "", "Issue title")
	createCmd.Flags().String("description", "", "Issue description")
	createCmd.Flags().String("project", "", "Project ID or identifier")
	createCmd.Flags().String("assignee", "", "Assignee user ID")
	createCmd.Flags().Int("tracker", 0, "Tracker ID")
	createCmd.Flags().Int("priority", 0, "Priority ID")
	createCmd.Flags().Int("status", 0, "Status ID")
	createCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	createCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD)")
	createCmd.Flags().StringSlice("field", []string{}, "Custom field (format: name=value)")
	createCmd.Flags().Bool("interactive", false, "Interactive mode")
}