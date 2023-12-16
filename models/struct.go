package models

type ProofStruct struct {
	Ar struct {
		X string `json:"X"`
		Y string `json:"Y"`
	} `json:"Ar"`
	Krs struct {
		X string `json:"X"`
		Y string `json:"Y"`
	} `json:"Krs"`
	Bs struct {
		X struct {
			A0 string `json:"A0"`
			A1 string `json:"A1"`
		} `json:"X"`
		Y struct {
			A0 string `json:"A0"`
			A1 string `json:"A1"`
		} `json:"Y"`
	} `json:"Bs"`
	Commitments   []interface{} `json:"Commitments"` // Assuming array of unspecified type
	CommitmentPok struct {
		X int `json:"X"`
		Y int `json:"Y"`
	} `json:"CommitmentPok"`
}

type APIResponseCelestia struct {
	Result int `json:"result"`
}
