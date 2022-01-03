package types

type Extractor interface {
	Get(query string) ([]Song, error)
}

type Dispatcher interface {
	Play(gID, vID, cmdID string, songs []Song)
	Queue(gID string) []Song
	Seek(gID string, seekTime int) (string, error)
	SkipTo(gID string, pos int) (string, error)
	Stop(gID string) (string, error)
	Skip(gID string) (string, error)
	Pause(gID string) (string, error)
	Resume(gID string) (string, error)
}
