# Fechamento Automático de Leilão

## Objetivo
Adicionar uma nova funcionalidade ao projeto existente para que o **leilão seja fechado automaticamente** após um tempo definido.

## Tarefas
Você deverá implementar:

1. **Função para cálculo do tempo do leilão**  
   - O tempo será baseado em parâmetros definidos em variáveis de ambiente.

2. **Nova goroutine**  
   - Validará a existência de leilões vencidos (tempo expirado).  
   - Atualizará o status do leilão, fechando-o.

3. **Teste automatizado**  
   - Garantirá que o fechamento automático esteja funcionando corretamente.

### Para rodar o banco de dados via docker `docker compose up --build -d`

### Para rodar a aplicação `go run cmd/auctions/main.go`

### Para realizar as chamadas http, basta abrir o arquivo `auctions.http` (com a extensão REST Client no VSCode)
- enviar request POST para criar o leilão
- coletar o ID da response e passar para a variável auctionId
- enviar request GET para obter os dados do leilão

#### O tempo de expiração do leilão está definido para **15s**

### Para rodar os testes, basta rodar o comando `go test ./tests/integration/...`
- o **testcontainers** vai usar uma imagem do mongo
- de preferência rode antes dos testes o comando `docker pull mongo` para baixar a imagem

