# HR XLXS to CSV Process

This is a simple utility that converts an automated Excel spreadsheet sent by UCF HR into a simplified set of CSVs for two email lists used with the PostMaster system. The first group includes all employees (so a simple conversion with no filtering) and the second group removes all OPS and Contingent Workers from the original list before creating the CSV.

## Usage
```sh
> go-hr-csv <input-file>
```
