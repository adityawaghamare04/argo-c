package webhook

import (
	"context"
	"fmt"
	"net/http"
	"github.com/google/go-github/v39/github"
)

// ReconcileWebhook checks if the repository webhook exists and reconciles state.
// It explicitly differentiates 404 Not Found errors from general API network failures.
func ReconcileWebhook(ctx context.Context, client *github.Client, owner, repo, targetUrl string) (bool, error) {
	hooks, resp, err := client.Repositories.ListHooks(ctx, owner, repo, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// Differentiate missing webhook (404) from API error
			return false, nil
		}
		// Return error for transient API failures to prevent erroneous deletion/re-creation
		return false, fmt.Errorf("github api error during webhook list: %w", err)
	}

	for _, h := range hooks {
		if h.Config != nil {
			if url, ok := h.Config["url"].(string); ok && url == targetUrl {
				return true, nil
			}
		}
	}
	return false, nil
}