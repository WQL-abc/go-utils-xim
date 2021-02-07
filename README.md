# xim
xim is means The **extra** **indexer** for **map**. inspired by xian.  
Package `github.com/go-utils/xim` generates Indexes and Filters to search NoSQL Firestore.

It's designed especially for [Google Cloud Firestore](https://cloud.google.com/firestore?hl=en) and [firestore-repo](https://github.com/go-generalize/firestore-repo). (firestore-repo is firestore-model generator)

However, it doesn't depend on any specific Firestore APIs so that you may be able to use for some other NoSQL databases which support list-property and merge-join.  

[日本語ドキュメント](docs/ja.md)

# Features

* prefix/suffix/partial match search
* IN search
* reduce composite indexes(esp. for Cloud Firestore)

# Note

* Search latency can increase depending on its result-set size and filter condition.
* Index storage size can be bigger especially with long text prefix/suffix/partial match.

# Usage

Code example is for Cloud Firestore.

## Installation

```
$ go get -u github.com/go-utils/xim
```

## Configuration

```go
var bookIndexesConfig = xim.MustValidateConfig(&xim.Config{
	IgnoreCase:         true, // search case-insensitive
	SaveNoFiltersIndex: true, // always save 'NoFilters' index
})

// configure IN-filter
var statusInBuilder *xim.InBuilder = xim.NewInBuilder()

var (
    BookStatusUnpublished  = statusInBuilder.NewBit()
    BookStatusPublished    = statusInBuilder.NewBit()
    BookStatusDiscontinued = statusInBuilder.NewBit()
)
```

This configuration should be used to initialize both Indexes and Filters.

## Label Constants

Define common labels for both Indexes and Filters.  
Constants are not necessary but recommended.  
Short label names would make index size smaller.

```go
const (
	BookQueryLabelTitlePartial = "ti"
	BookQueryLabelTitlePrefix  = "tp"
	BookQueryLabelTitleSuffix  = "ts"
	BookQueryLabelIsHobby      = "h"
	BookQueryLabelStatusIN     = "s"
	BookQueryLabelPriceRange   = "pr"
)
```

## Save indexes

```go
idxs := xim.NewIndexes(bookIndexesConfig)

idxs.AddBigrams(BookQueryLabelTitlePartial, book.Title)
idxs.AddBiunigrams(BookQueryLabelTitlePartial, book.Title)
idxs.AddPrefixes(BookQueryLabelTitlePrefix, book.Title)
idxs.AddSuffixes(BookQueryLabelTitleSuffix, book.Title)
idxs.AddSomething(BookQueryLabelIsHobby, book.Category == "sports" || book.Category == "cooking")
idxs.Add(BookQueryLabelStatusIN, statusInBuilder.Indexes(BookStatusUnpublished)...)

switch {
case book.Price < 3000:
	idxs.Add(BookQueryLabelPriceRange, "p<3000")
case book.Price < 5000:
	idxs.Add(BookQueryLabelPriceRange, "3000<=p<5000")
case book.Price < 10000:
	idxs.Add(BookQueryLabelPriceRange, "5000<=p<10000")
default:
	idxs.Add(BookQueryLabelPriceRange, "10000<=p")
}

// build and set indexes to the book's property
var err error
book.Indexes, err = idxs.Build()
if err != nil {
	// error handling
}

// save book
```

## Search (example of Cloud Firestore)

```go
q := firestore.Query{}

filters := NewFilters(bookIndexesConfig).
    AddSomething(BookQueryLabelIsHolly, true).
    Add(BookQueryLabelStatusIN, statusInBuilder.Filter(BookStatusUnpublished, BookStatusPublished)).
    Add(BookQueryLabelPriceRange, "5000<=p<10000").
    AddBigrams(BookQueryLabelTitlePartial, title).
    AddBiunigrams(BookQueryLabelTitlePartial, title).
    AddSuffix(BookQueryLabelTitleSuffix, title)

built, err := filters.Build()
if err != nil {
    // error handling
}

for idx := range built {
    q = q.WherePath(firestore.FieldPath{"Indexes", idx}, "==", true)
}

// query books
```
