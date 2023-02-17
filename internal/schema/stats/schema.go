/*
copyright 2020 the Escli authors

licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at

    http://www.apache.org/licenses/license-2.0

unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

package stats

type Stats struct {
	Shards Shards `json:"_shards"`
	All    All    `json:"_all"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

type All struct {
	Total Total `json:"total"`
}

type Total struct {
	Indexing Indexing `json:"indexing"`
	Search   Search   `json:"search"`
}

type Indexing struct {
	IndexTotal         int `json:"index_total"`
	IndexTimeInMillis  int `json:"index_time_in_millis"`
	DeleteTotal        int `json:"delete_total"`
	DeleteTimeInMillis int `json:"delete_time_in_millis"`
}

type Search struct {
	QueryTotal        int `json:"query_total"`
	QueryTimeInMillis int `json:"query_time_in_millis"`
	FetchTotal        int `json:"fetch_total"`
	FetchTimeInMillis int `json:"fetch_time_in_millis"`
}
