package action

import "github.com/malnick/cryptorious/vault"

var (
	expectedUsername   = "test-username"
	expectedPassword   = "test-password"
	expectedSecureNote = "test-secure-note"

	testVaultSet = vault.Set{
		Username: vault.EncryptedEntry{
			Ciphertext: "AAAAAAAAAAAAAAAAAAAAALthIXEjHFp2BzZFaM2NzZI=",
			Key:        "AQIBAHjjCeNKbXvpFm2HyIXoKqPvjyyCDoSbmOkO8uSNKqBm0wE/eSqj6B0NnrczvwcLSlzrAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMrzjLCH0Te9Egi1UsAgEQgDtiaSFOxBU3jpEU1+OBCg0fOMTCin1IO7uWMtPB5Rq+MArZiFDNNIl4V7P7WTGZk5p9nEVXZibi0evnKA==",
			IV:         "4VvVxXdouSf+ANN+PJ5+tg==",
		},

		Password: vault.EncryptedEntry{
			Ciphertext: "AAAAAAAAAAAAAAAAAAAAAGbAqmydcY5wALQOAz0/m28=",
			IV:         "fbAcI8NkQbFYWHCtOkCm/Q==",
			Key:        "AQIBAHjjCeNKbXvpFm2HyIXoKqPvjyyCDoSbmOkO8uSNKqBm0wFX/eO80576OpuMIJRk6ZvlAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMfYDlWalA7BKoXxJEAgEQgDvLKQigz+X1DDNGf0KAhWLcBXKh1DItzyZsPodwArUoU0wAj7tb3MWhQXJN4XcDFLfh4sln0ozr4vxtYA==",
		},

		SecureNote: vault.EncryptedEntry{
			Ciphertext: "AAAAAAAAAAAAAAAAAAAAAD+/u65AUMH576+xZgGZexMmaWcbhxwSkHm0SMt1Xh4Z",
			Key:        "AQIBAHjjCeNKbXvpFm2HyIXoKqPvjyyCDoSbmOkO8uSNKqBm0wG7UlVRaAOaTAbr9fqcochxAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMCupa+YgNV/EQizynAgEQgDuQ14WOf0gVBLYvjfXMAREMzygPGVfmNxRsx7SExpyHTRl0tQIE//7LztBvIFCjN7PNuDWACQuiHcW1Yw==`",
			IV:         "DRE6jVdAjY/cIsVsl1/lug==",
		},
	}
)
