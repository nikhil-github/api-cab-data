package output

// Result represents output response.
type Result struct {
	Medallion string `json:"medallion"`
	Trips     int    `json:"trips"`
}
