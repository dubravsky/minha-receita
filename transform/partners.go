package transform

import (
	"encoding/json"
	"fmt"
)

type partnerData struct {
	IdentificadorDeSocio                 *int    `json:"partner_code"`
	NomeSocio                            string  `json:"name"`
	CNPJCPFDoSocio                       string  `json:"cnpj_cpf"`
	CodigoQualificacaoSocio              *int    `json:"qualification_code"`
	QualificaoSocio                      *string `json:"qualification"`
	DataEntradaSociedade                 *date   `json:"company_entry_date"`
	CodigoPais                           *int    `json:"country_code"`
	Pais                                 *string `json:"code"`
	CPFRepresentanteLegal                string  `json:"cpf_legal_representative"`
	NomeRepresentanteLegal               string  `json:"name_legal_representative"`
	CodigoQualificacaoRepresentanteLegal *int    `json:"legal_representative_qualification_code"`
	QualificacaoRepresentanteLegal       *string `json:"legal_representative_qualification"`
	CodigoFaixaEtaria                    *int    `json:"age_range_code"`
	FaixaEtaria                          *string `json:"age_range"`
}

func (p *partnerData) faixaEtaria(v string) {
	c, err := toInt(v)
	if err != nil {
		return
	}
	p.CodigoFaixaEtaria = c

	var s string
	switch *c {
	case 1:
		s = "para os intervalos entre 0 a 12 anos"
	case 2:
		s = "Entre 13 a 20 ano"
	case 3:
		s = "Entre 21 a 30 anos"
	case 4:
		s = "Entre 31 a 40 anos"
	case 5:
		s = "Entre 41 a 50 anos"
	case 6:
		s = "Entre 51 a 60 anos"
	case 7:
		s = "Entre 61 a 70 anos"
	case 8:
		s = "Entre 71 a 80 anos"
	case 9:
		s = "Maiores de 80 anos"
	case 0:
		s = "Não se aplica"
	}
	if s != "" {
		p.FaixaEtaria = &s
	}
}

func (p *partnerData) pais(l *lookups, v string) error {
	i, err := toInt(v)
	if err != nil {
		return fmt.Errorf("error trying to parse CodigoPais %s: %w", v, err)
	}
	if i == nil {
		return nil
	}
	s := l.countries[*i]
	p.CodigoPais = i
	if s != "" {
		p.Pais = &s
	}
	return nil
}

func newPartnerData(l *lookups, r []string) (partnerData, error) {
	identificadorDeSocio, err := toInt(r[1])
	if err != nil {
		return partnerData{}, fmt.Errorf("error parsing IdentificadorDeSocio %s: %w", r[1], err)
	}

	dataEntradaSociedade, err := toDate(r[5])
	if err != nil {
		return partnerData{}, fmt.Errorf("error parsing DataEntradaSociedade %s: %w", r[5], err)
	}

	p := partnerData{
		IdentificadorDeSocio:   identificadorDeSocio,
		NomeSocio:              r[2],
		CNPJCPFDoSocio:         r[3],
		DataEntradaSociedade:   dataEntradaSociedade,
		CPFRepresentanteLegal:  r[7],
		NomeRepresentanteLegal: r[8],
	}
	p.pais(l, r[6])
	p.faixaEtaria(r[10])
	p.qualificacaoSocio(l, r[4], r[9])
	return p, nil
}

func loadPartnerRow(l *lookups, r []string) ([]byte, error) {
	p, err := newPartnerData(l, r)
	if err != nil {
		return nil, fmt.Errorf("error parsing taxes line: %w", err)
	}
	v, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling base: %w", err)
	}
	return v, nil
}
