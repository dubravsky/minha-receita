package transform

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cuducos/go-cnpj"
)

var companyNameClenupRegex = regexp.MustCompile(`(\D)(\d{3})(\d{5})(\d{3})$`) // masks CPF in MEI names

func companyNameClenup(n string) string {
	return strings.TrimSpace(companyNameClenupRegex.ReplaceAllString(n, "$1***$3***"))
}

type company struct {
	CNPJ                             string        `json:"cnpj"`
	IdentificadorMatrizFilial        *int          `json:"branch_type"`
	DescricaoMatrizFilial            *string       `json:"branch_type_description"`
	NomeFantasia                     string        `json:"trade_mark"`
	SituacaoCadastral                *int          `json:"registration_status_code"`
	DescricaoSituacaoCadastral       *string       `json:"registration_status"`
	DataSituacaoCadastral            *date         `json:"registration_update_date"`
	MotivoSituacaoCadastral          *int          `json:"closing_status_code"`
	DescricaoMotivoSituacaoCadastral *string       `json:"closing_status_reason"`
	NomeCidadeNoExterior             string        `json:"international_city_name"`
	CodigoPais                       *int          `json:"country_code"`
	Pais                             *string       `json:"country"`
	DataInicioAtividade              *date         `json:"activity_start_date"`
	CNAEFiscal                       *int          `json:"cnae_code"`
	CNAEFiscalDescricao              *string       `json:"cnae_description"`
	DescricaoTipoDeLogradouro        string        `json:"street_type"`
	Logradouro                       string        `json:"street"`
	Numero                           string        `json:"number"`
	Complemento                      string        `json:"additional"`
	Bairro                           string        `json:"bairro"`
	CEP                              string        `json:"cep"`
	UF                               string        `json:"uf"`
	CodigoMunicipio                  *int          `json:"municipality_code"`
	CodigoMunicipioIBGE              *int          `json:"ibge_municipality_code"`
	Municipio                        *string       `json:"municipality"`
	Telefone1                        string        `json:"phone1"`
	Telefone2                        string        `json:"phone2"`
	Fax                              string        `json:"fax"`
	Email                            *string       `json:"email"`
	SituacaoEspecial                 string        `json:"special_code"`
	DataSituacaoEspecial             *date         `json:"special_situation_date"`
	OpcaoPeloSimples                 *bool         `json:"simple_taxes_status"`
	DataOpcaoPeloSimples             *date         `json:"simple_taxes_start_date"`
	DataExclusaoDoSimples            *date         `json:"simple_taxes_exclusion_date"`
	OpcaoPeloMEI                     *bool         `json:"individual_taxpayer_status"`
	DataOpcaoPeloMEI                 *date         `json:"individual_taxpayer_start_date"`
	DataExclusaoDoMEI                *date         `json:"individual_taxpayer_delete_date"`
	RazaoSocial                      string        `json:"full_name"`
	CodigoNaturezaJuridica           *int          `json:"legal_entity_type_code"`
	NaturezaJuridica                 *string       `json:"legal_type"`
	QualificacaoDoResponsavel        *int          `json:"personal_responsability_code"`
	CapitalSocial                    *float32      `json:"charter_capital"`
	CodigoPorte                      *int          `json:"business_size_code"`
	Porte                            *string       `json:"business_size"`
	EnteFederativoResponsavel        string        `json:"responsible_federative_entity"`
	DescricaoPorte                   string        `json:"business_size_description"`
	QuadroSocietario                 []partnerData `json:"qsa"`
	CNAESecundarios                  []cnae        `json:"additional_cnae"`
}

func (c *company) situacaoCadastral(v string) error {
	i, err := toInt(v)
	if err != nil {
		return fmt.Errorf("error trying to parse SituacaoCadastral %s: %w", v, err)
	}

	var s string
	switch *i {
	case 1:
		s = "NULA"
	case 2:
		s = "ATIVA"
	case 3:
		s = "SUSPENSA"
	case 4:
		s = "INAPTA"
	case 8:
		s = "BAIXADA"
	}

	c.SituacaoCadastral = i
	if s != "" {
		c.DescricaoSituacaoCadastral = &s
	}
	return nil
}

