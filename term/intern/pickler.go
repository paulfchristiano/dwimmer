package intern

type Pickler interface {
	Pickle(Packer) Packed
	Unpickle(Packer, Packed) (Pickler, bool)
	Key() interface{} //Must be comparable
	Test(Pickler) bool
}

/*
//NOTE: if a packer stores both proxies and non-proxies,
//there might be conflicts with keys
type proxy struct {
	Pickler Pickler
	Keyer   Packer
}

func (p *proxy) Key() interface{} {
	return p.Keyer.PackPickler(p.Pickler)
}

func (p *proxy) Pickle(packer Packer) Packed {
	return p.Pickler.Pickle(packer)
}

func (p *proxy) Test(pickler Pickler) bool {
	other, ok := pickler.(*proxy)
	return ok && p.Pickler.Test(other.Pickler) && p.Keyer == other.Keyer
}

func (p *proxy) Unpickle(packer Packer, pickled Packed) (Pickler, bool) {
	var ok bool
	p.Pickler, ok = p.Pickler.Unpickle(packer, pickled)
	return p, ok
}

type Cacher struct {
	Packer
	Keyer Packer
}

func (k *Cacher) PackPickler(pickler Pickler) Packed {
	fmt.Println("[Cacher] Packing", pickler)
	return k.Packer.PackPickler(&proxy{pickler, k.Keyer})
}

func (k *Cacher) UnpackPickler(x Packed, pickler Pickler) (Pickler, bool) {
	result, ok := k.Packer.UnpackPickler(x, &proxy{pickler, k.Keyer})
	if !ok {
		return nil, false
	}
	return result.(*proxy).Pickler, true
}
*/
