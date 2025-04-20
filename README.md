# 🌦️ Weather by CEP - GoLang

Um microserviço em Go que recebe um **CEP**, consulta o endereço via [ViaCEP](https://viacep.com.br), identifica a **cidade** e retorna a **temperatura atual** em **Celsius**, **Fahrenheit** e **Kelvin**, utilizando a [WeatherAPI](https://www.weatherapi.com/).

---

## ✅ Funcionalidades

- Consulta de endereço a partir de um CEP válido.
- Consulta do clima atual com conversões de temperatura.
- API REST com tratamento de erros.
- Testes automatizados.
- Deploy via Docker e Google Cloud Run.

---

## 📦 Requisitos

- Go 1.20+
- Docker instalado
- Conta gratuita no [WeatherAPI](https://www.weatherapi.com/)
- Conta no [Google Cloud](https://cloud.google.com/)

---

## 🛠️ Variáveis de Ambiente

Crie um arquivo `.env` na pasta `cmd/weather/` com o seguinte conteúdo:

```
cep_api_url=https://viacep.com.br/ws
weather_api_url=http://api.weatherapi.com/v1/current.json
weather_api_key=SUA_CHAVE_DA_WEATHERAPI
```

---

## 🚀 Como rodar localmente com Docker



```
docker-compose up --build
```

A API estará disponível em:  
**http://localhost:8080**

---

## 🧪 Testando a API

### Sucesso

```
curl http://localhost:8080/?cep=01001000
```

**Resposta esperada:**

```json
{
  "celsius": 24.1,
  "fahrenheit": 75.38,
  "kelvin": 297.25
}
```

### Erros tratados

- **CEP inválido** (formato errado): `422 - invalid zipcode`
- **CEP não encontrado**: `404 - can not find zipcode`
- **Parâmetro ausente**: `400 - cep is required`

---

## 🧪 Rodando testes

```
go test ./...
```