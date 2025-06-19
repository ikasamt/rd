package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ikasamt/rd/pkg/config"
	"github.com/ikasamt/rd/pkg/redmine"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search Redmine resources",
	Long:  `Search issues and other resources in Redmine using the search API.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		urlFlag, _ := cmd.Root().Flags().GetString("url")
		keyFlag, _ := cmd.Root().Flags().GetString("key")
		debugFlag, _ := cmd.Root().Flags().GetBool("debug")
		
		cfg, err := config.Load(urlFlag, keyFlag)
		if err != nil {
			return err
		}

		client := redmine.NewClient(cfg.RedmineURL, cfg.APIKey)
		client.Debug = debugFlag

		// 検索オプションの設定
		opts := &redmine.SearchOptions{
			Query:      query,
			Limit:      100, // デフォルトのページサイズ
			Issues:     true, // デフォルトでIssuesを検索
		}

		// フラグから検索対象を設定
		if all, _ := cmd.Flags().GetBool("all-types"); all {
			opts.Issues = true
			opts.News = true
			opts.Documents = true
			opts.Changesets = true
			opts.WikiPages = true
			opts.Messages = true
			opts.Projects = true
		}

		// 個別の検索対象フラグ
		if v, _ := cmd.Flags().GetBool("issues"); v {
			opts.Issues = v
		}
		if v, _ := cmd.Flags().GetBool("wiki"); v {
			opts.WikiPages = v
		}
		if v, _ := cmd.Flags().GetBool("news"); v {
			opts.News = v
		}
		if v, _ := cmd.Flags().GetBool("documents"); v {
			opts.Documents = v
		}
		if v, _ := cmd.Flags().GetBool("changesets"); v {
			opts.Changesets = v
		}
		if v, _ := cmd.Flags().GetBool("messages"); v {
			opts.Messages = v
		}
		if v, _ := cmd.Flags().GetBool("projects"); v {
			opts.Projects = v
		}

		// 検索範囲
		if scope, _ := cmd.Flags().GetString("scope"); scope != "" {
			opts.Scope = scope
		}

		// 検索オプション
		if v, _ := cmd.Flags().GetBool("titles-only"); v {
			opts.TitlesOnly = v
		}
		if v, _ := cmd.Flags().GetBool("all-words"); v {
			opts.AllWords = v
		}

		// ページサイズ
		if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
			opts.Limit = limit
		}

		// 全件検索の処理
		allFlag, _ := cmd.Flags().GetBool("all")
		var allResults []redmine.SearchResult

		if allFlag {
			// ページネーションで全件取得
			offset := 0
			for {
				opts.Offset = offset
				result, err := client.Search(opts)
				if err != nil {
					return fmt.Errorf("search failed: %w", err)
				}

				if len(result.Results) == 0 {
					break
				}

				allResults = append(allResults, result.Results...)

				// 次のページがあるかチェック
				if offset+len(result.Results) >= result.TotalCount {
					break
				}

				offset += opts.Limit
			}
		} else {
			// 1ページのみ取得
			result, err := client.Search(opts)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}
			allResults = result.Results
		}

		// 出力形式の判定
		jsonFlag, _ := cmd.Root().Flags().GetBool("json")
		if jsonFlag {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(map[string]interface{}{
				"results": allResults,
				"total":   len(allResults),
			})
		}

		oneline, _ := cmd.Flags().GetBool("oneline")
		if oneline {
			for _, result := range allResults {
				// IDを抽出（URLから）
				parts := strings.Split(result.URL, "/")
				id := parts[len(parts)-1]
				fmt.Printf("%s: %s\n", id, result.Title)
			}
			return nil
		}

		// テーブル形式で出力
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Type\tID\tTitle\tDescription")
		fmt.Fprintln(w, strings.Repeat("-", 80))
		
		for _, result := range allResults {
			// URLからIDを抽出
			id := extractIDFromURL(result.URL)
			
			// タイトルとDescriptionを短縮
			title := result.Title
			if len(title) > 50 {
				title = title[:47] + "..."
			}
			
			desc := result.Description
			if len(desc) > 30 {
				desc = desc[:27] + "..."
			}
			// HTMLタグを除去
			desc = strings.ReplaceAll(desc, "<strong class=\"highlight\">", "")
			desc = strings.ReplaceAll(desc, "</strong>", "")
			
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				result.Type,
				id,
				title,
				desc,
			)
		}
		
		if err := w.Flush(); err != nil {
			return err
		}

		fmt.Printf("\nFound %d results matching '%s'\n", len(allResults), query)
		return nil
	},
}

// URLからIDを抽出するヘルパー関数
func extractIDFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		// 最後の部分がIDの場合
		if id := parts[len(parts)-1]; id != "" {
			return id
		}
		// issuesの場合、issues/123 の形式
		for i, part := range parts {
			if part == "issues" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	// IDを抽出できない場合
	if idInt, err := strconv.Atoi(url); err == nil {
		return strconv.Itoa(idInt)
	}
	return "-"
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// 検索対象の指定
	searchCmd.Flags().Bool("all-types", false, "Search all resource types")
	searchCmd.Flags().Bool("issues", false, "Search issues (default: true)")
	searchCmd.Flags().Bool("wiki", false, "Search wiki pages")
	searchCmd.Flags().Bool("news", false, "Search news")
	searchCmd.Flags().Bool("documents", false, "Search documents")
	searchCmd.Flags().Bool("changesets", false, "Search changesets")
	searchCmd.Flags().Bool("messages", false, "Search messages")
	searchCmd.Flags().Bool("projects", false, "Search projects")
	
	// 検索範囲
	searchCmd.Flags().String("scope", "", "Search scope: all, my_projects, subprojects")
	
	// 検索オプション
	searchCmd.Flags().Bool("titles-only", false, "Search in titles only")
	searchCmd.Flags().Bool("all-words", false, "Match all query words")
	
	// ページネーション
	searchCmd.Flags().Bool("all", false, "Fetch all results (may take longer)")
	searchCmd.Flags().Int("limit", 100, "Number of results per page")
	
	// 出力形式
	searchCmd.Flags().Bool("oneline", false, "Display in one line format")
}