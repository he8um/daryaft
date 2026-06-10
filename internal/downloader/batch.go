package downloader

import (
	"context"
	"fmt"

	"github.com/he8um/daryaft/internal/checksum"
	"github.com/he8um/daryaft/internal/download"
)

type BatchEventHandler func(BatchItem, Event)

type BatchHandlers struct {
	ItemStarted func(BatchItem)
	Event       BatchEventHandler
}

func (d *Downloader) DownloadBatch(plan download.Plan, handlers BatchHandlers) BatchResult {
	return d.DownloadBatchContext(context.Background(), plan, handlers)
}

func (d *Downloader) DownloadBatchContext(ctx context.Context, plan download.Plan, handlers BatchHandlers) BatchResult {
	if ctx == nil {
		ctx = context.Background()
	}
	items := make([]BatchItemResult, 0, len(plan.URLs))
	checksumVerified := 0

	for index, rawURL := range plan.URLs {
		if ctx.Err() != nil {
			break
		}

		item := BatchItem{
			Index: index + 1,
			Total: len(plan.URLs),
			URL:   rawURL,
		}

		if handlers.ItemStarted != nil {
			handlers.ItemStarted(item)
		}

		itemPlan := plan
		itemPlan.URLs = []string{rawURL}

		result, err := d.DownloadWithEventsContext(ctx, itemPlan, func(event Event) {
			if handlers.Event != nil {
				handlers.Event(item, event)
			}
		})

		checksumStatus := ""
		if err == nil {
			if spec, ok := checksumSpecForURL(plan, rawURL); ok {
				if verifyErr := verifyItemChecksum(result, spec); verifyErr != nil {
					err = verifyErr
					checksumStatus = "failed"
				} else {
					checksumStatus = "verified"
					checksumVerified++
				}
			}
		}

		items = append(items, BatchItemResult{
			Item:           item,
			Result:         result,
			Err:            err,
			ChecksumStatus: checksumStatus,
		})

		if isCancellationError(err) {
			break
		}
	}

	return BatchResult{Planned: len(plan.URLs), ChecksumVerified: checksumVerified, Items: items}
}

// checksumSpecForURL returns the checksum spec that applies to a download
// target, if any. Per-target checksum files take precedence; otherwise a
// single-target plan checksum applies.
func checksumSpecForURL(plan download.Plan, rawURL string) (checksum.Spec, bool) {
	if plan.TargetChecksums != nil {
		spec, ok := plan.TargetChecksums[rawURL]
		return spec, ok
	}
	if plan.Checksum != nil {
		return *plan.Checksum, true
	}
	return checksum.Spec{}, false
}

// verifyItemChecksum verifies a completed download against a checksum spec. It
// never removes the downloaded file on mismatch; the file is left in place for
// inspection. The returned error already includes expected and actual digests.
func verifyItemChecksum(result Result, spec checksum.Spec) error {
	if result.Path == "" {
		return fmt.Errorf("checksum verification skipped: no completed file path available")
	}
	if _, err := checksum.VerifyFile(result.Path, spec); err != nil {
		return err
	}
	return nil
}
