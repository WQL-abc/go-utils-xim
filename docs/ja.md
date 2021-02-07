# xim
ximはThe **extra** **indexer** for **map**の略称。 xianに影響を受けました。  
`github.com/go-utils/xim` はNoSQL Firestoreを検索するためのインデックスとフィルターを生成します。

[Google Cloud Firestore](https://cloud.google.com/firestore?hl=ja) と [firestore-repo](https://github.com/go-generalize/firestore-repo) 用に設計されています。  
(firestore-repoはモデルとなる構造体からCRUD/クエリ検索部分のコードを生成します)

ただし、Firestore APIに依存しないため、list-propertyとmerge-joinをサポートする他のNoSQLデータベースで使用できる場合があります。

# 特徴

* 前方/後方/部分 一致 検索
* IN 検索
* 複合インデックスを減らす(特にCloud Firestore)

# 備考

* 検索の待ち時間は、結果セットのサイズとフィルターの状態によっては長くなる可能性があります。
* インデックスのストレージサイズは、特に長いテキストは前方/後方/部分一致の場合に大きくなる可能性があります。

# 使い方

コード例はCloud Firestore用です。

## Installation

```
$ go get -u github.com/go-utils/xim
```

## 構成

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

この構成は、インデックスとフィルターの両方を初期化するために使用する必要があります。

## ラベル定数

インデックスとフィルターの両方に共通のラベルを定義します。  
定数は必須ではありませんが、お勧めします。  
ラベル名が短いと、インデックスサイズが小さくなります。

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

## インデックスを保存する

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

## 検索 (Cloud Firestoreの例)

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
