# Technical Debts - Frontend

This file is intended to list all technical debts identified during development.
These may have arisen due to time constraints or limited knowledge at the time of implementation.

## List

| File / Topic  | Debt Description                                                                                                                                                                                   |
|---------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `tabRowHeigt` | The `tabRowHeigt` is a property of the struct AppModel in `main.go`. However, this value is used in other files (`ospfMonitoring/view.go`) and not handed over. Just redeclared in the other file. |
