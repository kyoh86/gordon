package history

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alessio/shellescape"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/context"
)

type Store struct{}

var mkdirAllOnce sync.Once

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
	dir := filepath.Dir(ctx.HistoryFile())
	mkdirAllOnce.Do(func() {
		log.Printf("debug: create directory %s", dir)
		retErr = os.MkdirAll(dir, 0700)
	})
	if retErr != nil {
		return
	}
	file, err := os.OpenFile(ctx.HistoryFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		retErr = fmt.Errorf("open history file: %w", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = fmt.Errorf("close history file: %w", err)
		}
	}()
	if _, err := file.WriteString(strings.Join(opts, " ") + "\n"); err != nil {
		return err
	}
	return nil
}
