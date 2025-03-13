# Fast Data Processing Using Golang

When we first received study data in the form of `json` files, there was a need to flatten the files, i.e. convert it to a table. Given the volume of data, Go was used to fast parsing. However, this approach was abandoned for the same reason why database systems were developed in the 1960s, namely that applications tightly couple code and data. If changes need to be made to the data, application logic needs to implement these changes. There is an element of reinventing the wheel constantly.

Instead, databases and database-like approaches were used.

> [!TIP]
> Always try to use a database where possible.
