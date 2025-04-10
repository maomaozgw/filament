package model

type Statistic struct {
	Kind   string           `json:"kind"`
	Title  string           `json:"title"`
	Values []StatisticValue `json:"values"`
}

type StatisticValue struct {
	Name     string           `json:"name"`
	Value    int64            `json:"value"`
	Chileren []StatisticValue `gorm:"-" json:"children,omitempty"`
}
