package downloader

import "github.com/he8um/daryaft/internal/download"

type BatchEventHandler func(BatchItem, Event)

type BatchHandlers struct {
	ItemStarted func(BatchItem)
	Event       BatchEventHandler
}

func (d *Downloader) DownloadBatch(plan download.Plan, handlers BatchHandlers) BatchResult {
	items := make([]BatchItemResult, 0, len(plan.URLs))

	for index, rawURL := range plan.URLs {
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

		result, err := d.DownloadWithEvents(itemPlan, func(event Event) {
			if handlers.Event != nil {
				handlers.Event(item, event)
			}
		})

		items = append(items, BatchItemResult{
			Item:   item,
			Result: result,
			Err:    err,
		})
	}

	return BatchResult{Items: items}
}
