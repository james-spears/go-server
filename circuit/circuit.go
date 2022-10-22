package circuit

type Attributes struct {
	NO      bool  `json:"no"`
	State   bool  `json:"state"`
	Address []int `json:"address"`
}

func (a *Attributes) IsNormallyClosed() bool {
	return !a.NO
}

func (a *Attributes) IsNormallyOpen() bool {
	return a.NO
}

func (a *Attributes) SetNormallyOpen() {
	a.NO = true
}

func (a *Attributes) SetNormallyClosed() {
	a.NO = false
}

func (a *Attributes) ToggleDefaultPosition() {
	a.NO = !a.NO
}

func (a *Attributes) IsNotEnergized() bool {
	return !a.State
}

func (a *Attributes) IsEnergized() bool {
	return a.State
}

func (a *Attributes) Energize() {
	a.State = true
}

func (a *Attributes) Deenergize() {
	a.State = false
}

func (a *Attributes) TogglePower() {
	a.State = !a.State
}

func (a *Attributes) IsOpen() bool {
	return a.NO == a.State
}

func (a *Attributes) IsClosed() bool {
	return a.NO != a.State
}

type Circuit struct {
	Name       string     `json:"name"`
	Attributes Attributes `json:"attributes"`
	Children   []Circuit  `json:"children"`
}

func (t *Circuit) Append(circuits ...Circuit) {
	if len(circuits) > 0 {
		firstChild := circuits[0]

		for i := 0; i < len(firstChild.Children); i++ {
			addr := append(t.Attributes.Address, len(t.Children))
			firstChild.Children[i].Attributes.Address = addr
		}

		t.Children = append(t.Children, circuits...)
		t.Append(circuits[1:]...)
	}
}
