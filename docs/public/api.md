---
title: Library API v1.0.0
language_tabs:
  - go: go
  - javascript: javascript
language_clients:
  - go: ""
  - javascript: ""
toc_footers: []
includes: []
search: false
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id="library-api">Library API v1.0.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="https://localhost:3000/api/v1">https://localhost:3000/api/v1</a>

<h1 id="library-api-default">Default</h1>

## List books in the library

<a id="opIdlistBooks"></a>

> Code samples

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "https://localhost:3000/api/v1/books", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```javascript

const headers = {
  'Accept':'application/json'
};

fetch('https://localhost:3000/api/v1/books',
{
  method: 'GET',

  headers: headers
})
.then(function(res) {
    return res.json();
}).then(function(body) {
    console.log(body);
});

```

`GET /books`

<h3 id="list-books-in-the-library-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|page_token|query|string|false|a pagination placeholder|
|total_size|query|integer(int32)|false|a pagination limit|

#### Detailed descriptions

**page_token**: a pagination placeholder

**total_size**: a pagination limit

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "title": "string",
      "isbn": 0
    }
  ],
  "next_page_token": "string"
}
```

<h3 id="list-books-in-the-library-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|success|[BookList](#schemabooklist)|
|default|Default|unexpected error|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_BookList">BookList</h2>
<!-- backwards compatibility -->
<a id="schemabooklist"></a>
<a id="schema_BookList"></a>
<a id="tocSbooklist"></a>
<a id="tocsbooklist"></a>

```json
{
  "items": [
    {
      "title": "string",
      "isbn": 0
    }
  ],
  "next_page_token": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[Book](#schemabook)]|true|none|none|
|next_page_token|string|true|none|none|

<h2 id="tocS_Book">Book</h2>
<!-- backwards compatibility -->
<a id="schemabook"></a>
<a id="schema_Book"></a>
<a id="tocSbook"></a>
<a id="tocsbook"></a>

```json
{
  "title": "string",
  "isbn": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|title|string|true|none|none|
|isbn|integer(int64)|true|none|none|

<h2 id="tocS_Error">Error</h2>
<!-- backwards compatibility -->
<a id="schemaerror"></a>
<a id="schema_Error"></a>
<a id="tocSerror"></a>
<a id="tocserror"></a>

```json
{
  "code": 0,
  "message": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|integer|true|none|none|
|message|string|true|none|none|

