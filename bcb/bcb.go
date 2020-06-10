package bcb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	SPREAD float32 = 4.0
	IOF float32 = 6.38
)

type cotacaoDolarFetch struct {
	Context string          `json:"@odata.context"`
	Cotacoes []CotacaoDolar `json:"value"`
}

type CotacaoDolar struct {
	Compra float32 `json:"cotacaoCompra"`
	Venda float32 `json:"cotacaoVenda"`
	DataHora string `json:"dataHoraCotacao"`
	TotalCartaoCredito float32 `json:"totalCartaoCredito"`
}

func (c *CotacaoDolar) calcularTotal() {
	spread := (c.Venda * SPREAD) / 100
	iof := ((c.Venda + spread) * IOF) / 100
	c.TotalCartaoCredito = c.Venda + spread + iof
}

func ObterCotacao(data time.Time, cotacao chan CotacaoDolar) (err error) {
	url := "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata/CotacaoDolarDia(dataCotacao=@dataCotacao)?@dataCotacao='" + data.Format("01-02-2006") + "'&$top=1&$format=json&$select=cotacaoCompra,cotacaoVenda,dataHoraCotacao"
	request, err := http.Get(url)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	fetch := new(cotacaoDolarFetch)
	if err = json.Unmarshal(body, &fetch); err != nil {
		return
	}

	if len(fetch.Cotacoes) == 0 {
		ObterCotacao(data.Add(-1 * time.Hour * 24), cotacao)
		return
	}

	c := fetch.Cotacoes[0]
	c.calcularTotal()
	cotacao <- c

	return
}
