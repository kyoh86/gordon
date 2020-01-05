package history

import (
	"os"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/context"
)

type Store struct{}

func SaveHistory(ctx context.Context, repo *gogh.Repo, tag string) (retErr error) {
	if !ctx.HistorySave() {
		return nil
	}
	opts := []string{
		shellescape.Quote(repo.URL(ctx, false).String()),
	}
	if tag != "" {
		opts = append(opts, "--tag", shellescape.Quote(tag))
	}
	file, err := os.OpenFile(ctx.HistoryFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	if _, err := file.WriteString(strings.Join(opts, " ")); err != nil {
		return err
	}
	return nil
}
