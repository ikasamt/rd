package config

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type Config struct {
    RedmineURL string
    APIKey     string
}

// Load resolves configuration in the following priority:
// 1) Flags (--url, --key)
// 2) Environment variables (REDMINE_URL, REDMINE_API_KEY)
// 3) .rd in current directory
// 4) .rd in Git repository root
// 5) ~/.rd
func Load(urlFlag, keyFlag string) (*Config, error) {
    cfg := &Config{}

    // 1) Flags
    if urlFlag != "" {
        cfg.RedmineURL = urlFlag
    }
    if keyFlag != "" {
        cfg.APIKey = keyFlag
    }

    // 2) Environment variables
    applyIfEmpty(cfg, &Config{
        RedmineURL: os.Getenv("REDMINE_URL"),
        APIKey:     os.Getenv("REDMINE_API_KEY"),
    })

    // 3-5) .rd files
    for _, p := range candidateConfigPaths() {
        if p == "" {
            continue
        }
        if _, err := os.Stat(p); err == nil {
            if fileCfg, err := loadFromRD(p); err == nil {
                applyIfEmpty(cfg, fileCfg)
            }
        }
    }

    // Validation
    if cfg.RedmineURL == "" {
        return nil, fmt.Errorf("Redmine URL is not set. Set via --url, REDMINE_URL, or .rd")
    }
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("Redmine API key is not set. Set via --key, REDMINE_API_KEY, or .rd")
    }
    return cfg, nil
}

// applyIfEmpty copies non-empty fields from src to dst, but only when dst fields are empty.
func applyIfEmpty(dst, src *Config) {
    if dst.RedmineURL == "" && src.RedmineURL != "" {
        dst.RedmineURL = src.RedmineURL
    }
    if dst.APIKey == "" && src.APIKey != "" {
        dst.APIKey = src.APIKey
    }
}

// candidateConfigPaths returns .rd candidate paths in priority order (highest first)
// CWD -> Git root -> Home dir.
func candidateConfigPaths() []string {
    var paths []string

    // Current working directory
    if wd, err := os.Getwd(); err == nil {
        paths = append(paths, filepath.Join(wd, ".rd"))
    }

    // Git root (by walking up for .git)
    if gitRoot, err := findGitRoot(); err == nil && gitRoot != "" {
        // Avoid duplicating CWD if it's already the same path
        if wd, err := os.Getwd(); err == nil {
            if filepath.Clean(gitRoot) != filepath.Clean(wd) {
                paths = append(paths, filepath.Join(gitRoot, ".rd"))
            }
        } else {
            paths = append(paths, filepath.Join(gitRoot, ".rd"))
        }
    }

    // Home directory
    if home, err := os.UserHomeDir(); err == nil {
        paths = append(paths, filepath.Join(home, ".rd"))
    }

    return paths
}

// findGitRoot walks up from the current directory to find a directory containing .git.
func findGitRoot() (string, error) {
    start, err := os.Getwd()
    if err != nil {
        return "", err
    }
    dir := start
    for {
        if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
            return dir, nil
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            break
        }
        dir = parent
    }
    return "", errors.New("git root not found")
}

// loadFromRD parses a simple key-value file format:
// - Lines: KEY=VALUE or KEY: VALUE
// - Comments start with # or ;
// - Supported keys (case-insensitive):
//     REDMINE_URL, URL
//     REDMINE_API_KEY, API_KEY, KEY
func loadFromRD(path string) (*Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    cfg := &Config{}
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
            continue
        }
        // Split by = or : (first occurrence)
        var key, val string
        if i := strings.IndexAny(line, "=:"); i >= 0 {
            key = strings.TrimSpace(line[:i])
            val = strings.TrimSpace(line[i+1:])
        } else {
            // Unknown format, skip
            continue
        }
        lk := strings.ToUpper(strings.TrimSpace(key))
        switch lk {
        case "REDMINE_URL", "URL":
            if cfg.RedmineURL == "" {
                cfg.RedmineURL = val
            }
        case "REDMINE_API_KEY", "API_KEY", "KEY":
            if cfg.APIKey == "" {
                cfg.APIKey = val
            }
        }
    }
    // Ignore scanner.Err() to keep robust; caller treats empty cfg as no data
    return cfg, nil
}
