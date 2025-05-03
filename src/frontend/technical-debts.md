# Technical Debts - Frontend

This file is intended to list all technical debts identified during development.
These may have arisen due to time constraints or limited knowledge at the time of implementation.

## List

| File / Topic                        | Debt Description                                                                                                                                                                                   |
|-------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `tabRowHeigt`                       | The `tabRowHeigt` is a property of the struct AppModel in `main.go`. However, this value is used in other files (`ospfMonitoring/view.go`) and not handed over. Just redeclared in the other file. |
| `main.go`, `ospfMonitoring/view.go` | Variables that are initialized outside of any function.                                                                                                                                            |
| `shell/*`                           | Currently only simple bash commands work, but not commands with \| or >.                                                                                                                           |
| `ospfMonitoring/view.go`            | When scrolling through the tables (e.g. Router LSAs) the table entries and the area are switched randomly. --> Probably triggered with all kinds of page reloads.                                  |
