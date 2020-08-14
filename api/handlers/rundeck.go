package handlers

type rundeckResult struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Selected bool   `json:"selected"`
}

func fillRundeckResult(results []string, selected bool) []rundeckResult {
	var rundeckR []rundeckResult
	rundeckR = make([]rundeckResult, len(results), len(results))
	for i, r := range results {
		rundeckR[i] = rundeckResult{r, r, selected}
	}

	return rundeckR
}
