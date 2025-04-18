# Technical Debts - Frontend

This file is intended to list all technical debts identified during development.
These may have arisen due to time constraints or limited knowledge at the time of implementation.

## List

| File / Topic                        | Debt Description                                                                                                                                                                                   |
|-------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `tabRowHeigt`                       | The `tabRowHeigt` is a property of the struct AppModel in `main.go`. However, this value is used in other files (`ospfMonitoring/view.go`) and not handed over. Just redeclared in the other file. |
| `main.go`, `ospfMonitoring/view.go` | Variables that are initialized outside of any function.                                                                                                                                            |
| all `view.go`                       | In all view files the box size is calculated separately. It would be cleaner to have const for boxSizeForOneBox, boxSizeForTwoBox, etc. but this is hard to accomplish with dynamic resizing.     |
| `shell/*`                           | Currently the bash custom command does not work on the routers.                                                                                                                                    |
