package downloader

import (
	"context"

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

		items = append(items, BatchItemResult{
			Item:   item,
			Result: result,
			Err:    err,
		})

		if isCancellationError(err) {
			break
		}
	}

	return BatchResult{Planned: len(plan.URLs), Items: items}
}
