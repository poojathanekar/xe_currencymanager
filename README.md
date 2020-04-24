# CurrencyConversion
XE currency data API

It provides up to date  information regarding currencies.

Overview:
1] Creates table to add currency value if table is not exist
2] Inserts entries for countriecodes and its currency value with regards to from exchange value
3] Updating one by one with respect to different exchange format
4] Still finding solution for postgresql in go to update for batch operations

Execution flow
1] sample .env file is provided please add your credetials to execute this program
2] if table is not exist it will create that table in your postgresql db
3] it will insert currency details with respect to each iterated currency value
4] If data with that exchange value already not exists it updates to that currency and mid value.


Benchmarking :

Total Time taken for 10 currencies  3.974s