func (c *company) identificadorMatrizFilial(v string) error {
	i, err := toInt(v)
	if err != nil {
		return fmt.Errorf("error trying to parse IdentificadorMatrizFilial %s: %w", v, err)
	}

	var s string
	switch *i {
	case 1:
		s = "MATRIZ"
	case 2:
		s = "FILIAL"
	}

	c.IdentificadorMatrizFilial = i
	if s != "" {
		c.DescricaoMatrizFilial = &s
	}
	return nil
}

func newCompany(row []string, l *lookups, kv kvStorage, privacy bool) (company, error) {
	var c company
	c.CNPJ = row[0] + row[1] + row[2]
	c.NomeFantasia = row[4]
	c.NomeCidadeNoExterior = row[8]
	c.DescricaoTipoDeLogradouro = row[13]
	c.Logradouro = row[14]
	c.Numero = row[15]
	c.Complemento = row[16]
	c.Bairro = row[17]
	c.CEP = row[18]
	c.UF = row[19]
	c.Telefone1 = row[21] + row[22]
	c.Telefone2 = row[23] + row[24]
	c.Fax = row[25] + row[26]
	c.Email = &row[27]
	c.SituacaoEspecial = row[28]

	if privacy {
		c.NomeFantasia = companyNameClenup(row[4])
		c.Email = nil
		if c.CodigoNaturezaJuridica != nil && strings.Contains(strings.ToLower(*c.NaturezaJuridica), "individual") {
			c.DescricaoTipoDeLogradouro = ""
			c.Logradouro = ""
			c.Numero = ""
			c.Complemento = ""
			c.Telefone1 = ""
			c.Telefone2 = ""
			c.Fax = ""
		}
	}

	if err := c.identificadorMatrizFilial(row[3]); err != nil {
		return c, fmt.Errorf("error trying to parse IdentificadorMatrizFilial: %w", err)
	}

	if err := c.situacaoCadastral(row[5]); err != nil {
		return c, fmt.Errorf("error trying to parse SituacaoCadastral: %w", err)
	}

	dataSituacaoCadastral, err := toDate(row[6])
	if err != nil {
		return c, fmt.Errorf("error trying to parse DataSituacaoCadastral %s: %w", row[3], err)
	}
	c.DataSituacaoCadastral = dataSituacaoCadastral

	if err := c.motivoSituacaoCadastral(l, row[7]); err != nil {
		return c, fmt.Errorf("error trying to parse MotivoSituacaoCadastral: %w", err)
	}

	if err := c.pais(l, row[9]); err != nil {
		return c, fmt.Errorf("error trying to parse CodigoPais: %w", err)
	}

	dataInicioAtividade, err := toDate(row[10])
	if err != nil {
		return c, fmt.Errorf("error trying to parse DataInicioAtividade %s: %w", row[10], err)
	}
	c.DataInicioAtividade = dataInicioAtividade

	if err := c.cnaes(l, row[11], row[12]); err != nil {
		return c, fmt.Errorf("error trying to parse cnae: %w", err)
	}

	if err := c.municipio(l, row[20]); err != nil {
		return c, fmt.Errorf("error trying to parse CodigoMunicipio %s: %w", row[20], err)
	}

	dataSituacaoEspecial, err := toDate(row[29])
	if err != nil {
		return c, fmt.Errorf("error trying to parse DataSituacaoEspecial %s: %w", row[20], err)
	}
	c.DataSituacaoEspecial = dataSituacaoEspecial

	if err := kv.enrichCompany(&c); err != nil {
		return c, fmt.Errorf("error enriching company %s: %w", cnpj.Mask(c.CNPJ), err)
	}
	return c, nil
}

func (c *company) JSON() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("error while mashaling company JSON: %w", err)
	}
	return string(b), nil
}
