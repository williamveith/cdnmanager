# Cloudflare CDN Manager

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
