# gournal
Gournal is a small error tracking system intended to simplify error reporting and logging.

Gournal works with execution cycles where you can easily store errors encountered during the program execution.
At every moment you know if an error has occured during the current cycle, or execute your own custom reports.

## Usage
```

			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)