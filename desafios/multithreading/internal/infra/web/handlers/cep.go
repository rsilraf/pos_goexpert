package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"golang.org/x/exp/rand"
)

type CepInterface interface {
	ToCepInfo() *CepInfo
}
type CepInfo struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	API         string `json:"api"`
}

type viaCepAPI struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func (v *viaCepAPI) ToCepInfo() *CepInfo {
	return &CepInfo{
		Cep:         v.Cep,
		Logradouro:  v.Logradouro,
		Complemento: v.Complemento,
		Unidade:     v.Unidade,
		Bairro:      v.Bairro,
		Localidade:  v.Localidade,
		Uf:          v.Uf,
		Ibge:        v.Ibge,
		Gia:         v.Gia,
		Ddd:         v.Ddd,
		Siafi:       v.Siafi,
		API:         "ViaCEP",
	}
}

type brasilAPIAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func (b *brasilAPIAPI) ToCepInfo() *CepInfo {
	return &CepInfo{
		Cep:        b.Cep,
		Uf:         b.State,
		Bairro:     b.Neighborhood,
		Logradouro: b.Street,
		API:        "BrasilAPI",
	}
}

const HTTP_TIMEOUT = 1 * time.Second

type APIs struct{}

type APICall struct {
	URL   string
	Body  []byte
	Cep   *CepInfo
	Error error
}

type CepHandler struct{}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

// GetCEP	godoc
// @Title		Get CEP information
// @Tags		cep
// @Accept		json
// @Produce		json
// @Param 		cep		path		string	true	"cep" default(01311200)
// @Success		200		{object}	CepInfo
// @Failure		400		{object}	Error
// @Failure		500		{object}	Error
// @Security	ApiKeyAuth
// @Router		/cep/{cep}	[get]
func (h *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	_, err := strconv.Atoi(cep)

	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	cepInfo, err := getCepInfo(cep)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&cepInfo)
}

func (a *APIs) request(ctx context.Context, url string, cepStruct CepInterface) *APICall {

	call := &APICall{
		URL: url,
	}

	// request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		call.Error = err
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		call.Error = err
		return nil
	}
	defer resp.Body.Close()

	call.Body, err = io.ReadAll(resp.Body)

	if err != nil {
		call.Error = err
	}

	if json.Unmarshal(call.Body, cepStruct) != nil {
		return nil
	}

	call.Cep = cepStruct.ToCepInfo()

	return call

}

func (a *APIs) BrasilAPI(ctx context.Context, cep string) chan *APICall {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	ch := make(chan *APICall)
	go func() {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		ch <- a.request(ctx, url, &brasilAPIAPI{})
	}()
	return ch
}

func (a *APIs) ViaCEP(ctx context.Context, cep string) chan *APICall {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	ch := make(chan *APICall)
	go func() {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		ch <- a.request(ctx, url, &viaCepAPI{})
	}()
	return ch
}

func getCepInfo(cep string) (*CepInfo, error) {

	apis := &APIs{}
	ctx, cancel := context.WithTimeout(context.Background(), HTTP_TIMEOUT)
	defer cancel()

	resA := apis.BrasilAPI(ctx, cep)
	resB := apis.ViaCEP(ctx, cep)

	var answer *APICall

	select {
	case answer = <-resA:
		cancel()
	case answer = <-resB:
		cancel()

	case <-time.After(HTTP_TIMEOUT):
		return nil, errors.New("timeout")
	}

	if answer.Error != nil {
		log.Println(answer.Error)
		return nil, answer.Error
	}
	fmt.Println(answer)
	return answer.Cep, nil
}
