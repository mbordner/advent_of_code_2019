package nanofactory

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	ORE  = "ORE"
	FUEL = "FUEL"
)

type Chemical struct {
	ID       string
	Quantity uint64
}

func (c *Chemical) init(q, id string) {
	n, e := strconv.Atoi(q)
	if e != nil {
		panic(e)
	}
	c.Quantity = uint64(n)
	c.ID = id
}

type InputChemical struct {
	Chemical
}

func NewInputChemical(q, id string) *InputChemical {
	ic := new(InputChemical)
	ic.init(q, id)
	return ic
}

type OutputChemical struct {
	Chemical
	Inputs []*InputChemical
}

func (oc *OutputChemical) ProducedFromOre() bool {
	if len(oc.Inputs) == 1 {
		if oc.Inputs[0].ID == ORE {
			return true
		}
	}
	return false
}

func NewOutputChemical(q, id string) *OutputChemical {
	oc := new(OutputChemical)
	oc.init(q, id)
	oc.Inputs = make([]*InputChemical, 0, 10)
	return oc
}

func (oc *OutputChemical) AddInput(ic *InputChemical) {
	oc.Inputs = append(oc.Inputs, ic)
}

type ProductionChemical struct {
	CreatedForProduction uint64
	OreRequired          uint64
	Extra                uint64
}

func (pc *ProductionChemical) Clone() *ProductionChemical {
	o := NewProductionChemical()
	o.CreatedForProduction = pc.CreatedForProduction
	o.OreRequired = pc.OreRequired
	o.Extra = pc.Extra
	return o
}

func NewProductionChemical() *ProductionChemical {
	pc := new(ProductionChemical)
	return pc
}

type Production map[string]*ProductionChemical

func (p Production) Clone() Production {
	o := make(Production)
	for k,v := range p {
		o[k] = v.Clone()
	}
	return o
}

func (p Production) ClearUsed() {
	for _, pc := range p {
		pc.CreatedForProduction = 0
		pc.OreRequired = 0
	}
}

func (p Production) Scale( scale uint64 ) {
	for _, pc := range p {
		pc.OreRequired *= scale
		pc.CreatedForProduction *= scale
		pc.Extra *= scale
	}
}

func (p Production) OreUsed() uint64 {
	oreUsed := uint64(0)
	for _, pc := range p {
		oreUsed += pc.OreRequired
	}
	return oreUsed
}

func (p Production) GetProductionChemical(id string) *ProductionChemical {
	if pc, ok := p[id]; ok {
		return pc
	}
	return nil
}

type Factory struct {
	reactionsList map[string]*OutputChemical
}

func (f *Factory) GetOutputChemical(id string) *OutputChemical {
	if oc, ok := f.reactionsList[id]; ok {
		return oc
	}
	return nil
}


func (f *Factory) FuelFromOre( ore uint64 ) uint64 {
	fuel := uint64(0)

	production := make(Production)

	batch := uint64(1000)
	for batch > ore {
		batch /= uint64(10)
	}

	for {
		savedProduction := production.Clone()
		c := f.CostInOre(FUEL, batch, production)
		if c <= ore {
			ore -= c
			fuel += batch
			fmt.Println("produced ",fuel," fuel, with ", ore, " ore remaining.")
			production.ClearUsed()
		} else {
			if batch > uint64(1) {
				production = savedProduction
				batch /= uint64(10)
			} else {
				break
			}
		}
	}
	return fuel
}

func (f *Factory) CostInOre(id string, quantity uint64, prevProduction Production) uint64 {
	var production Production
	if prevProduction != nil {
		production = prevProduction
	} else {
		production = make(Production)
	}

	oc := f.GetOutputChemical(id)

	opc := production.GetProductionChemical(oc.ID)
	if opc == nil {
		production[oc.ID] = NewProductionChemical()
		opc = production[oc.ID]
	}

	if oc.ProducedFromOre() {

		if opc.Extra >= quantity {
			opc.CreatedForProduction += quantity
			opc.Extra -= quantity
		} else {
			produced := opc.Extra
			opc.Extra = 0

			needed := quantity - produced

			batches := uint64(math.Ceil(float64(needed) / float64(oc.Quantity)))
			produced += batches * oc.Quantity

			opc.OreRequired += batches * oc.Inputs[0].Quantity // we know this is ORE
			opc.CreatedForProduction += quantity
			opc.Extra += produced - quantity
		}

	} else {
		produced := opc.Extra
		opc.Extra = 0

		for produced < quantity {
			for _, ic := range oc.Inputs {

				ipc := production.GetProductionChemical(ic.ID)
				if ipc == nil {
					production[ic.ID] = NewProductionChemical()
					ipc = production[ic.ID]
				}

				if ipc.Extra >= ic.Quantity {
					ipc.CreatedForProduction += ic.Quantity
					ipc.Extra -= ic.Quantity
				} else {
					f.CostInOre(ic.ID, ic.Quantity, production)
				}

			}
			produced += oc.Quantity
			opc.CreatedForProduction += oc.Quantity
		}

		opc.Extra += produced - quantity

	}

	return production.OreUsed()
}

func NewFactory(reactionsList string) *Factory {
	reComponents := regexp.MustCompile(`(\d+)\s([A-Z]+)`)

	f := new(Factory)
	f.reactionsList = make(map[string]*OutputChemical)

	list := strings.Split(reactionsList, "\n")
	for _, reaction := range list {
		matches := reComponents.FindAllStringSubmatch(reaction, -1)
		oconf := matches[len(matches)-1]
		oc := NewOutputChemical(oconf[1], oconf[2])
		for _, iconf := range matches[:len(matches)-1] {
			ic := NewInputChemical(iconf[1], iconf[2])
			oc.AddInput(ic)
		}
		f.reactionsList[oc.ID] = oc
	}
	return f
}
