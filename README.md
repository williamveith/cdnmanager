# README

## About

This is the official Wails Vanilla template.

You can configure the project by editing [`wails.json`](https://wails.io/docs/reference/project-config)

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on [http://localhost:34115]. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## Cloudflare CDN Manager

[API Documentation](https://developers.cloudflare.com/api-next/go/)
[Examples](https://github.com/cloudflare/cloudflare-go/blob/master/workers_kv_example_test.go)

## **Functions: database.go*

- **`NewDatabase(dbName string) *Database`** Creates and returns a new `Database` instance connected to the specified SQLite database.

- **`CreateTable()`** Creates the `records` table if it does not already exist.

- **`DropTable()`** Drops the `records` table if it exists.

- **`InsertEntry(datavalues session.Entry)`** Inserts or replaces a single entry into the `records` table.

- **`InsertEntries(datavalues []session.Entry)`** Inserts or replaces multiple entries into the `records` table in a single transaction.

- **`GetEntryByName(name string) session.Entry`** Retrieves a single entry from the `records` table by its `name`.

- **`GetEntryByValue(value string) session.Entry`** Retrieves a single entry from the `records` table by its `value`.

- **`GetEntriesByValue(value string) []session.Entry`** Retrieves all entries from the `records` table with a matching `value`.

- **`GetAllEntries() []session.Entry`** Retrieves all entries from the `records` table.

- **`DeleteName(key string)`** Deletes a single entry from the `records` table by its `name`.

- **`DeleteNames(names []string)`** Deletes multiple entries from the `records` table where their `name` matches any of the specified names.

- **`DeleteEntry(entry session.Entry)`** Deletes a single entry from the `records` table by the entry's `name`.

- **`DeleteEntries(entries []session.Entry)`** Deletes multiple entries from the `records` table by the `name` of each provided entry.

- **`Size() int`** Returns the total number of entries in the `records` table.
