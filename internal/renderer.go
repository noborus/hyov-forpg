package internal

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/noborus/termhyo"
)

func printResults(ctx context.Context, w io.Writer, rows *sql.Rows, useOviewer bool) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	thCols := make([]termhyo.Column, len(cols))
	for i, c := range cols {
		thCols[i] = termhyo.Column{Title: c}
	}

	style := termhyo.VerticalBarStyle
	autoAlign := !useOviewer // oviewer uses false, others use true
	tw := termhyo.NewTable(w, thCols, termhyo.Border(style), termhyo.AutoAlign(autoAlign))

	for rows.Next() {
		// Check for cancellation during row processing
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}
		row := formatRow(columns, cols)
		tw.AddRow(row...)
	}

	tw.Render()
	return rows.Err()
}

func formatRow(columns []any, cols []string) []string {
	row := make([]string, len(cols))
	for i := range cols {
		v := columns[i]
		if v == nil {
			row[i] = "NULL"
			continue
		}
		switch val := v.(type) {
		case []byte:
			row[i] = string(val)
		default:
			row[i] = fmt.Sprintf("%v", val)
		}
	}
	return row
}
