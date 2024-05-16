package rotors

const (
	Rot2Prog string = "rot2prog"
)

type Rotor interface {
	Init() error
	GetPos() (float64, float64, error)
	SetPos(float64, float64) error
}

