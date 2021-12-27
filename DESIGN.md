# Rationale

I use Firefly III to record transactions, but an Excel sheet for budgetting. I want to have yearly budgets and associated reporting (like in the Excel sheet) in the same app that handles recording data.

# Desired features

- Yearly budgets per category, with totals reported per month.
- Can create, edit, and view past, present, and future budgets
- Record transactions.
- Support email reports
- Nice mobile UI

# Implementation ideas

- Write a whole new app from scratch
    - Lots of work, possibly lots of unknowns around accounting
- Write a budget app that uses the Firefly API for transactions, but stores its own budget data
    - No need to handle data import / export
    - No need to implement transaction or account data models
    - Could decide to add my own transaction backend later
    - Can still create my own UI
    - Could even choose to query Firefly's DB directly for transactions / categories / accounts

# Data model

Table: Budgets
    - ID
    - Start datetime
    - End datetime
    - Reporting interval (e.g. monthly)
Table: CategoryBudgets
    - ID
    - BudgetID (FK)
    - CategoryID (or key, whatever Firefly uses)
    - Amount (budgeted amount over the budget end-start dates)


# Implementation plan

What I'll do first:

- Make a budget in the backend, talk to the database
- Make a front-end that can show a budget status homepage, including current totals
- Add to UI:
    - Create & edit budget
    - Create transactions

# FAQ

Q: How do we manage multiple budgets?

A: When showing the status page, the budget containing the current date is selected. We prevent a new budget from being created if it overlaps an existing budget. Previous budgets can be browsed, and future budgets can be created.