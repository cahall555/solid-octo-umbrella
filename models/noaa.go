package models

type NOAAActiveAlertsResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Features []struct {
		ID         string `json:"id"`
		Properties struct {
			Event       string `json:"event"`
			Severity    string `json:"severity"`
			Urgency     string `json:"urgency"`
			Certainty   string `json:"certainty"`
			Headline    string `json:"headline"`
			AreaDesc    string `json:"areaDesc"`
			Sent        string `json:"sent"`
			Effective   string `json:"effective"`
			Expires     string `json:"expires"`
			Instruction string `json:"instruction"`
		} `json:"properties"`
	} `json:"features"`
}
