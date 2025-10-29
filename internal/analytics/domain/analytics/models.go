package analytics

// Analiticas de nuevos clientes por mes
type NewClientRate struct {
	TotalCount int64 `json:"total_count"`
	Data       []struct {
		Month      string `json:"month"`
		NewClients int    `json:"new_clients"`
	} `json:"month_data"`
}

// Analiticas del total de ingresos por mes
type MonthlyRevenue struct {
	TotalRevenue float64 `json:"total_revenue"`
	Data         []struct {
		Month        string  `json:"month"`
		TotalRevenue float64 `json:"total_revenue"`
	} `json:"month_data"`
}
